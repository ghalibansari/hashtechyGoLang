package server

import (
	"encoding/json"
	"hashtechy/src/logger"
	"hashtechy/src/postgres"
	"net/http"
	"time"
)

func handleSSE(w http.ResponseWriter, r *http.Request) {
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create channel for client disconnect detection
	notify := w.(http.CloseNotifier).CloseNotify()

	for {
		select {
		case <-notify:
			logger.Info("Client disconnected")
			return
		default:
			// Get latest data
			users, err := postgres.GetAllUsers()
			if err != nil {
				logger.Error("SSE: Failed to get users: %v", err)
				continue
			}

			// Create event data
			data := map[string]interface{}{
				"timestamp": time.Now().Format(time.RFC3339),
				"users":     users,
			}

			jsonData, err := json.Marshal(data)
			if err != nil {
				logger.Error("SSE: Failed to marshal data: %v", err)
				continue
			}

			// Send event
			_, err = w.Write([]byte("data: " + string(jsonData) + "\n\n"))
			if err != nil {
				logger.Error("SSE: Failed to write data: %v", err)
				return
			}

			w.(http.Flusher).Flush()
			time.Sleep(5 * time.Second) // Update every 5 seconds
		}
	}
}
