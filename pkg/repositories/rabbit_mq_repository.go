package repositories

import (
	"encoding/json"
	"fmt"

	"github.com/banyar/go-packages/pkg/common"
	entities "github.com/banyar/go-packages/pkg/entities"

	"github.com/streadway/amqp"
)

type RabbitMQRepository struct {
	Conn   *amqp.Connection
	dsnRBQ *entities.DSNRabbitMQ
}

func ConnectRabbitMQ(DSNRBQ *entities.DSNRabbitMQ) *RabbitMQRepository {
	// Format the RabbitMQ URL
	rabbitURL := fmt.Sprintf("amqp://%s:%s@%s:%d%s/", DSNRBQ.User, DSNRBQ.Password, DSNRBQ.Host, DSNRBQ.Port, DSNRBQ.VirtualHost)
	// Connect to RabbitMQ server
	conn, err := amqp.Dial(rabbitURL)
	common.FailOnError(err, "Failed to connect to RabbitMQ")
	return &RabbitMQRepository{
		Conn:   conn,
		dsnRBQ: DSNRBQ,
	}
}

func (r *RabbitMQRepository) PostMessage(payloadObj interface{}) (int, string) {
	var statusMessage string
	var statusCode int = 0
	payload, err := json.Marshal(payloadObj)
	common.FailOnError(err, "Error converting struct to JSON")
	messageBody := string(payload)
	fmt.Println("messageBody", messageBody)
	common.DisplayJsonFormat("PostMessage dsnRBQ", r.dsnRBQ)
	defer r.Conn.Close()
	// Create a channel
	ch, err := r.Conn.Channel()
	if err != nil {
		fmt.Printf("Failed to open a channel: %v\n", ch)
		statusCode = 500
		statusMessage = "Failed to open a channel"
	}

	defer ch.Close()
	// Declare a exchange with custom settings
	err = r.getExchangeDeclare(ch)
	if err != nil {
		fmt.Printf("Failed to declare a exchange: %v\n", ch)
		statusCode = 500
		statusMessage = "Failed to declare a exchange"
	}

	// Declare a queue with custom settings
	q, err := r.getQueueDeclare(ch)
	if err != nil {
		fmt.Printf("Failed to declare a queue: %v\n", q)
		statusCode = 500
		statusMessage = "Failed to declare a queue"
	}

	// Publish a message with custom settings
	err = ch.Publish(
		r.dsnRBQ.Exchange,   // exchange Use default exchange
		r.dsnRBQ.RoutingKey, // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType: r.dsnRBQ.ContentType,
			Body:        []byte(messageBody),
		})

	if err != nil {
		fmt.Printf("Delivery failed: %v\n", messageBody)
		statusCode = 500
		statusMessage = "Delivery failed" + messageBody

	} else {
		fmt.Printf("Delivered message to %v\n", messageBody)
		statusCode = 200
		statusMessage = "Delivered success" + messageBody
	}
	return statusCode, statusMessage
}

func (r *RabbitMQRepository) getQueueDeclare(ch *amqp.Channel) (amqp.Queue, error) {
	// Declare an queue
	return ch.QueueDeclare(
		r.dsnRBQ.Queue, // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
}

func (r *RabbitMQRepository) getExchangeDeclare(ch *amqp.Channel) error {
	// Declare an exchange
	return ch.ExchangeDeclare(
		r.dsnRBQ.Exchange,     // name
		r.dsnRBQ.ExchangeType, // type
		true,                  // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)
}
