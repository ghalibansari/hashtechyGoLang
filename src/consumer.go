package src

import (
	"encoding/json"
	"hashtechy/src/errors"
	"hashtechy/src/logger"
	"hashtechy/src/postgres"
	"hashtechy/src/rabbitmq"
	"hashtechy/src/redis"
	"hashtechy/src/user"
	"runtime"
	"sync"
)

func consumer() error {
	conn, ch, _, consumer, err := rabbitmq.ConnectToRabbitMQ("csv_queue")
	if err != nil {
		logger.Error("Failed to connect to RabbitMQ: %v", err)
		return errors.New(errors.ErrNetwork, "failed to connect to RabbitMQ", err)
	}
	defer conn.Close()
	defer ch.Close()

	setValue, _ := redis.ConnectToRedis()

	var wg sync.WaitGroup
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	// defer cancel()

	numCPU := runtime.NumCPU()
	logger.Info("Starting %d consumer goroutines", numCPU*numCPU)

	for i := 0; i < numCPU*numCPU; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case msg := <-consumer:
					logger.Debug("Received message: %s", msg.Body)
					userData := string(msg.Body)
					var user user.User
					if err := json.Unmarshal([]byte(userData), &user); err != nil {
						logger.Error("Failed to unmarshal user data: %v", err)
						continue
					}

					user, err = postgres.InsertUser(user)
					if err != nil {
						logger.Error("Failed to insert user: %v", err)
						continue
					}

					userJson, err := json.Marshal(user)
					if err != nil {
						logger.Error("Failed to marshal user data: %v", err)
						continue
					}

					err = setValue("users/"+user.ID, string(userJson))
					if err != nil {
						logger.Error("Failed to set value in Redis: %v", err)
					} else {
						logger.Debug("Successfully set value in Redis: %s", string(userJson))
					}

					if err := msg.Ack(false); err != nil {
						logger.Error("Failed to acknowledge message: %v", err)
					} else {
						logger.Debug("Message acknowledged successfully")
					}
					// case <-ctx.Done():
					// 	logger.Info("No messages received for 10 seconds, closing...")
					// 	return
				}
			}
		}()
	}

	wg.Wait()
	return nil
}
