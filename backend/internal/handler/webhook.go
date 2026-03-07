package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"nemoris/internal/model"
	"nemoris/internal/service"
)

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	var msg model.IncomingMessage

	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	log.Println("Incoming message:")
	log.Printf("From: %s | Body: %s\n", msg.From, msg.Body)

	response := service.ProcessMessage(msg.From, msg.Body)

	w.Write([]byte(response))
}