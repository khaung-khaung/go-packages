package interfaces

type IRabbitMQService interface {
	Produce(payload interface{}) (int, string)
}
