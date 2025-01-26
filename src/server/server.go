package server

import (
	"encoding/json"
	"hashtechy/src/postgres"
	"net/http"
	"strconv"

	"golang.org/x/time/rate"
)

func Server() *http.ServeMux {
	mux := http.NewServeMux()
	limiter := NewRateLimiter(rate.Limit(2), 2)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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
			http.Error(w, "Failed to get users", http.StatusInternalServerError)
			return
		}

		jsonData, err := json.Marshal(users)
		if err != nil {
			http.Error(w, "Failed to encode data", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})

	mux.Handle("/users", limiter.RateLimit(handler))

	// TODO: temp remove this later, just for testing
	getAllDataHandler := getAllData()
	mux.Handle("/all", getAllDataHandler)

	return mux
}
