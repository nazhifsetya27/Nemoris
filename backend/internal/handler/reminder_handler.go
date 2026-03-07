package handler

import (
	"encoding/json"
	"net/http"

	"nemoris/internal/service"
)

func GetReminders(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")

	reminders, err := service.ListAllReminders(from)
	if err != nil {
		http.Error(w, "failed to fetch reminders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reminders)
}

func GetPendingReminders(w http.ResponseWriter, r *http.Request) {
	reminders, err := service.ListPendingReminders()
	if err != nil {
		http.Error(w, "failed to fetch pending reminders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reminders)
}