package adapters

import (
	entities "github.com/banyar/go-packages/pkg/entities"
	"github.com/banyar/go-packages/pkg/repositories"
)

type KafkaAdapter struct {
	KafkaRepo *repositories.KafkaRepository
}

func NewKafkaAdapter(kafkaProducer *entities.KafkaProducerDSN, kafkaConsumer *entities.KafkaConsumerDSN) *KafkaAdapter {
	kafkaRepo, _ := repositories.ConnectKafka(kafkaProducer, kafkaConsumer)
	return &KafkaAdapter{
		KafkaRepo: kafkaRepo,
	}
}
