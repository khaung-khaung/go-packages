package adapters

import (
	"github.com/banyar/go-packages/pkg/repositories"
)

type KafkaAdapter struct {
	KafkaRepo *repositories.KafkaRepository
}

func NewKafkaAdapter() *KafkaAdapter {
	kafkaRepo := repositories.ConnectKafka()
	return &KafkaAdapter{
		KafkaRepo: kafkaRepo,
	}
}
