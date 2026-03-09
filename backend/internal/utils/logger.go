package utils

import "log"

func LogInbound(message string) {
	log.Println("[INBOUND]", message)
}

// LogInboundWithRequestID logs inbound events with optional request correlation id.
// When requestID is non-empty, appends "request_id=<id>" for traceability.
// Use when handler has request id from context; existing LogInbound remains for callers without it.
func LogInboundWithRequestID(requestID, message string) {
	if requestID != "" {
		log.Println("[INBOUND]", message, "request_id="+requestID)
	} else {
		log.Println("[INBOUND]", message)
	}
}

func LogOutbound(message string) {
	log.Println("[OUTBOUND]", message)
}

func LogDB(message string) {
	log.Println("[DB]", message)
}

func LogScheduler(message string) {
	log.Println("[SCHEDULER]", message)
}

func LogAI(message string) {
	log.Println("[AI]", message)
}

// LogSecurity logs security events (category: inbound/security).
// Never log secrets—only authorized/unauthorized outcomes.
func LogSecurity(message string) {
	log.Println("[INBOUND/SECURITY]", message)
}

// LogSystem logs system-level events (e.g. panic recovery).
// Category: [SYSTEM]
func LogSystem(message string) {
	log.Println("[SYSTEM]", message)
}