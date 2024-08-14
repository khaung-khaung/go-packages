package entities

type DSNKafka struct {
	Brokers   string
	Topics    string
	GroupId   string
	User      string
	Password  string
	Mechanism string
	Protocol  string
}
