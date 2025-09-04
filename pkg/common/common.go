package common

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/banyar/go-packages/pkg/entities"

	"golang.org/x/exp/rand"
)

// letters is a constant string of characters that the random string will use.
const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func DisplayJsonFormat(message string, class interface{}) {
	if message == "" {
		message = "CLASS"
	}
	p, err := json.Marshal(class)
	LogError(err)
	log.Println(message+" ===>", string(p))
}

func LogError(err error) {
	if err != nil {
		log.Fatal("ERROR : ", err)
	}
}

func EnsureOutputFolderExists(folderPath string) error {
	// Check if the folder already exists
	_, err := os.Stat(folderPath)
	if err == nil {
		// Folder already exists, nothing to do
		return nil
	}

	// If the error is not due to the folder not existing, return the error
	if !os.IsNotExist(err) {
		return err
	}

	// Folder doesn't exist, create it
	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return err
	}

	// Folder created successfully
	return nil
}

type UnitConfig struct {
	Prefix    string
	Suffix    string
	Separator string
	AddRandom bool
}

func GenerateUnitName(data string, config UnitConfig) string {
	timestampStr := TimestampString()

	parts := []string{}

	if config.Prefix != "" {
		parts = append(parts, config.Prefix)
	}

	parts = append(parts, data)
	parts = append(parts, timestampStr)

	if config.Suffix != "" {
		parts = append(parts, config.Suffix)
	}

	if config.AddRandom {
		randomStr := GenerateRandomString(4)
		parts = append(parts, randomStr)
	}

	separator := config.Separator
	if separator == "" {
		separator = "_"
	}

	return strings.Join(parts, separator)
}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func PreparePayload(data string) entities.CustomPublish {
	// Payload
	timestampStr := TimestampString()
	config := UnitConfig{Prefix: "unit", Suffix: "prod", Separator: "-", AddRandom: true}
	unitName := GenerateUnitName(data, config)
	publish := entities.CustomPublish{
		Name:      fmt.Sprintf("user_%s_%s", unitName, timestampStr),
		Email:     fmt.Sprintf("%s.%s@gmail.com", data, timestampStr),
		Timestamp: timestampStr,
		UserID:    fmt.Sprintf("usr_%d", time.Now().UnixNano()),
		IsActive:  true,
		CreatedAt: time.Now(),
		Metadata: map[string]interface{}{
			"source":  "dynamic_generator",
			"version": "1.0",
		},
	}
	return publish
}

func TimestampString() string {
	return time.Now().Format("20060102150405")
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func CreateAndSetStruct(fields map[string]reflect.Type, values map[string]interface{}) interface{} {
	structType := CreateStructType(fields)
	structValue := reflect.New(structType).Elem()

	for name, value := range values {
		field := structValue.FieldByName(name)
		if field.IsValid() && field.CanSet() {
			val := reflect.ValueOf(value)
			if field.Type() == val.Type() {
				field.Set(val)
			}
		}
	}

	return structValue.Interface()
}

// CreateStructType creates a struct type with the given field names and types
func CreateStructType(fields map[string]reflect.Type) reflect.Type {
	structFields := make([]reflect.StructField, len(fields))
	i := 0
	for name, fieldType := range fields {
		structFields[i] = reflect.StructField{
			Name: name,
			Type: fieldType,
			Tag:  reflect.StructTag(`json:"` + name + `"`),
		}
		i++
	}
	return reflect.StructOf(structFields)
}

func RandString(length int) string {
	src := rand.NewSource(uint64(time.Now().UnixNano()))

	rnd := rand.New(src)
	// Create a slice of runes to hold the characters of the random string.
	r := make([]rune, length)
	// Generate random characters from the letters constant and add them to the slice.
	for i := range r {
		r[i] = rune(letters[rnd.Intn(len(letters))])
	}
	// Convert the slice of runes to a string and return it.
	return string(r)
}

func RandInt() int {

	rand.NewSource(uint64(time.Now().UnixNano()))

	// Generate a random integer between 0 and 99 (inclusive)
	randomInt := rand.Intn(100)
	fmt.Println("Random Integer between 0 and 99:", randomInt)

	// Generate a random integer between min and max (inclusive)
	min := 10
	max := 100
	randomIntInRange := min + rand.Intn(max-min+1)
	return randomIntInRange
}

func GetDynamicPayLoad() interface{} {

	fields := map[string]reflect.Type{
		"Name":  reflect.TypeOf(""),
		"Age":   reflect.TypeOf(0),
		"Email": reflect.TypeOf(""),
	}

	// Define the values to set in the struct
	values := map[string]interface{}{
		"Name":  RandString(10),
		"Age":   RandInt(),
		"Email": RandString(10) + "@gmail.com",
	}

	// Create and set the struct
	dynamicStruct := CreateAndSetStruct(fields, values)
	fmt.Printf("Dynamic Struct: %+v\n", dynamicStruct)
	return dynamicStruct
}

func GetDSNRabbitMQ() entities.DSNRabbitMQ {
	rbqPort, err := strconv.Atoi(os.Getenv("RABBIT_MQ_PORT"))
	if err != nil {
		log.Fatalf("Error converting RBQ to integer: %v", err)
	}
	return entities.DSNRabbitMQ{
		Host:         os.Getenv("RABBIT_MQ_HOST"),
		User:         os.Getenv("RABBIT_MQ_USER"),
		Password:     os.Getenv("RABBIT_MQ_PASSWORD"),
		Port:         rbqPort,
		RoutingKey:   os.Getenv("RABBIT_MQ_ROUTING_KEY"),
		Queue:        os.Getenv("RABBIT_MQ_QUEUE"),
		Exchange:     os.Getenv("RABBIT_MQ_EXCHANGE"),
		ExchangeType: os.Getenv("RABBIT_MQ_EXCHANGE_TYPE"),
		ContentType:  os.Getenv("RABBIT_MQ_CONTENT_TYPE"),
		VirtualHost:  os.Getenv("RABBIT_MQ_VIRTUAL_HOST"),
	}
}

func GetDSNRabbitMQ1() entities.DSNRabbitMQ {
	rbqPort, err := strconv.Atoi(os.Getenv("RABBIT_MQ_PORT1"))
	if err != nil {
		log.Fatalf("Error converting RBQ to integer: %v", err)
	}
	return entities.DSNRabbitMQ{
		Host:         os.Getenv("RABBIT_MQ_HOST"),
		User:         os.Getenv("RABBIT_MQ_USER"),
		Password:     os.Getenv("RABBIT_MQ_PASSWORD"),
		Port:         rbqPort,
		RoutingKey:   os.Getenv("RABBIT_MQ_ROUTING_KEY"),
		Queue:        os.Getenv("queue2"),
		Exchange:     os.Getenv("RABBIT_MQ_EXCHANGE"),
		ExchangeType: os.Getenv("RABBIT_MQ_EXCHANGE_TYPE"),
		ContentType:  os.Getenv("RABBIT_MQ_CONTENT_TYPE"),
	}
}
