package whatsapp

import (
	"strings"

	"nemoris/internal/config"
	"nemoris/internal/utils"
)

func SendText(to string, text string) bool {
	if strings.TrimSpace(text) == "" {
		utils.LogOutbound("Blocked empty message")
		return false
	}

	if to == config.App.BotNumber {
		utils.LogOutbound("Blocked self message")
		return false
	}

	utils.LogOutbound("To: " + to + " | Text: " + text)

	// WAHA HTTP send will replace this later

	utils.LogOutbound("WAHA response: simulated success")

	return true
}