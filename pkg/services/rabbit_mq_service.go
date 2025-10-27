package services

import (
	"github.com/banyar/go-packages/pkg/repositories"
	"github.com/rabbitmq/amqp091-go"
)

type rabbitMQService struct {
	RabbitMQRepo *repositories.RabbitMQRepository
}

func NewRabbitMQService(rbqRepo *repositories.RabbitMQRepository) *rabbitMQService {
	return &rabbitMQService{
		RabbitMQRepo: rbqRepo,
	}
}

func (r *rabbitMQService) Produce(payload any, headers map[string]any) (int, string) {
	return r.RabbitMQRepo.PostMessage(payload, headers)
}

func (r *rabbitMQService) Consume(fn func(*amqp091.Delivery) bool) {
	r.RabbitMQRepo.ConsumerLoop(fn)
}
