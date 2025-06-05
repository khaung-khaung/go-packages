package repositories

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/banyar/go-packages/pkg/common"
	entities "github.com/banyar/go-packages/pkg/entities"
	"github.com/streadway/amqp"
)

type RabbitMQRepository struct {
	Conn        *amqp.Connection
	dsnRBQ      *entities.DSNRabbitMQ
	channelPool chan *amqp.Channel
	maxPoolSize int
	mu          sync.Mutex
}

// Initialize the repository with a connection and channel pool
func ConnectRabbitMQ(DSNRBQ *entities.DSNRabbitMQ, poolSize int) *RabbitMQRepository {
	rabbitURL := fmt.Sprintf("amqp://%s:%s@%s:%d%s", DSNRBQ.User, DSNRBQ.Password, DSNRBQ.Host, DSNRBQ.Port, DSNRBQ.VirtualHost)
	conn, err := amqp.Dial(rabbitURL)
	common.FailOnError(err, "Failed to connect to RabbitMQ")

	repo := &RabbitMQRepository{
		Conn:        conn,
		dsnRBQ:      DSNRBQ,
		channelPool: make(chan *amqp.Channel, poolSize),
		maxPoolSize: poolSize,
	}

	// Pre-fill the channel pool
	for i := 0; i < poolSize; i++ {
		ch, err := repo.createChannel()
		if err != nil {
			common.FailOnError(err, "Failed to create channel")
		}
		repo.channelPool <- ch
	}

	return repo
}

// Create a new channel
func (r *RabbitMQRepository) createChannel() (*amqp.Channel, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	ch, err := r.Conn.Channel()
	if err != nil {
		return nil, err
	}
	return ch, nil
}

// Get a channel from the pool, waiting if necessary
func (r *RabbitMQRepository) GetChannel() (*amqp.Channel, error) {
	select {
	case ch := <-r.channelPool:
		return ch, nil
	default:
		// If no channels are available, wait for one to be released
		for {
			time.Sleep(500 * time.Millisecond) // Adjust the wait time as needed
			select {
			case ch := <-r.channelPool:
				return ch, nil
			default:
				continue
			}
		}
	}
}

// Release the channel back to the pool
func (r *RabbitMQRepository) ReleaseChannel(ch *amqp.Channel) {
	r.channelPool <- ch
}

// Post a message to the queue
func (r *RabbitMQRepository) PostMessage(payloadObj interface{}, headers map[string]interface{}) (int, string) {
	ch, err := r.GetChannel()
	if err != nil {
		return 500, "Failed to get channel: " + err.Error()
	}
	defer r.ReleaseChannel(ch)

	payload, err := json.Marshal(payloadObj)
	if err != nil {
		return 500, "Error converting struct to JSON"
	}
	messageBody := string(payload)

	err = r.getExchangeDeclare(ch)
	if err != nil {
		return 500, "Failed to declare exchange: " + err.Error()
	}
	if len(r.dsnRBQ.Queue) > 0 {
		_, err = r.getQueueDeclare(ch)
		if err != nil {
			return 500, "Failed to declare queue: " + err.Error()
		}
	}

	err = ch.Publish(
		r.dsnRBQ.Exchange,
		r.dsnRBQ.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  r.dsnRBQ.ContentType,
			Body:         []byte(messageBody),
			DeliveryMode: 2,
			Headers:      headers,
		})

	if err != nil {
		return 500, "Delivery failed: " + err.Error()
	}

	return 200, "Delivered successfully"
}

func (r *RabbitMQRepository) getQueueDeclare(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		r.dsnRBQ.Queue,
		true,
		false,
		false,
		false,
		nil,
	)
}

func (r *RabbitMQRepository) getExchangeDeclare(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		r.dsnRBQ.Exchange,
		r.dsnRBQ.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
}
