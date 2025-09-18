package main

import (
	"fmt"
	"log"
	"time"

	"github.com/banyar/go-packages/pkg/adapters"
	"github.com/banyar/go-packages/pkg/entities"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	dsn := entities.DSNRabbitMQ{
		Host:         "localhost",
		Port:         5672,
		User:         "guest",
		Password:     "guest",
		VirtualHost:  "/",
		Exchange:     "test.fanout.integration",
		ExchangeType: "fanout",
		Queue:        "test.fanout.integration.queue",
		Timeout:      30,
	}
	consumer := adapters.NewRabbitMQAdapter(&dsn, 4).RabbitMQService
	if consumer == nil {
		log.Fatalf("Failed to create consumer")
	}

	// Start consumer loop
	consumer.Consume(func(delivery *amqp091.Delivery) bool {
		// Print the received message details
		fmt.Println("--- New Message Received ---")
		// fmt.Printf("Message ID: %s\n", delivery.MessageId)
		// fmt.Printf("App ID: %s\n", delivery.AppId)
		// fmt.Printf("Timestamp: %v\n", delivery.Timestamp)
		// fmt.Println("Headers:")
		// for key, value := range delivery.Headers {
		// 	fmt.Printf("  - %s: %v\n", key, value)
		// }
		fmt.Println("        Body: ", string(delivery.Body))
		// fmt.Println("----------------------------")
		time.Sleep(1 * time.Second)
		fmt.Println("--- Job finished ---")
		return true
	})
}
