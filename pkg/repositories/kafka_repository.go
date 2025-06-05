package repositories

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"os"

// 	entities "github.com/banyar/go-packages/pkg/entities"
// 	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
// )

// type KafkaRepository struct {
// 	configMap *kafka.ConfigMap
// 	dsnKafka  *entities.DSNKafka
// }

// func ConnectKafka(DSNKafka *entities.DSNKafka) *KafkaRepository {

// 	config := &kafka.ConfigMap{
// 		"bootstrap.servers": DSNKafka.Brokers,
// 		"security.protocol": DSNKafka.Protocol,
// 		"sasl.mechanism":    DSNKafka.Mechanism,
// 		"sasl.username":     DSNKafka.User,
// 		"sasl.password":     DSNKafka.Password,
// 		"group.id":          DSNKafka.GroupId,
// 		"auto.offset.reset": "earliest",
// 	}

// 	return &KafkaRepository{
// 		configMap: config,
// 		dsnKafka:  DSNKafka,
// 	}
// }

// func (r *KafkaRepository) GetMessage() (string, []kafka.Header) {
// 	var messageValue string
// 	var messageHeader []kafka.Header
// 	consumer, err := kafka.NewConsumer(r.configMap)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
// 		os.Exit(1)
// 		return err.Error(), messageHeader
// 	}

// 	consumer.SubscribeTopics([]string{r.dsnKafka.Topics}, nil)
// 	run := true

// 	for run {
// 		ev := consumer.Poll(100)
// 		switch e := ev.(type) {
// 		case *kafka.Message:
// 			messageValue = string(e.Value)
// 			fmt.Printf("Message on %s: %s\n", e.TopicPartition, string(e.Value))
// 			if e.Headers != nil {
// 				fmt.Printf("Headers: %v\n", e.Headers)
// 				messageHeader = e.Headers
// 			}

// 		case kafka.Error:
// 			fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
// 			run = false
// 			messageValue = e.Error()
// 		default:
// 			// fmt.Printf("Ignored %v\n", e)
// 		}
// 	}

// 	consumer.Close()
// 	return messageValue, messageHeader

// }

// func (r *KafkaRepository) PostMessage(payloadObj interface{}) (int, string, kafka.TopicPartition) {

// 	// Produce messages to topic (asynchronously)
// 	payload, err := json.Marshal(payloadObj)
// 	jsonString := string(payload)

// 	// fmt.Println("PayLoad Json", jsonString)
// 	if err != nil {
// 		log.Fatalf("Error converting struct to JSON: %s", err)
// 	}
// 	p, err := kafka.NewProducer(r.configMap)
// 	if err != nil {
// 		panic(err)
// 	}

// 	var topicPartition kafka.TopicPartition
// 	var statusMessage string
// 	var statusCode int = 0

// 	defer p.Close()

// 	// Delivery report handler for produced messages
// 	go func() {
// 		for e := range p.Events() {
// 			switch ev := e.(type) {
// 			case *kafka.Message:
// 				if ev.TopicPartition.Error != nil {
// 					// fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
// 					statusCode = 500
// 					statusMessage = "Delivery failed"
// 					topicPartition = ev.TopicPartition
// 				} else {
// 					// fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
// 					statusCode = 200
// 					statusMessage = "Delivered success"
// 					topicPartition = ev.TopicPartition
// 				}
// 			}
// 		}
// 	}()

// 	for _, word := range []string{jsonString} {
// 		p.Produce(&kafka.Message{
// 			TopicPartition: kafka.TopicPartition{Topic: &r.dsnKafka.Topics, Partition: kafka.PartitionAny},
// 			Value:          []byte(word),
// 			Headers: []kafka.Header{
// 				{Key: "group_id", Value: []byte(r.dsnKafka.GroupId)}, // we will decide later for header value set
// 			},
// 		}, nil)
// 	}

// 	// Wait for message deliveries before shutting down
// 	p.Flush(15 * 1000)

// 	return statusCode, statusMessage, topicPartition
// }
