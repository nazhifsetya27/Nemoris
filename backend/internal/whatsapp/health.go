package whatsapp

import (
	"net/http"
	"strings"
	"time"

	"nemoris/internal/config"
)

var healthClient = &http.Client{Timeout: 3 * time.Second}

// CheckWAHA performs a lightweight GET to WAHA /api/sessions.
// Returns true if HTTP 200, false otherwise.
// Does not panic if WAHA is unavailable.
func CheckWAHA() bool {
	url := strings.TrimSuffix(config.App.WAHABaseURL, "/") + "/api/sessions"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false
	}
	req.Header.Set("Accept", "application/json")
	if config.App.WAHAAPIKey != "" {
		req.Header.Set("X-Api-Key", config.App.WAHAAPIKey)
	}

	resp, err := healthClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// CheckWAHALatency performs a lightweight GET to WAHA /api/sessions.
// Returns (ok, latency). Use for health checks that need latency.
func CheckWAHALatency() (bool, time.Duration) {
	url := strings.TrimSuffix(config.App.WAHABaseURL, "/") + "/api/sessions"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, 0
	}
	req.Header.Set("Accept", "application/json")
	if config.App.WAHAAPIKey != "" {
		req.Header.Set("X-Api-Key", config.App.WAHAAPIKey)
	}

	start := time.Now()
	resp, err := healthClient.Do(req)
	latency := time.Since(start)
	if err != nil {
		return false, latency
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, latency
}
