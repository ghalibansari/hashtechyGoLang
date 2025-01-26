package server

import (
	"encoding/json"
	"fmt"
	"hashtechy/src/logger"
	"hashtechy/src/postgres"
	"net/http"
	"time"
)

func handleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	done := r.Context().Done()
	ticker := time.NewTicker(3 * time.Second) // update every 3 seconds
	defer ticker.Stop()

	// Send initial data immediately
	sendUpdate(w)

	for {
		select {
		case <-done:
			logger.Info("Client disconnected")
			return
		case <-ticker.C:
			if err := sendUpdate(w); err != nil {
				logger.Error("SSE: Failed to send update: %v", err)
				return
			}
		}
	}
}

func sendUpdate(w http.ResponseWriter) error {
	users, err := postgres.GetAllUsers()
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"users":     users,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Write event, data, and newlines in ONE operation
	_, err = fmt.Fprintf(w, "event: userUpdate\ndata: %s\n\n", string(jsonData))
	if err != nil {
		return err
	}

	// Flush immediately
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	return nil
}
