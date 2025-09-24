package tests

// import (
// 	"fmt"
// 	"os"
// 	"testing"

// 	"github.com/banyar/go-packages/pkg/adapters"
// 	"github.com/banyar/go-packages/pkg/common"
// 	"github.com/banyar/go-packages/pkg/entities"
// )

// func TestKafkaConsume(t *testing.T) {
// 	DSNKafka := GetDSNKafka()
// 	common.DisplayJsonFormat("DSNKafka", &DSNKafka)
// 	kafkaAdapter := adapters.NewKafkaAdapter(&DSNKafka)
// 	message, header := kafkaAdapter.KafkaService.Consume()
// 	fmt.Println("message ", message, header)
// }

// func TestKafkaProduce(t *testing.T) {
// 	DSNKafka := GetDSNKafka()
// 	common.DisplayJsonFormat("DSNKafka", &DSNKafka)
// 	kafkaAdapter := adapters.NewKafkaAdapter(&DSNKafka)
// 	payload := common.PreparePayload()
// 	status, message, topic := kafkaAdapter.KafkaService.Produce(payload)
// 	fmt.Println("message ", status, message, topic)
// }

// func GetDSNKafka() entities.DSNKafka {
// 	return entities.DSNKafka{
// 		Brokers:   os.Getenv("KAFKA_HOST"),
// 		Topics:    os.Getenv("KAFKA_TOPIC"),
// 		GroupId:   os.Getenv("KAFKA_GROUP_ID"),
// 		User:      os.Getenv("KAFKA_USER"),
// 		Password:  os.Getenv("KAFKA_PASSWORD"),
// 		Mechanism: os.Getenv("KAFKA_MECHANISM"),
// 		Protocol:  os.Getenv("KAFKA_PROTOCOL"),
// 	}
// }
