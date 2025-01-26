package server

import (
	"encoding/json"
	"hashtechy/src/logger"
	"hashtechy/src/postgres"
	"net/http"
	"strconv"

	"golang.org/x/time/rate"
)

// @Summary Get users with filtering
// @Description Get users with optional filtering by name and age
// @Tags users
// @Accept json
// @Produce json
// @Param name query string false "Filter by name"
// @Param min_age query integer false "Minimum age"
// @Param max_age query integer false "Maximum age"
// @Param limit query integer false "Number of records to return"
// @Param skip query integer false "Number of records to skip"
// @Success 200 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 429 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users [get]
func Server() *http.ServeMux {
	mux := http.NewServeMux()
	limiter := NewRateLimiter(rate.Limit(2), 2)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		query := r.URL.Query()
		name := query.Get("name")
		min_age, _ := strconv.Atoi(query.Get("min_age"))
		max_age, _ := strconv.Atoi(query.Get("max_age"))
		limit, _ := strconv.Atoi(query.Get("limit"))
		skip, _ := strconv.Atoi(query.Get("skip"))

		users, err := postgres.GetAllUsersByQuery(name, min_age, max_age, limit, skip)
		if err != nil {
			logger.Error("Failed to get users: %v", err)
			sendError(w, "Failed to get users", http.StatusInternalServerError)
			return
		}

		sendResponse(w, UserResponse{
			Status: "success",
			Data:   users,
		})
	})

	mux.Handle("/users", limiter.RateLimit(handler))
	mux.Handle("/all", getAllData())

	// Add SSE endpoint
	mux.HandleFunc("/events", handleSSE)

	return mux
}

func sendResponse(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(response)
	if err != nil {
		logger.Error("Failed to marshal response: %v", err)
		sendError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func sendError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := ErrorResponse{
		Status:  "error",
		Message: message,
	}
	json.NewEncoder(w).Encode(response)
}
