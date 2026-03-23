//go:build cgo
// +build cgo

package interfaces

type IKafkaService interface {
	// GetConsumer() (*kafka.Consumer, error)
	// GetProducer() (*kafka.Producer, error)
}
