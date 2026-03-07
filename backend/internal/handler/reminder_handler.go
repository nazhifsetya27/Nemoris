package handler

import (
	"encoding/json"
	"net/http"

	"nemoris/internal/repository"
)

func GetReminders(w http.ResponseWriter, r *http.Request) {
	reminders, err := repository.GetAllReminders()
	if err != nil {
		http.Error(w, "failed to fetch reminders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reminders)
}

func GetPendingReminders(w http.ResponseWriter, r *http.Request) {
	reminders, err := repository.GetPendingReminders()
	if err != nil {
		http.Error(w, "failed to fetch pending reminders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reminders)
}