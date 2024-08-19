package interfaces

type IRabbitMQService interface {
	Produce(payload interface{}, headers map[string]interface{}) (int, string)
}
