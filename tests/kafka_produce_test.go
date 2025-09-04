package tests

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/banyar/go-packages/pkg/adapters"
	"github.com/banyar/go-packages/pkg/common"
	"github.com/banyar/go-packages/pkg/frontlog"
	"github.com/banyar/go-packages/pkg/repositories"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap"
)

func TestKafkaProduce(t *testing.T) {

	producerDSN := common.ProducerDSNFromEnv()

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// logger.Info("Producing test message...", zap.Any("", producerDSN))

	// Create adapter with error handling
	kafkaAdapter := adapters.NewKafkaAdapter()

	// Defer close immediately after successful creation
	defer func() {
		kafkaAdapter.KafkaRepo.ProducerClose()
	}()

	// Access producer safely
	kafkaAdapter.KafkaRepo.ConnectProducer(producerDSN)
	producer := kafkaAdapter.KafkaRepo.GetProducer()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	currentTime := time.Now()

	// Custom format: "YYYY-MM-DD HH:MM:SS"
	formatted := currentTime.Format("2006-01-02 15:04:05")

	payloadObj := common.PreparePayload(formatted)

	// Produce messages to topic (asynchronously)
	payload, err := json.Marshal(payloadObj)
	if err != nil {
		frontlog.Logger.Error("Error converting struct to JSON", zap.Any(":", err))
	}
	jsonString := string(payload)

	headers := []kafka.Header{
		{Key: "Trace-ID", Value: []byte("12345")},
		{Key: "Source", Value: []byte("order-service")},
		{Key: "Event-Type", Value: []byte("order.created")},
	}

	err = SendMessage(producer, ctx, os.Getenv("KAFKA_TOPIC"), []byte(jsonString), headers)
	if err != nil {

		frontlog.Logger.Error("Failed to send message", zap.Any(":", err))
	}

}

func SendMessage(kp *repositories.KafkaProducer, ctx context.Context, topic string, message []byte, headers []kafka.Header) error {
	kp.Mu.Lock()
	defer kp.Mu.Unlock()

	if !kp.Running {
		return errors.New("producer is closed")
	}

	deliveryChan := make(chan kafka.Event, 1)
	defer close(deliveryChan)

	// Create message with headers
	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value:   message,
		Headers: headers,
	}

	err := kp.Client.Produce(kafkaMsg, deliveryChan)

	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case ev := <-deliveryChan:
		msg := ev.(*kafka.Message)
		if msg.TopicPartition.Error != nil {
			return fmt.Errorf("delivery failed: %w", msg.TopicPartition.Error)
		}

		frontlog.Logger.Info("", zap.Any("Delivered message", msg.TopicPartition))

		return nil
	}
}
