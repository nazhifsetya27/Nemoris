package whatsapp

import (
	"encoding/json"
	"strings"

	"nemoris/internal/model"
)

// ExtractMessageFromWAHA parses WAHA webhook envelope and returns a normalized message.
// Only "message" events with non-fromMe payload are returned; others return empty ok=false.
// Falls back to flat {from, body} for test scripts.
func ExtractMessageFromWAHA(raw []byte) (model.IncomingMessage, bool) {
	var envelope model.WAHAWebhookEnvelope
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return model.IncomingMessage{}, false
	}
	if envelope.Event != "" && envelope.Event != "message" {
		return model.IncomingMessage{}, false
	}
	// WAHA format: payload has from, body
	if len(envelope.Payload) > 0 {
		var p model.WAHAMessagePayload
		if err := json.Unmarshal(envelope.Payload, &p); err == nil {
			if p.FromMe {
				return model.IncomingMessage{}, false
			}
			return model.IncomingMessage{
				From: strings.TrimSpace(p.From),
				Body: strings.TrimSpace(p.Body),
			}, true
		}
	}
	// Fallback: flat {from, body} for test scripts
	var flat model.IncomingMessage
	if err := json.Unmarshal(raw, &flat); err != nil {
		return model.IncomingMessage{}, false
	}
	if flat.From == "" && flat.Body == "" {
		return model.IncomingMessage{}, false
	}
	return model.IncomingMessage{
		From: strings.TrimSpace(flat.From),
		Body: strings.TrimSpace(flat.Body),
	}, true
}

func NormalizePayload(raw model.IncomingMessage) model.IncomingMessage {
	return model.IncomingMessage{
		From: raw.From,
		Body: raw.Body,
	}
}