package src

import (
	"encoding/json"
	"fmt"
	"hashtechy/src/postgres"
	"hashtechy/src/rabbitmq"
	"hashtechy/src/redis"
	"hashtechy/src/user"
	"log"
	"runtime"
	"sync"
)

func consumer() error {
	conn, ch, _, consumer, err := rabbitmq.ConnectToRabbitMQ("csv_queue")
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	defer conn.Close()
	defer ch.Close()

	setValue, _ := redis.ConnectToRedis()

	var wg sync.WaitGroup
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	// defer cancel()

	numCPU := runtime.NumCPU()
	for i := 0; i < numCPU*numCPU; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case msg := <-consumer:
					log.Printf("received message: %s", msg.Body)
					userData := string(msg.Body)
					var user user.User
					if err := json.Unmarshal([]byte(userData), &user); err != nil {
						log.Printf("failed to unmarshal user data: %v", err)
					}

					user, err = postgres.InsertUser(user)
					if err != nil {
						log.Printf("panic: %v", err)
					}

					userJson, _ := json.Marshal(user)
					err = setValue("users/"+user.ID, string(userJson)) // Changed msg.body to msg.Body
					if err != nil {
						log.Printf("panic: %v", err)
					} else {
						log.Printf("value set in redis: %s", string(userJson))
					}

					if err := msg.Ack(false); err != nil {
						log.Printf("Failed to acknowledge message: %s", err)
					} else {
						log.Println("Message acknowledged")
					}
					// case <-ctx.Done():
					// 	log.Println("No messages received for 10 seconds, closing...")
					// 	return
				}
			}
		}()
	}

	wg.Wait()
	return nil
}
