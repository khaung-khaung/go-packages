package tests

import (
	"fmt"
	"log"
	"testing"

	"github.com/banyar/go-packages/pkg/adapters"
	"github.com/banyar/go-packages/pkg/common"

	"github.com/streadway/amqp"
)

func TestRabbitProduce(t *testing.T) {
	DSNRBQ := common.GetDSNRabbitMQ()
	rbqAdapter := adapters.NewRabbitMQAdapter(&DSNRBQ)
	payload := common.GetDynamicPayLoad()
	statusCode, message := rbqAdapter.RabbitMQService.Produce(payload)
	fmt.Println("produce message ", statusCode, message)

}

func TestRabbitConsume(t *testing.T) {
	// Custom configuration settings
	DSNRBQ := common.GetDSNRabbitMQ()
	rabbitURL := fmt.Sprintf("amqp://%s:%s@%s:%d/", DSNRBQ.User, DSNRBQ.Password, DSNRBQ.Host, DSNRBQ.Port)
	exchangeName := DSNRBQ.Exchange
	exchangeType := DSNRBQ.ExchangeType
	routingKey := DSNRBQ.RoutingKey
	queueName := DSNRBQ.Queue

	// Connect to RabbitMQ server
	conn, err := amqp.Dial(rabbitURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare an exchange
	err = ch.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	// Declare a queue
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Bind the queue to the exchange with the routing key
	err = ch.QueueBind(
		q.Name,       // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to bind a queue")

	// Consume messages
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	// Create a channel to signal when the program should exit
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			// Process the message
			log.Printf("Received a message: %s", d.Body)

			// Return the message content
			fmt.Println("Returned message:", string(d.Body))
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
