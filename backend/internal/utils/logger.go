package utils

import "log"

func LogInbound(message string) {
	log.Println("[INBOUND]", message)
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