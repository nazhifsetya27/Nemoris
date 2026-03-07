package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"nemoris/internal/config"
	"nemoris/internal/utils"
)

type sendTextPayload struct {
	Session string `json:"session"`
	ChatID  string `json:"chatId"`
	Text    string `json:"text"`
}

// SendResult holds the outcome of a send attempt.
type SendResult struct {
	Accepted bool
	Err      error
}

func SendText(to string, text string) SendResult {
	if strings.TrimSpace(text) == "" {
		utils.LogOutbound("Blocked empty message")
		return SendResult{Accepted: false, Err: nil}
	}

	if to == config.App.BotNumber {
		utils.LogOutbound("Blocked self message")
		return SendResult{Accepted: false, Err: nil}
	}

	// EPIC 4: Temporary simulated failure for retry logic testing — set to false for production
	// if true {
	// 	return SendResult{Accepted: false, Err: errors.New("simulated failure")}
	// }

	utils.LogOutbound("To: " + to + " | Text: " + text)

	chatID := to
	if !strings.Contains(chatID, "@") {
		chatID = chatID + "@c.us"
	}

	payload := sendTextPayload{
		Session: "default",
		ChatID:  chatID,
		Text:    text,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return SendResult{Accepted: false, Err: fmt.Errorf("marshal payload: %w", err)}
	}

	url := strings.TrimSuffix(config.App.WAHABaseURL, "/") + "/api/sendText"
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return SendResult{Accepted: false, Err: fmt.Errorf("create request: %w", err)}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if config.App.WAHAAPIKey != "" {
		req.Header.Set("X-Api-Key", config.App.WAHAAPIKey)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return SendResult{Accepted: false, Err: fmt.Errorf("WAHA request failed: %w", err)}
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return SendResult{Accepted: false, Err: fmt.Errorf("WAHA sendText failed (status %d): %s", resp.StatusCode, string(respBody))}
	}

	utils.LogOutbound("WAHA response: " + string(respBody))
	return SendResult{Accepted: true, Err: nil}
}
