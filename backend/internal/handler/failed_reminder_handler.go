package handler

import (
	"encoding/json"
	"net/http"

	"nemoris/internal/service"
)

func GetFailedReminders(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")

	reminders, err := service.ListFailedReminders(from)
	if err != nil {
		http.Error(w, "failed to fetch failed reminders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reminders)
}
