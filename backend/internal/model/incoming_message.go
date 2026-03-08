package model

import "encoding/json"

// IncomingMessage is the normalized message used by handlers.
type IncomingMessage struct {
	From string `json:"from"`
	Body string `json:"body"`
}

// WAHAWebhookEnvelope is the raw WAHA webhook payload structure.
// WAHA sends { event, session, payload: { from, body, ... } }.
type WAHAWebhookEnvelope struct {
	Event   string          `json:"event"`
	Session string          `json:"session"`
	Payload json.RawMessage `json:"payload"`
}

// WAHAMessagePayload is the payload for event "message".
type WAHAMessagePayload struct {
	From   string `json:"from"`
	Body   string `json:"body"`
	FromMe bool   `json:"fromMe"`
}