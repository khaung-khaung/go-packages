package services

import (
	"github.com/banyar/go-packages/pkg/repositories"
)

type rabbitMQService struct {
	RabbitMQRepo *repositories.RabbitMQRepository
}

func NewRabbitMQService(rbqRepo *repositories.RabbitMQRepository) *rabbitMQService {
	return &rabbitMQService{
		RabbitMQRepo: rbqRepo,
	}
}

func (r *rabbitMQService) Produce(payload interface{}) (int, string) {
	return r.RabbitMQRepo.PostMessage(payload)
}
