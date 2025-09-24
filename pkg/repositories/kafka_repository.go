package repositories

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	entities "github.com/banyar/go-packages/pkg/entities"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaRepository struct {
	producer *KafkaProducer
	consumer *KafkaConsumer
}

type KafkaConsumer struct {
	Client  *kafka.Consumer
	Mu      sync.RWMutex
	Running bool
	Closed  chan struct{}
}

type KafkaProducer struct {
	Client  *kafka.Producer
	Mu      sync.Mutex
	Running bool
}

func ConnectKafka() *KafkaRepository {
	return &KafkaRepository{}
}

func (kr *KafkaRepository) ConnectProducer(producerDSN *entities.KafkaProducerDSN) error {
	producerConfig := createProducerConfig(producerDSN)
	producer, err := kafka.NewProducer(producerConfig)
	if err != nil {
		return fmt.Errorf("producer creation failed: %w", err)
	}

	kr.producer = &KafkaProducer{
		Client:  producer,
		Running: true,
	}
	return nil
}

func (kr *KafkaRepository) ConnectConsumer(consumerDSN *entities.KafkaConsumerDSN) error {
	consumerConfig := createConsumerConfig(consumerDSN)
	consumer, err := kafka.NewConsumer(consumerConfig)
	if err != nil {
		return fmt.Errorf("consumer creation failed: %w", err)
	}

	kr.consumer = &KafkaConsumer{
		Client:  consumer,
		Running: true,
		Closed:  make(chan struct{}),
	}
	return nil
}

func (kr *KafkaRepository) ProducerClose() {
	if kr.producer == nil {
		return
	}

	kr.producer.Mu.Lock()
	defer kr.producer.Mu.Unlock()

	if kr.producer.Running {
		kr.producer.Client.Close()
		kr.producer.Running = false
	}
}

func (kr *KafkaRepository) ConsumerClose() {
	if kr.consumer == nil {
		return
	}

	kr.consumer.Mu.Lock()
	defer kr.consumer.Mu.Unlock()

	if kr.consumer.Running {
		close(kr.consumer.Closed)
		kr.consumer.Client.Close()
		kr.consumer.Running = false
	}
}

// Getter Producer methods
func (kr *KafkaRepository) GetProducer() *KafkaProducer {
	return kr.producer
}

// Getter Consumer methods
func (kr *KafkaRepository) GetConsumer() *KafkaConsumer {
	return kr.consumer
}

// Helper functions to create configs from entity structs
func createProducerConfig(dsn *entities.KafkaProducerDSN) *kafka.ConfigMap {

	retries, err := strconv.Atoi(dsn.Retries)
	if err != nil {
		// Set default value if conversion fails or empty
		retries = 3 // Default retry count
	}

	return &kafka.ConfigMap{
		"bootstrap.servers": dsn.Brokers,
		"client.id":         getWithDefault(dsn.ClientID, "default-consumer"),
		"security.protocol": getWithDefault(dsn.Protocol, "PLAINTEXT"),
		"sasl.mechanism":    dsn.Mechanism,
		"sasl.username":     dsn.Username,
		"sasl.password":     dsn.Password,
		"acks":              getWithDefault(dsn.Acks, "all"),
		"retries":           retries,
		"compression.type":  getWithDefault(dsn.Compression, "snappy"),
	}
}

func createConsumerConfig(dsn *entities.KafkaConsumerDSN) *kafka.ConfigMap {
	return &kafka.ConfigMap{
		"bootstrap.servers":  dsn.Brokers,
		"group.id":           dsn.GroupID,
		"client.id":          getWithDefault(dsn.ClientID, "default-consumer"),
		"security.protocol":  getWithDefault(dsn.Protocol, "PLAINTEXT"),
		"sasl.mechanism":     dsn.Mechanism,
		"sasl.username":      dsn.Username,
		"sasl.password":      dsn.Password,
		"auto.offset.reset":  getWithDefault(dsn.AutoOffset, "earliest"),
		"enable.auto.commit": getEnvAsBool(dsn.AutoCommit, false),
	}
}

func getWithDefault(key, defaultValue string) string {
	val := key
	if val == "" {
		return defaultValue
	}
	return val
}

// Helper function for boolean environment variables
func getEnvAsBool(key string, defaultValue bool) bool {
	val := strings.ToLower(key)
	switch val {
	case "true", "1", "yes", "on":
		return true
	case "false", "0", "no", "off":
		return false
	default:
		return defaultValue
	}
}
