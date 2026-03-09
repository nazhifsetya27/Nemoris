package handler

import (
	"encoding/json"
	"net/http"

	"nemoris/internal/service"
)

func GetFailedReminders(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	status := r.URL.Query().Get("status")
	failureType := r.URL.Query().Get("failure_type")
	date := r.URL.Query().Get("date")

	reminders, err := service.ListFailedReminders(from, status, failureType, date)
	if err != nil {
		http.Error(w, "failed to fetch failed reminders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reminders)
}
