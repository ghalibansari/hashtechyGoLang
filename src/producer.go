package src

import (
	"encoding/json"
	"fmt"
	"hashtechy/src/rabbitmq"
	"log"
	"strings"
	"sync"
)

func producer() error {
	err, header, csvChannel := readCsv("./user.csv")
	if err != nil {
		return fmt.Errorf("failed to read CSV: %w", err)
	}

	fmt.Println("Header:", strings.Join(header, ","))

	conn, ch, publisher, _, err := rabbitmq.ConnectToRabbitMQ("csv_queue")
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	defer conn.Close()
	defer ch.Close()

	var wg sync.WaitGroup

	for msg := range csvChannel {
		body, err := json.Marshal(msg)
		if err != nil {
			log.Printf("failed to marshal message to JSON: %v", err)
			continue
		}

		// test scalability
		// sem := make(chan struct{}, 1_000) // Create a semaphore with a limit of 1000
		// for i := 0; i < 10_00_000; i++ {
		// 	sem <- struct{}{} // Acquire a semaphore
		// 	wg.Add(1)
		// 	go func(index int) {
		// 		defer wg.Done()
		// 		defer func() { <-sem }() // Release the semaphore
		// 		err := publisher([]byte(body))
		// 		if err != nil {
		// 			log.Printf("failed to publish message: %v", err)
		// 		} else {
		// 			log.Printf("sent message %d: %s", index, body)
		// 		}
		// 	}(i)
		// }

		wg.Add(1)
		go func() {
			defer wg.Done()
			err = publisher([]byte(body))
			if err != nil {
				log.Printf("failed to publish message: %v", err)
			} else {
				log.Printf("sent message: %s", body)
			}
		}()

	}

	wg.Wait()
	return nil
}
