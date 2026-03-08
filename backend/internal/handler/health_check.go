package handler

import (
	"encoding/json"
	"net/http"

	"nemoris/internal/service"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	result := service.CheckHealth()

	statusCode := http.StatusOK
	if result.Status != "ok" {
		statusCode = http.StatusServiceUnavailable
	}
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(result)
}
