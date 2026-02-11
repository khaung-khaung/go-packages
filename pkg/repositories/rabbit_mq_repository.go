package repositories

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	entities "github.com/banyar/go-packages/pkg/entities"
	"github.com/rabbitmq/amqp091-go"
)

func isTransientError(err error) bool {
	if errors.Is(err, amqp091.ErrClosed) {
		return true
	}
	var amqpErr *amqp091.Error
	if errors.As(err, &amqpErr) {
		if amqpErr.Code == 501 {
			return true
		}
	}
	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}
	if errors.Is(err, io.EOF) {
		return true
	}
	return false
}

const (
	networkStateDisconnected = iota
	networkStateConnected
	networkStateConnecting
	networkStateFailed
)

type RabbitMQRepository struct {
	Conn                *amqp091.Connection
	dsnRBQ              *entities.DSNRabbitMQ
	channelPool         chan *amqp091.Channel
	maxPoolSize         int
	consumerWorkerCount int
	mu                  sync.Mutex
	netErrFlag          int32
	connID              int32
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
	t1 := time.Now()
	for {
		conn, err = amqp091.Dial(uri)
		if err == nil {
			break
		} else {
			if retryCount == 0 {
				log.Printf("Failed to connect to RabbitMQ: %v", err)
			}
			if !isTransientError(err) {
				log.Printf("Failed to connect to RabbitMQ: %v", err)
				return nil
			}
		}

		t2 := time.Now()
		dt := int(t2.Sub(t1).Seconds())
		if dt > DSNRBQ.Timeout {
			log.Println("RabbitMQ connection timeout")
			return nil
		}

		time.Sleep(delay)
		delay *= 2
		if delay > 16*time.Second {
			delay = 16 * time.Second
		}
		retryCount++
		log.Printf("Retry %d\n", retryCount)
	}

	repo := &RabbitMQRepository{
		Conn:                conn,
		dsnRBQ:              DSNRBQ,
		maxPoolSize:         poolSize,
		consumerWorkerCount: poolSize,
	}

	// Initialize channels
	repo.channelPool = make(chan *amqp091.Channel, poolSize)
	for i := range poolSize {
		ch, err := repo.createChannel()
		if err != nil {
			log.Printf("Failed to initialize RabbitMQ channel %d: %v", i, err)
			return nil
		}
		repo.channelPool <- ch
	}
	repo.netErrFlag = networkStateConnected

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
		if atomic.LoadInt32(&r.netErrFlag) == networkStateFailed {
			return false
		} else {
			return true
		}
	}
	atomic.StoreInt32(&r.netErrFlag, networkStateConnecting)
	time.Sleep(time.Second)

	// Disconnect
	if !r.Conn.IsClosed() {
		close(r.channelPool)
		for ch := range r.channelPool {
			ch.Close()
		}
		r.Conn.Close()
	}
	defer atomic.AddInt32(&r.connID, networkStateFailed)

	host := r.dsnRBQ.Host
	port := r.dsnRBQ.Port
	user := r.dsnRBQ.User
	pass := r.dsnRBQ.Password
	vhost := r.dsnRBQ.VirtualHost
	uri := fmt.Sprintf("amqp://%s:%s@%s:%d%s", user, pass, host, port, vhost)
	retryCount := 1
	delay := time.Second
	log.Println("Reconnecting RabbitMQ")

	t1 := time.Now()
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
					atomic.StoreInt32(&r.netErrFlag, networkStateFailed)
					return false
				}

				r.channelPool <- ch
			}
			atomic.StoreInt32(&r.netErrFlag, networkStateConnected)
			return true
		}

		t2 := time.Now()
		dt := int(t2.Sub(t1).Seconds())
		if dt > r.dsnRBQ.Timeout {
			log.Println("RabbitMQ connection timeout")
			atomic.StoreInt32(&r.netErrFlag, networkStateFailed)
			return false
		}

		delay *= 2
		if delay > 16*time.Second {
			delay = 16 * time.Second
		}
		time.Sleep(delay)
		retryCount++
	}
}

