package whatsapp

import "nemoris/internal/model"

func NormalizePayload(raw model.IncomingMessage) model.IncomingMessage {
	return model.IncomingMessage{
		From: raw.From,
		Body: raw.Body,
	}
}