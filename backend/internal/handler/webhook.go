package handler

import (
	"io"
	"net/http"

	"nemoris/internal/middleware"
	"nemoris/internal/service"
	"nemoris/internal/whatsapp"
	"nemoris/internal/utils"
)

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	msg, ok := whatsapp.ExtractMessageFromWAHA(body)
	if !ok {
		w.Write([]byte("ok"))
		return
	}

	// Resolve LID to phone for storage and sending
	msg.From = whatsapp.ResolveSendTarget(msg.From)

	requestID := middleware.GetRequestID(r.Context())
	utils.LogInboundWithRequestID(requestID, "From: "+msg.From+" | Body: "+msg.Body)

	response := service.ProcessMessage(msg.From, msg.Body)

	if response != "" && msg.From != "" {
		result := whatsapp.SendText(msg.From, response)
		if result.Err != nil {
			utils.LogInboundWithRequestID(requestID, "send failed: "+result.Err.Error())
		}
	}

	w.Write([]byte("ok"))
}