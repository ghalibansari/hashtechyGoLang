package src

import (
	"hashtechy/src/postgres"
	"hashtechy/src/server"
	"log"
	"net/http"
	"sync"
)

func App() error {

	err := postgres.Connect()
	if err != nil {
		return err
	}
	defer postgres.DB.Close()

	postgres.DropUserTable()
	postgres.CreateUserTable()
	postgres.ShowDatabases()
	postgres.ShowTables()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		consumer()
		defer wg.Done()
	}()

	err = producer()
	if err != nil {
		return err
	}

	mux := server.Server()
	log.Println("Server starting on port 3000...")
	http.ListenAndServe(":3000", mux)

	wg.Wait()
	return nil
}
