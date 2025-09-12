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
