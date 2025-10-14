package entities

type HttpResponse struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"msg"`
}

type KafkaResponse struct {
	Status         string      `json:"status"`
	Message        string      `json:"msg"`
	TopicPartition interface{} `json:"topic_partition"`
}

func (h HttpResponse) Error() string {
	return h.Message
}
