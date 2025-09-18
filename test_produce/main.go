package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/banyar/go-packages/pkg/adapters"
	"github.com/banyar/go-packages/pkg/entities"
)

type TestPayload struct {
	Index   int    `json:"index"`
	Message string `json:"message"`
}

func generateRandomText(length int) string {
	const charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+"
	b := make([]byte, length)
	for i := range b {
		b[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(b)
}

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
	producer := adapters.NewRabbitMQAdapter(&dsn, 4).RabbitMQService
	if producer == nil {
		log.Fatalf("Failed to create producer")
	}

	minLength := 30
	maxLength := 40
	var index int

	for {
		randomLength := rand.Intn(maxLength-minLength+1) + minLength
		randomText := generateRandomText(randomLength)

		fmt.Println("Publishing message ", index)

		//go func(i int) {
		code, message := producer.Produce(TestPayload{Index: index, Message: randomText}, nil)
		fmt.Printf("Response: %s,  Code: %d\n", message, code)
		//}(index)
		time.Sleep(time.Second)
		index++
	}
}
