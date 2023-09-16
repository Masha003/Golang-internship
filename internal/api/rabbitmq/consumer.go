package rabbitmq

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Masha003/Golang-internship/internal/config"
	"github.com/Masha003/Golang-internship/internal/models"
	"github.com/Masha003/Golang-internship/internal/service"
	"github.com/wagslane/go-rabbitmq"
)

type Consumer interface {
	Close() error
}

type consumer struct {
	conn        *rabbitmq.Conn
	consumer    *rabbitmq.Consumer
	userService service.UserService
}

const queueName = "users"

func NewConsumer(cfg config.Config, userService service.UserService) Consumer {
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

	c := &consumer{
		conn:        conn,
		userService: userService,
	}

	cons, err := rabbitmq.NewConsumer(
		conn,
		c.handleDelivery,
		queueName,
		rabbitmq.WithConsumerOptionsQueueDurable,
	)
	if err != nil {
		log.Fatal("Failed to create consumer")
	}

	c.consumer = cons

	return c
}

func (c *consumer) Close() error {
	c.consumer.Close()
	return c.conn.Close()
}

func (c *consumer) handleDelivery(d rabbitmq.Delivery) rabbitmq.Action {
	log.Print("Processing delivery")

	var user models.RegisterUser

	err := json.Unmarshal(d.Body, &user)
	if err != nil {
		log.Print("Failed to unmarshal message")
		return rabbitmq.Ack
	}

	c.userService.Register(user)

	return rabbitmq.Ack
}
