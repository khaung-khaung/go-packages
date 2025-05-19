package repositories

import (
	"fmt"
	"sync"

	entities "github.com/banyar/go-packages/pkg/entities"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaRepository struct {
	producer *KafkaProducer // Use pointer
	consumer *KafkaConsumer // Use pointer
}

type KafkaConsumer struct {
	Client  *kafka.Consumer
	Mu      sync.RWMutex
	Running bool
	Closed  chan struct{} // Add closed channel
}

type KafkaProducer struct {
	Client  *kafka.Producer
	Running bool
	Mu      sync.Mutex
}

// ConnectKafka establishes connections using entity configurations
func ConnectKafka(
	producerDSN *entities.KafkaProducerDSN,
	consumerDSN *entities.KafkaConsumerDSN,
) (*KafkaRepository, error) {

	// Create producer config from entities
	producerConfig := createProducerConfig(producerDSN)
	producer, err := kafka.NewProducer(producerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	// Create consumer config from entities
	consumerConfig := createConsumerConfig(consumerDSN)
	consumer, err := kafka.NewConsumer(consumerConfig)
	if err != nil {
		producer.Close()
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	return &KafkaRepository{
		producer: &KafkaProducer{
			Client:  producer,
			Running: true,
		},
		consumer: &KafkaConsumer{
			Client:  consumer,
			Running: true,
			Closed:  make(chan struct{}), // Initialize closed channel

		},
	}, nil
}

func (r *KafkaRepository) Producer() *KafkaProducer {
	return r.producer
}

func (r *KafkaRepository) Consumer() *KafkaConsumer {
	return r.consumer
}

// Helper functions to create configs from entity structs
func createProducerConfig(dsn *entities.KafkaProducerDSN) *kafka.ConfigMap {
	return &kafka.ConfigMap{
		"bootstrap.servers": dsn.Brokers,
		"client.id":         dsn.ClientID,
		"security.protocol": dsn.Protocol,
		"sasl.mechanism":    dsn.Mechanism,
		"sasl.username":     dsn.Username,
		"sasl.password":     dsn.Password,
		"acks":              dsn.Acks,
		"retries":           dsn.Retries,
		"compression.type":  dsn.Compression,
	}
}

func createConsumerConfig(dsn *entities.KafkaConsumerDSN) *kafka.ConfigMap {
	return &kafka.ConfigMap{
		"bootstrap.servers":  dsn.Brokers,
		"group.id":           dsn.GroupID,
		"client.id":          dsn.ClientID,
		"security.protocol":  dsn.Protocol,
		"sasl.mechanism":     dsn.Mechanism,
		"sasl.username":      dsn.Username,
		"sasl.password":      dsn.Password,
		"auto.offset.reset":  dsn.AutoOffset,
		"enable.auto.commit": false,
	}
}

// Close safely shuts down connections
func (kr *KafkaRepository) ProducerClose() {
	kr.producer.Mu.Lock()
	defer kr.producer.Mu.Unlock()
	if kr.producer.Running {
		kr.producer.Client.Close()
		kr.producer.Running = false
	}
}

func (kr *KafkaRepository) ConsumerClose() {
	kr.consumer.Mu.Lock()
	defer kr.consumer.Mu.Unlock()

	if kr.consumer.Running {
		close(kr.consumer.Closed) // Signal closure
		kr.consumer.Client.Close()
		kr.consumer.Running = false
	}
}
