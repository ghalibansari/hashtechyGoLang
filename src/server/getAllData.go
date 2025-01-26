package server

import (
	"encoding/json"
	"hashtechy/src/postgres"
	"hashtechy/src/redis"
	"net/http"
)

func getAllData() http.HandlerFunc {
	_, getValue := redis.ConnectToRedis()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		users, err := postgres.GetAllUsers()
		if err != nil {
			http.Error(w, "Failed to get users", http.StatusInternalServerError)
			return
		}

		redisMap := make(map[string]string)

		for _, user := range users {
			value, err := getValue("users/" + user.ID)
			if err != nil {
				continue
			}
			redisMap[user.ID] = value
		}

		// Create combined response map
		combinedData := map[string]interface{}{
			"users":     users,
			"redisData": redisMap,
		}

		// Convert combined data to JSON
		jsonData, err := json.Marshal(combinedData)
		if err != nil {
			http.Error(w, "Failed to encode combined data", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})

	return handler
}
