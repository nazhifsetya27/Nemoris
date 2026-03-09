package handler

import (
	"encoding/json"
	"net/http"

	"nemoris/internal/service"
)

func GetMemories(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	if from == "" {
		http.Error(w, "from is required", http.StatusBadRequest)
		return
	}

	memories, err := service.ListMemories(from)
	if err != nil {
		http.Error(w, "failed to fetch memories", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(memories)
}
