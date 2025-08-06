package adapters

import (
	entities "github.com/banyar/go-packages/pkg/entities"
	"github.com/banyar/go-packages/pkg/interfaces"
	"github.com/banyar/go-packages/pkg/repositories"
	"github.com/banyar/go-packages/pkg/services"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQAdapter struct {
	RabbitMQService interfaces.IRabbitMQService
	Conn            *amqp091.Connection
}

func NewRabbitMQAdapter(DSNRBQ *entities.DSNRabbitMQ, poolSize int) *RabbitMQAdapter {
	rbqRepo := repositories.ConnectRabbitMQ(DSNRBQ, poolSize)
	return &RabbitMQAdapter{
		RabbitMQService: services.NewRabbitMQService(rbqRepo),
		Conn:            rbqRepo.Conn,
	}
}
