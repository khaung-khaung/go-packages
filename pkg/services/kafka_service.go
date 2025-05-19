package services

import (
	"github.com/banyar/go-packages/pkg/repositories"
)

type kafkaService struct {
	KafkaRepo *repositories.KafkaRepository
}

func NewKafkaService(kafkaRepo *repositories.KafkaRepository) *kafkaService {
	return &kafkaService{
		KafkaRepo: kafkaRepo,
	}
}

// func (k *kafkaService) GetConsumer() (*kafka.Consumer, error) {
// 	return k.KafkaRepo.CreateConsumer()
// }

// func (k *kafkaService) GetProducer() (*kafka.Producer, error) {
// 	return k.KafkaRepo.CreateProducer()
// }
