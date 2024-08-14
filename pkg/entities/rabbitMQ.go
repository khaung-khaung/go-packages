package entities

type DSNRabbitMQ struct {
	Host          string
	User          string
	Password      string
	RoutingKey    string
	Port          int
	Queue         string
	Exchange      string
	ExchangeType  string
	ContentType   string
	StatusMessage string
	VirtualHost   string
}
