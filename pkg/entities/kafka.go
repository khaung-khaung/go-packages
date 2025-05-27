package entities

type KafkaCommonDSN struct {
	Brokers   string // Comma-separated list of broker addresses
	ClientID  string // Client identifier
	Protocol  string // Security protocol (e.g., SASL_SSL)
	Mechanism string // SASL mechanism (e.g., PLAIN, SCRAM-SHA-256)
	Username  string // SASL username
	Password  string // SASL password
}

type KafkaProducerDSN struct {
	KafkaCommonDSN        // Embedded common configuration
	Acks           string // Acknowledgement mode (e.g., "all")
	Retries        string // Number of retries for failed sends
	Compression    string // Compression type (e.g., "snappy")
}

type KafkaConsumerDSN struct {
	KafkaCommonDSN        // Embedded common configuration
	GroupID        string // Consumer group identifier
	AutoOffset     string // Offset reset policy (e.g., "earliest")
	AutoCommit     string // Enable/disable auto-commit of offsets
}
