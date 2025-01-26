package redis

import (
	"context"
	"crypto/tls"
	"hashtechy/src/errors"
	"hashtechy/src/logger"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func ConnectToRedis() (func(string, string) error, func(string) (string, error)) {
	if err := godotenv.Load(); err != nil {
		logger.Debug("No .env file found: %v", err)
	}

	// Load TLS certificates
	cert, err := tls.LoadX509KeyPair("/app/certs/redis.crt", "/app/certs/redis.key")
	if err != nil {
		logger.Error("Failed to load Redis certificates: %v", err)
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
		logger.Error("Failed to connect to Redis: %v", err)
		return nil, nil
	}

	logger.Info("Successfully connected to Redis")

	setValue := func(key string, value string) error {
		if err := rdb.Set(ctx, key, value, 5*time.Minute).Err(); err != nil {
			logger.Error("Failed to set value in Redis: %v", err)
			return errors.New(errors.ErrDatabase, "failed to set value in Redis", err)
		}
		logger.Debug("Successfully set value for key: %s", key)
		return nil
	}

	getValue := func(key string) (string, error) {
		value, err := rdb.Get(ctx, key).Result()
		if err != nil {
			logger.Error("Failed to get value from Redis: %v", err)
			return "", errors.New(errors.ErrDatabase, "failed to get value from Redis", err)
		}
		logger.Debug("Successfully retrieved value for key: %s", key)
		return value, nil
	}

	return setValue, getValue
}
