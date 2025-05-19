package common

import (
	"os"
	"strconv"
	"strings"

	"github.com/banyar/go-packages/pkg/entities"
)

func ProducerDSNFromEnv() *entities.KafkaProducerDSN {
	// Handle numeric conversion with error checking
	retries, err := strconv.Atoi(os.Getenv("KAFKA_RETRIES"))
	if err != nil {
		// Set default value if conversion fails or empty
		retries = 3 // Default retry count
	}

	return &entities.KafkaProducerDSN{
		KafkaCommonDSN: entities.KafkaCommonDSN{
			Brokers:   os.Getenv("KAFKA_HOST"),
			ClientID:  os.Getenv("KAFKA_GROUP_ID"),
			Protocol:  os.Getenv("KAFKA_PROTOCOL"),
			Mechanism: os.Getenv("KAFKA_MECHANISM"),
			Username:  os.Getenv("KAFKA_USERNAME"),
			Password:  os.Getenv("KAFKA_PASSWORD"),
		},
		Acks:        getEnvWithDefault("KAFKA_ACKS", "all"),
		Retries:     retries,
		Compression: getEnvWithDefault("KAFKA_COMPRESSION", "snappy"),
	}
}

func getEnvWithDefault(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}

func ConsumerDSNFromEnv() *entities.KafkaConsumerDSN {
	return &entities.KafkaConsumerDSN{
		KafkaCommonDSN: entities.KafkaCommonDSN{
			Brokers:   os.Getenv("KAFKA_HOST"),
			ClientID:  getEnvWithDefault("KAFKA_CLIENT_ID", "default-consumer"),
			Protocol:  getEnvWithDefault("KAFKA_PROTOCOL", "PLAINTEXT"),
			Mechanism: os.Getenv("KAFKA_MECHANISM"),
			Username:  os.Getenv("KAFKA_USERNAME"),
			Password:  os.Getenv("KAFKA_PASSWORD"),
		},
		GroupID:    os.Getenv("KAFKA_GROUP_ID"),
		AutoOffset: getEnvWithDefault("KAFKA_AUTO_OFFSET", "earliest"),
		AutoCommit: getEnvAsBool("KAFKA_AUTO_COMMIT", false),
	}
}

// Helper function for boolean environment variables
func getEnvAsBool(key string, defaultValue bool) bool {
	val := strings.ToLower(os.Getenv(key))
	switch val {
	case "true", "1", "yes", "on":
		return true
	case "false", "0", "no", "off":
		return false
	default:
		return defaultValue
	}
}
