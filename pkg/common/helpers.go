package common

import (
	"os"

	"github.com/banyar/go-packages/pkg/entities"
)

func ProducerDSNFromEnv() *entities.KafkaProducerDSN {
	// Handle numeric conversion with error checking

	return &entities.KafkaProducerDSN{
		KafkaCommonDSN: entities.KafkaCommonDSN{
			Brokers:   os.Getenv("KAFKA_HOST"),
			ClientID:  os.Getenv("KAFKA_GROUP_ID"),
			Protocol:  os.Getenv("KAFKA_PROTOCOL"),
			Mechanism: os.Getenv("KAFKA_MECHANISM"),
			Username:  os.Getenv("KAFKA_USERNAME"),
			Password:  os.Getenv("KAFKA_PASSWORD"),
		},
		Acks:        os.Getenv("KAFKA_ACKS"),
		Retries:     os.Getenv("KAFKA_RETRIES"),
		Compression: os.Getenv("KAFKA_COMPRESSION"),
	}
}

func ConsumerDSNFromEnv() *entities.KafkaConsumerDSN {
	return &entities.KafkaConsumerDSN{
		KafkaCommonDSN: entities.KafkaCommonDSN{
			Brokers:   os.Getenv("KAFKA_HOST"),
			ClientID:  os.Getenv("KAFKA_CLIENT_ID"),
			Protocol:  os.Getenv("KAFKA_PROTOCOL"),
			Mechanism: os.Getenv("KAFKA_MECHANISM"),
			Username:  os.Getenv("KAFKA_USERNAME"),
			Password:  os.Getenv("KAFKA_PASSWORD"),
		},
		GroupID:    os.Getenv("KAFKA_GROUP_ID"),
		AutoOffset: os.Getenv("KAFKA_AUTO_OFFSET"),
		AutoCommit: os.Getenv("KAFKA_AUTO_COMMIT"),
	}
}
