package src

import (
	"encoding/json"
	"hashtechy/src/logger"
	"hashtechy/src/rabbitmq"
	"strings"
	"sync"
)

func producer() error {
	err, header, csvChannel := readCsv("./user.csv")
	if err != nil {
		logger.Error("Failed to read CSV: %v", err)
		return err
	}

	logger.Info("Header: %s", strings.Join(header, ","))

	conn, ch, publisher, _, err := rabbitmq.ConnectToRabbitMQ("csv_queue")
	if err != nil {
		logger.Error("Failed to connect to RabbitMQ: %v", err)
		return err
	}
	defer conn.Close()
	defer ch.Close()

	var wg sync.WaitGroup

	for msg := range csvChannel {
		body, err := json.Marshal(msg)
		if err != nil {
			logger.Error("Failed to marshal message to JSON: %v", err)
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			err = publisher([]byte(body))
			if err != nil {
				logger.Error("Failed to publish message: %v", err)
			} else {
				logger.Info("Sent message: %s", body)
			}
		}()
	}

	wg.Wait()
	return nil
}
