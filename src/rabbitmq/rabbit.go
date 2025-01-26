package rabbitmq

import (
	"crypto/tls"
	"hashtechy/src/errors"
	"hashtechy/src/logger"
	"os"
	"time"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectToRabbitMQ(name string) (*amqp.Connection, *amqp.Channel, func([]byte) error, <-chan amqp.Delivery, error) {
	if err := godotenv.Load(); err != nil {
		logger.Debug("No .env file found: %v", err)
	}

	amqpURL := os.Getenv("AMQP_URL")
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@rabbitmq:5672/"
		logger.Debug("Using default RabbitMQ URL: %s", amqpURL)
	}

	config := &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS12,
		RootCAs:            nil,
		Certificates:       nil,
	}

	var conn *amqp.Connection
	var err error
	for i := 0; i < 10; i++ {
		conn, err = amqp.DialTLS(amqpURL, config)
		if err == nil {
			break
		}
		logger.Error("Attempt %d: Failed to connect to RabbitMQ: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return nil, nil, nil, nil, errors.New(errors.ErrNetwork, "failed to connect to RabbitMQ after retries", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, nil, nil, errors.New(errors.ErrNetwork, "failed to open a channel", err)
	}

	q, err := ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, nil, nil, nil, errors.New(errors.ErrNetwork, "failed to declare a queue", err)
	}

	// Define the publisher function
	publisher := func(body []byte) error {
		if err := ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{ContentType: "text/plain", Body: body},
		); err != nil {
			logger.Error("Failed to publish message: %v", err)
			return errors.New(errors.ErrNetwork, "failed to publish message", err)
		}
		logger.Debug("Successfully published message to queue: %s", q.Name)
		return nil
	}

	// Create a consumer channel
	consumer, err := ch.Consume(
		name,  // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, nil, nil, nil, errors.New(errors.ErrNetwork, "failed to register a consumer", err)
	}

	logger.Info("Successfully connected to RabbitMQ and set up queue: %s", name)
	return conn, ch, publisher, consumer, nil
}
