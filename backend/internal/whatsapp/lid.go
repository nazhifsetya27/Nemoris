package whatsapp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"nemoris/internal/config"
	"nemoris/internal/utils"
)

const sessionName = "default"

// lidResponse is the WAHA GET /api/{session}/lids/{lid} response.
type lidResponse struct {
	LID string  `json:"lid"`
	PN  *string `json:"pn"` // phone number @c.us, or null if not resolvable
}

// ResolveLIDToPhone calls WAHA to resolve LID (e.g. 204917302104303@lid) to phone number (e.g. 628123456789@c.us).
// Returns empty string if WAHA returns pn: null or on error.
func ResolveLIDToPhone(lid string) (string, error) {
	lid = strings.TrimSpace(lid)
	if lid == "" || !strings.HasSuffix(lid, "@lid") {
		return "", nil
	}

	// WAHA expects @ escaped as %40 in path
	escaped := url.PathEscape(lid)
	apiURL := strings.TrimSuffix(config.App.WAHABaseURL, "/") + "/api/" + sessionName + "/lids/" + escaped

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	if config.App.WAHAAPIKey != "" {
		req.Header.Set("X-Api-Key", config.App.WAHAAPIKey)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("WAHA lids request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("WAHA lids status %d", resp.StatusCode)
	}

	var out lidResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("decode lids response: %w", err)
	}

	if out.PN == nil || strings.TrimSpace(*out.PN) == "" {
		utils.LogOutbound("LID not resolvable: " + lid)
		return "", nil
	}
	return strings.TrimSpace(*out.PN), nil
}

// ResolveSendTarget returns the chatId to use for sending. If target is a LID (@lid),
// tries to resolve to phone number via WAHA; otherwise returns target as-is.
func ResolveSendTarget(target string) string {
	target = strings.TrimSpace(target)
	if target == "" {
		return target
	}
	if !strings.HasSuffix(target, "@lid") {
		return target
	}
	pn, err := ResolveLIDToPhone(target)
	if err != nil {
		utils.LogOutbound("LID resolve failed: " + err.Error() + ", using LID")
		return target
	}
	if pn != "" {
		return pn
	}
	return target
}
