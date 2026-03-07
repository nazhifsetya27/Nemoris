package handler

import (
	"encoding/json"
	"net/http"

	"nemoris/internal/model"
	"nemoris/internal/service"
	"nemoris/internal/whatsapp"
	"nemoris/internal/utils"
)

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	var raw model.IncomingMessage

	err := json.NewDecoder(r.Body).Decode(&raw)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	msg := whatsapp.NormalizePayload(raw)

	utils.LogInbound("From: " + msg.From + " | Body: " + msg.Body)

	response := service.ProcessMessage(msg.From, msg.Body)

	if response != "" {
		result := whatsapp.SendText(msg.From, response)
		if result.Err != nil {
			utils.LogInbound("send failed: " + result.Err.Error())
		}
	}

	w.Write([]byte("ok"))
}