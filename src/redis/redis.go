package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func ConnectToRedis() (func(string, string) error, func(string) (string, error)) {
	godotenv.Load()

	// Load TLS certificates
	cert, err := tls.LoadX509KeyPair("/app/certs/redis.crt", "/app/certs/redis.key")
	if err != nil {
		log.Printf("Error loading Redis certificates: %v", err)
		return nil, nil
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
		TLSConfig: &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: true,
			Certificates:       []tls.Certificate{cert},
			RootCAs:            nil,
		},
	})

	// Test the connection
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("Error connecting to Redis: %v", err)
		return nil, nil
	}

	fmt.Println("Connected to Redis")

	setValue := func(key string, value string) error {
		return rdb.Set(ctx, key, value, 5*time.Minute).Err()
	}

	getValue := func(key string) (string, error) {
		return rdb.Get(ctx, key).Result()
	}

	return setValue, getValue
}
