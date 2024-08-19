package adapters

import (
	entities "github.com/banyar/go-packages/pkg/entities"
	"github.com/banyar/go-packages/pkg/interfaces"
	"github.com/banyar/go-packages/pkg/repositories"
	"github.com/banyar/go-packages/pkg/services"

	"github.com/streadway/amqp"
)

type RabbitMQAdapter struct {
	RabbitMQService interfaces.IRabbitMQService
	Conn            *amqp.Connection
}

func NewRabbitMQAdapter(DSNMySQL *entities.DSNRabbitMQ, poolSize int) *RabbitMQAdapter {
	rbqRepo := repositories.ConnectRabbitMQ(DSNMySQL, poolSize)
	return &RabbitMQAdapter{
		RabbitMQService: services.NewRabbitMQService(rbqRepo),
		Conn:            rbqRepo.Conn,
	}
}