// Post a message to the queue
func (r *RabbitMQRepository) PostMessage(payloadObj any, headers map[string]any) (int, string) {
	payload, err := json.Marshal(payloadObj)
	if err != nil {
		return 500, "Error converting struct to JSON"
	}

	switch atomic.LoadInt32(&r.netErrFlag) {
	case networkStateFailed:
		log.Println("RabbitMQ disconnected")
		return 500, "RabbitMQ disconnected"
	case networkStateConnecting:
		r.mu.Lock()
		r.mu.Unlock()
		if atomic.LoadInt32(&r.netErrFlag) == networkStateFailed {
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
			r.dsnRBQ.RoutingKey,
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
			if isTransientError(err) {
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

func (r *RabbitMQRepository) ConsumerLoop(fn func(*amqp091.Delivery) bool) {
	for {
		for {
			time.Sleep(100 * time.Millisecond)
			flag := atomic.LoadInt32(&r.netErrFlag)
			if flag == networkStateConnected {
				break
			} else if flag == networkStateFailed {
				return
			}
		}

		ch, err := r.declareQueue()
		if err != nil {
			if isTransientError(err) {
				if r.reconnectRBMQ() {
					continue
				}
			} else {
				return
			}
		}

		msgs, err := ch.Consume(
			r.dsnRBQ.Queue,
			"",
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			if isTransientError(err) {
				if r.reconnectRBMQ() {
					continue
				}
			}
			log.Printf("Failed to consume message: %v", err)
			break
		}
		var wg sync.WaitGroup
		limiter := make(chan struct{}, r.consumerWorkerCount)
		index := 0
		recovery := true
		for d := range msgs {
			limiter <- struct{}{}
			wg.Add(1)
			go func(i int, msg *amqp091.Delivery) {
				finalized := false

				defer func() {
					wg.Done()
					<-limiter
				}()
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Goroutine %d panicked: %v\n%s", i, r, debug.Stack())
						if !finalized {
							msg.Ack(false)
						}
					}
				}()

				ack := fn(msg)
				finalized = true
				switch atomic.LoadInt32(&r.netErrFlag) {
				case networkStateFailed:
					return
				case networkStateConnecting:
					return
				}

				// Acknowledge
				var err error
				if ack {
					err = msg.Ack(false)
				} else {
					err = msg.Nack(false, true)
				}

				if err != nil {
					if isTransientError(err) {
						if r.reconnectRBMQ() {
							return
						}
					}
					log.Printf("Failed to acknowledge message: %v", err)
					return
				}
			}(index, &d)
			index++
			if index == r.consumerWorkerCount {
				if recovery {
					wg.Wait()
					recovery = false
				}
				index = 0
			}
		}
		wg.Wait()
	}
}

func (r *RabbitMQRepository) declareQueue() (*amqp091.Channel, error) {
	delay := 4 * time.Second
	for {
		ch, err := r.Conn.Channel()
		if err != nil {
			return nil, fmt.Errorf("failed to create channel: %w", err)
		}
		_, err = ch.QueueDeclare(
			r.dsnRBQ.Queue,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			ch.Close()
			return nil, fmt.Errorf("failed to declare queue: %w", err)
		}

		err = ch.QueueBind(
			r.dsnRBQ.Queue,
			"",
			r.dsnRBQ.Exchange,
			false,
			nil,
		)
		if err != nil {
			var amqpErr *amqp091.Error
			if errors.As(err, &amqpErr) && amqpErr.Code == 404 {
				log.Printf("Failed to bind with exchange: %v", err)
				time.Sleep(delay)
				continue
			}
			ch.Close()
			return nil, fmt.Errorf("failed to bind queue: %w", err)
		}

		err = ch.Qos(r.consumerWorkerCount, 0, true)
		if err != nil {
			ch.Close()
			return nil, fmt.Errorf("failed to set queue qos: %w", err)
		}
		return ch, nil
	}
}
