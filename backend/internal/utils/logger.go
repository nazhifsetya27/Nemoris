package utils

import (
	"log"
	"strings"
)

const (
	sevInfo    = "INFO"
	sevWarning = "WARNING"
	sevError   = "ERROR"
)

func logEvent(severity, category, message string, extra ...string) {
	parts := []string{"[" + severity + "]", "[" + category + "]", message}
	parts = append(parts, extra...)
	log.Println(strings.Join(parts, " "))
}

func LogInbound(message string) {
	logEvent(sevInfo, "INBOUND", message)
}

func LogInboundWarn(message string) {
	logEvent(sevWarning, "INBOUND", message)
}

func LogInboundError(message string) {
	logEvent(sevError, "INBOUND", message)
}

// LogInboundWithRequestID logs inbound events with optional request correlation id.
// When requestID is non-empty, appends "request_id=<id>" for traceability.
// Use when handler has request id from context; existing LogInbound remains for callers without it.
func LogInboundWithRequestID(requestID, message string) {
	if requestID != "" {
		logEvent(sevInfo, "INBOUND", message, "request_id="+requestID)
	} else {
		logEvent(sevInfo, "INBOUND", message)
	}
}

func LogInboundWithRequestIDError(requestID, message string) {
	if requestID != "" {
		logEvent(sevError, "INBOUND", message, "request_id="+requestID)
	} else {
		logEvent(sevError, "INBOUND", message)
	}
}

func LogOutbound(message string) {
	logEvent(sevInfo, "OUTBOUND", message)
}

func LogOutboundWarn(message string) {
	logEvent(sevWarning, "OUTBOUND", message)
}

func LogDB(message string) {
	logEvent(sevInfo, "DB", message)
}

func LogDBError(message string) {
	logEvent(sevError, "DB", message)
}

func LogScheduler(message string) {
	logEvent(sevInfo, "SCHEDULER", message)
}

func LogSchedulerWarn(message string) {
	logEvent(sevWarning, "SCHEDULER", message)
}

func LogSchedulerError(message string) {
	logEvent(sevError, "SCHEDULER", message)
}

func LogAI(message string) {
	logEvent(sevInfo, "AI", message)
}

func LogAIWarn(message string) {
	logEvent(sevWarning, "AI", message)
}

// LogSecurity logs security events (category: inbound/security).
// Never log secrets—only authorized/unauthorized outcomes.
func LogSecurity(message string) {
	logEvent(sevInfo, "INBOUND/SECURITY", message)
}

func LogSecurityWarn(message string) {
	logEvent(sevWarning, "INBOUND/SECURITY", message)
}

// LogSystem logs system-level events (e.g. panic recovery).
// Category: [SYSTEM]
func LogSystem(message string) {
	logEvent(sevInfo, "SYSTEM", message)
}

func LogSystemError(message string) {
	logEvent(sevError, "SYSTEM", message)
}