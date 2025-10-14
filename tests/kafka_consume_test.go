package tests

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/banyar/go-packages/pkg/adapters"
	"github.com/banyar/go-packages/pkg/common"
	"github.com/banyar/go-packages/pkg/frontlog"
	"github.com/banyar/go-packages/pkg/repositories"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap"
)

func TestKafkaConsume(t *testing.T) {
	consumerDSN := common.ConsumerDSNFromEnv()

	kafkaAdapter := adapters.NewKafkaAdapter()
	defer kafkaAdapter.KafkaRepo.ConsumerClose()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	kafkaAdapter.KafkaRepo.ConnectConsumer(consumerDSN)
	consumer := kafkaAdapter.KafkaRepo.GetConsumer()

	// Use wait group for proper synchronization
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		Consume(consumer, ctx, "fibermaps")
	}()

	select {
	case <-ctx.Done():
		t.Log("Test completed successfully")
	case <-time.After(25 * time.Second):
		t.Error("Test timed out waiting for messages")
	}

	// Wait for consumer to fully stop
	wg.Wait()
}
func Consume(kc *repositories.KafkaConsumer, ctx context.Context, topic string) {
	kc.Mu.Lock()
	if !kc.Running {
		kc.Mu.Unlock()
		return
	}

	// Use temporary unsubscribe to prevent race
	defer func() {
		kc.Mu.Lock()
		kc.Client.Unsubscribe()
		kc.Mu.Unlock()
	}()

	err := kc.Client.SubscribeTopics([]string{topic}, nil)
	kc.Mu.Unlock()

	if err != nil {
		frontlog.Logger.Error("Subscribe error: :", zap.Any("", err))
		return
	}

	log.Printf("Consuming from %s", topic)

	for {
		select {
		case <-ctx.Done():
			frontlog.Logger.Info("Context cancelled - stopping consumer")
			return
		case <-kc.Closed: // Listen to closure signal
			frontlog.Logger.Info("Consumer closed - stopping")

			return
		default:
			kc.Mu.RLock()
			ev := kc.Client.Poll(100)
			kc.Mu.RUnlock()

			if ev == nil {
				continue
			}

			// Handle message in a goroutine with panic recovery
			go func(e kafka.Event) {
				defer func() {
					if r := recover(); r != nil {
						frontlog.Logger.Info("Recovered from panic:", zap.Any("", r))
					}
				}()

				switch msg := e.(type) {
				case *kafka.Message:
					handleMessage(kc, msg)
				case kafka.Error:
					handleKafkaError(msg)
				}
			}(ev)
		}
	}
}

func handleMessage(kc *repositories.KafkaConsumer, msg *kafka.Message) {
	headers := headersToMap(msg.Headers)

	log.Printf("Received message:\n"+
		"  Topic: %s\n"+
		"  Partition: %d\n"+
		"  Offset: %d\n"+
		"  Key: %s\n"+
		"  Headers: %+v\n"+
		"  Value: %s\n",
		*msg.TopicPartition.Topic,
		msg.TopicPartition.Partition,
		msg.TopicPartition.Offset,
		string(msg.Key),
		headers,
		string(msg.Value),
	)

	// Process message with headers
	err := processMessage(msg.Value, headers)
	if err != nil {
		log.Printf("Message processing failed: %v", err)
		return
	}

	// Commit message offset
	_, err = kc.Client.CommitMessage(msg)
	if err != nil {
		log.Printf("Failed to commit message offset: %v", err)
	}
}

func headersToMap(headers []kafka.Header) map[string]string {
	headerMap := make(map[string]string)
	for _, h := range headers {
		headerMap[h.Key] = string(h.Value)
	}
	return headerMap
}

func processMessage(body []byte, headers map[string]string) error {
	// Your business logic here
	// Access headers like headers["Trace-ID"], headers["Content-Type"], etc.
	return nil
}

func handleKafkaError(e kafka.Error) {
	if e.Code() == kafka.ErrAllBrokersDown {
		log.Fatal("Fatal Kafka error: all brokers down")
	}
	log.Printf("Kafka error: %v (%v)", e.Code(), e)
}
