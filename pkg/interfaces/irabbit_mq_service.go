package interfaces

import "github.com/rabbitmq/amqp091-go"

type IRabbitMQService interface {
	Produce(payload any, headers map[string]any) (int, string)
	Consume(fn func(*amqp091.Delivery) bool)
}
