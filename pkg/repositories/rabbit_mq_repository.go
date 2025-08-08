package repositories

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	entities "github.com/banyar/go-packages/pkg/entities"
	"github.com/rabbitmq/amqp091-go"
)

func isNetworkError(err error) bool {
	// Check if the underlying connection/channel is gone.
	if errors.Is(err, amqp091.ErrClosed) {
		return true
	}

	// Check for network-related errors
	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}

	// Check for io.EOF, which often signals a closed connection or stream
	if errors.Is(err, io.EOF) {
		return true
	}

	return false
}

type RabbitMQRepository struct {
	Conn        *amqp091.Connection
	dsnRBQ      *entities.DSNRabbitMQ
	channelPool chan *amqp091.Channel
	maxPoolSize int
	mu          sync.Mutex
	netErrFlag  int32
	connID      int32
}

// Initialize the repository with a connection and channel pool
func ConnectRabbitMQ(DSNRBQ *entities.DSNRabbitMQ, poolSize int) *RabbitMQRepository {
	host := DSNRBQ.Host
	port := DSNRBQ.Port
	user := DSNRBQ.User
	pass := DSNRBQ.Password
	vhost := DSNRBQ.VirtualHost
	uri := fmt.Sprintf("amqp://%s:%s@%s:%d%s", user, pass, host, port, vhost)
	retryCount := 0
	delay := time.Second

	var conn *amqp091.Connection
	var err error

	// Connect to server
	for {
		conn, err = amqp091.Dial(uri)
		if err == nil {
			break
		} else {
			if retryCount == 0 {
				log.Printf("Failed to connect to RabbitMQ: %v", err)
			}
			if !isNetworkError(err) {
				log.Fatalf("Failed to connect to RabbitMQ: %v", err)
			}
		}

		if retryCount == 6 {
			log.Fatalln("RabbitMQ connection max retries reached")
		}
		retryCount++

		time.Sleep(delay)
		delay *= 2
		if delay > 16*time.Second {
			delay = 16 * time.Second
		}
		log.Printf("Retry %d\n", retryCount)
	}

	repo := &RabbitMQRepository{
		Conn:        conn,
		dsnRBQ:      DSNRBQ,
		maxPoolSize: poolSize,
	}

	// Initialize channels
	repo.channelPool = make(chan *amqp091.Channel, poolSize)
	for i := 0; i < poolSize; i++ {
		ch, err := repo.createChannel()
		if err != nil {
			log.Fatalf("Failed to initialize RabbitMQ channel %d: %v", i, err)
		}
		repo.channelPool <- ch
	}

	return repo
}

// Create a new channel
func (r *RabbitMQRepository) createChannel() (*amqp091.Channel, error) {
	ch, err := r.Conn.Channel()
	if err != nil {
		return nil, err
	}

	// Declare exchange
	err = ch.ExchangeDeclare(
		r.dsnRBQ.Exchange,
		r.dsnRBQ.ExchangeType,
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	// Declare queue
	if len(r.dsnRBQ.Queue) > 0 {
		_, err = ch.QueueDeclare(
			r.dsnRBQ.Queue,
			true,
			false,
			false,
			false,
			nil,
		)
	}

	return ch, err
}

// Get a channel from the pool, waiting if necessary
func (r *RabbitMQRepository) GetChannel() (*amqp091.Channel, error) {
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
func (r *RabbitMQRepository) ReleaseChannel(ch *amqp091.Channel) {
	select {
	case r.channelPool <- ch:
		// Return to pool
	default:
		ch.Close()
	}
}

// Reconnect when a network error occur
func (r *RabbitMQRepository) reconnectRBMQ() bool {
	id := atomic.LoadInt32(&r.connID)
	r.mu.Lock()
	defer r.mu.Unlock()
	if id != atomic.LoadInt32(&r.connID) {
		if atomic.LoadInt32(&r.netErrFlag) == 1 {
			return false
		} else {
			return true
		}
	}
	atomic.StoreInt32(&r.netErrFlag, 2)
	time.Sleep(time.Second)

	// Disconnect
	if !r.Conn.IsClosed() {
		close(r.channelPool)
		for ch := range r.channelPool {
			ch.Close()
		}
		r.Conn.Close()
	}
	defer atomic.AddInt32(&r.connID, 1)

	host := r.dsnRBQ.Host
	port := r.dsnRBQ.Port
	user := r.dsnRBQ.User
	pass := r.dsnRBQ.Password
	vhost := r.dsnRBQ.VirtualHost
	uri := fmt.Sprintf("amqp://%s:%s@%s:%d%s", user, pass, host, port, vhost)
	retryCount := 1
	delay := time.Second
	log.Println("Reconnecting RabbitMQ")

	for {
		log.Printf("Retry %d\n", retryCount)

		// Reconnect
		conn, err := amqp091.Dial(uri)
		if err == nil {
			// Initialize channels
			r.Conn = conn
			r.channelPool = make(chan *amqp091.Channel, r.maxPoolSize)
			for i := 0; i < r.maxPoolSize; i++ {
				ch, err := r.createChannel()
				if err != nil {
					atomic.StoreInt32(&r.netErrFlag, 1)
					return false
				}

				r.channelPool <- ch
			}
			atomic.StoreInt32(&r.netErrFlag, 0)
			return true
		}

		if retryCount == 6 {
			atomic.StoreInt32(&r.netErrFlag, 1)
			return false
		}
		retryCount++

		delay *= 2
		if delay > 16*time.Second {
			delay = 16 * time.Second
		}
		time.Sleep(delay)
	}
}

// Post a message to the queue
func (r *RabbitMQRepository) PostMessage(payloadObj interface{}, headers map[string]interface{}) (int, string) {
	payload, err := json.Marshal(payloadObj)
	if err != nil {
		return 500, "Error converting struct to JSON"
	}

	switch atomic.LoadInt32(&r.netErrFlag) {
	case 1:
		log.Println("RabbitMQ disconnected")
		return 500, "RabbitMQ disconnected"
	case 2:
		r.mu.Lock()
		r.mu.Unlock()
		if atomic.LoadInt32(&r.netErrFlag) == 1 {
			log.Println("RabbitMQ disconnected")
			return 500, "RabbitMQ disconnected"
		}
	}

	for {
		// Get RBMQ channel
		ch, err := r.GetChannel()
		if err != nil {
			return 500, err.Error()
		}

		// Publish message
		err = ch.Publish(
			r.dsnRBQ.Exchange,
			"",
			false,
			false,
			amqp091.Publishing{
				ContentType:  "application/json",
				Headers:      headers,
				Body:         payload,
				DeliveryMode: 2,
			},
		)
		if err != nil {
			log.Printf("Failed to publish message to RabbitMQ: %v\n", err)
			if isNetworkError(err) {
				if ok := r.reconnectRBMQ(); ok {
					continue
				}
			}
			return 500, err.Error()
		}

		// Return channel to channel pool
		r.ReleaseChannel(ch)
		break
	}
	return 200, "Sent successfully"
}
