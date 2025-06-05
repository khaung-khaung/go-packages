package services

// import (
// 	"github.com/banyar/go-packages/pkg/repositories"
// 	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
// )

// type kafkaService struct {
// 	KafkaRepo *repositories.KafkaRepository
// }

// func NewKafkaService(kafkaRepo *repositories.KafkaRepository) *kafkaService {
// 	return &kafkaService{
// 		KafkaRepo: kafkaRepo,
// 	}
// }

// func (k *kafkaService) Consume() (string, []kafka.Header) {
// 	return k.KafkaRepo.GetMessage()
// }

// func (k *kafkaService) Produce(payload interface{}) (int, string, kafka.TopicPartition) {
// 	return k.KafkaRepo.PostMessage(payload)
// }
