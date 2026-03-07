package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"nemoris/internal/config"
	"nemoris/internal/utils"
)

var (
	duplicateMu    sync.Mutex
	duplicateStore = make(map[string]time.Time)
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

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
	// Protection 3: safe empty normalization
	text = strings.TrimSpace(text)
	if text == "" {
		utils.LogOutbound("Blocked empty message")
		return SendResult{Accepted: false, Err: nil}
	}

	if to == config.App.BotNumber {
		utils.LogOutbound("Blocked self message")
		return SendResult{Accepted: false, Err: nil}
	}

	// Protection 2: duplicate outbound suppression
	key := to + "|" + text
	duplicateBlocked := false
	func() {
		duplicateMu.Lock()
		defer duplicateMu.Unlock()
		if t, ok := duplicateStore[key]; ok && time.Since(t) < 10*time.Second {
			utils.LogOutbound("duplicate blocked")
			duplicateBlocked = true
			return
		}
		for k, v := range duplicateStore {
			if time.Since(v) > 10*time.Second {
				delete(duplicateStore, k)
			}
		}
	}()
	if duplicateBlocked {
		return SendResult{Accepted: false, Err: nil}
	}

	// Protection 1: outbound pacing (300ms + random 0..600ms)
	time.Sleep(300*time.Millisecond + time.Duration(rand.Intn(601))*time.Millisecond)

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

	// Record duplicate only after successful HTTP send
	duplicateMu.Lock()
	defer duplicateMu.Unlock()
	now := time.Now()
	for k, v := range duplicateStore {
		if time.Since(v) > 10*time.Second {
			delete(duplicateStore, k)
		}
	}
	duplicateStore[key] = now

	utils.LogOutbound("WAHA response: " + string(respBody))
	return SendResult{Accepted: true, Err: nil}
}
