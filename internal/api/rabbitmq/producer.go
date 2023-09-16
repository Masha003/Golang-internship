package rabbitmq

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/Masha003/Golang-internship/internal/config"
	"github.com/Masha003/Golang-internship/internal/models"
	"github.com/wagslane/go-rabbitmq"
)

type Producer interface {
	SendUser(models.User)
	Close() error
}

func NewProducer(cfg config.Config) Producer {
	log.Print("Initializing mail service")

	const retries = 5
	var conn *rabbitmq.Conn
	var err error

	for i := 0; i < retries; i++ {
		conn, err = rabbitmq.NewConn(
			cfg.RabbitMQUrl,
			rabbitmq.WithConnectionOptionsLogging,
		)
		if err != nil {
			log.Print("Failed to connect to RabbitMQ, retrying in 3 seconds...")
			time.Sleep(3 * time.Second)
		} else {
			break
		}
	}
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ")
	}

	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
	)
	if err != nil {
		log.Fatal("Failed to create publisher")
	}

	return &producer{
		conn:      conn,
		publisher: publisher,
	}
}

type producer struct {
	conn      *rabbitmq.Conn
	publisher *rabbitmq.Publisher
}

func (s *producer) Close() error {
	s.publisher.Close()
	return s.conn.Close()
}

func (p *producer) SendUser(user models.User) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		body, err := json.Marshal(user)
		if err != nil {
			log.Print("Failed to marshal message")
			return
		}

		err = p.send(ctx, body)
		if err != nil {
			log.Print("Failed to send message")
		}
	}()
}

func (p *producer) send(ctx context.Context, body []byte) error {

	return p.publisher.Publish(
		body,
		[]string{queueName},
		rabbitmq.WithPublishOptionsContentType("application/json"),
	)
}
