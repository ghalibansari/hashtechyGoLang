package src

import (
	_ "hashtechy/docs" // Add this line to import Swagger docs
	"hashtechy/src/logger"
	"hashtechy/src/postgres"
	"hashtechy/src/server"
	"net/http"
	"sync"
	"time"
)

func App() error {
	defer logger.Close()
	logger.Info("Starting application...")

	err := postgres.Connect()
	if err != nil {
		logger.Error("Failed to connect to database: %v", err)
		return err
	}
	defer postgres.DB.Close()

	postgres.DropUserTable()
	postgres.CreateUserTable()
	if err := postgres.CreateIndexes(); err != nil {
		logger.Error("Failed to create indexes: %v", err)
		return err
	}
	postgres.ShowDatabases()
	postgres.ShowTables()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		if err := consumer(); err != nil {
			logger.Error("Consumer error: %v", err)
		}
		defer wg.Done()
	}()

	wg.Add(1)
	go func() {
		time.Sleep(time.Second * 3)
		if err := Producer(); err != nil {
			logger.Error("Producer error: %v", err)
		}
		defer wg.Done()
	}()

	mux := server.Server()
	server.AddSwaggerHandler(mux) // swagger

	logger.Info("Server starting on port 3000...")
	if err := http.ListenAndServe(":3000", mux); err != nil {
		logger.Error("Server error: %v", err)
		return err
	}

	wg.Wait()
	return nil
}
