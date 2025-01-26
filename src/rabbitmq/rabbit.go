package rabbitmq

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectToRabbitMQ(name string) (*amqp.Connection, *amqp.Channel, func([]byte) error, <-chan amqp.Delivery, error) {
	godotenv.Load()
	amqpURL := os.Getenv("AMQP_URL")
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@rabbitmq:5672/"
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
		log.Printf("Attempt %d: Failed to connect to RabbitMQ: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to connect to RabbitMQ after retries: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, nil, nil, fmt.Errorf("failed to open a channel: %w", err)
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
		return nil, nil, nil, nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	// Define the publisher function
	publisher := func(body []byte) error {
		return ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{ContentType: "text/plain", Body: body},
		)
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
		return nil, nil, nil, nil, fmt.Errorf("failed to register a consumer: %w", err)
	}

	return conn, ch, publisher, consumer, nil
}
