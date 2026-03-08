package middleware

import (
	"net/http"
	"strings"

	"nemoris/internal/config"
	"nemoris/internal/utils"
)

const webhookSecretHeader = "X-Webhook-Secret"

// WebhookAuth validates X-Webhook-Secret header against config.WebhookSecret.
// Returns 401 if server secret is unset, header is missing, or secret does not match.
// Safe: never logs the secret, only authorized/unauthorized outcome.
//
// Go middleware (http.Handler) vs Express: Go uses the standard Handler interface—
// func(next http.Handler) http.Handler. Express uses (req, res, next). Both chain
// request through middleware before handler. Go's form works with any Handler,
// not just HandlerFunc, so it's more generic and composable.
func WebhookAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := strings.TrimSpace(config.App.WebhookSecret)

		if expected == "" {
			utils.LogSecurity("webhook blocked: server secret missing")
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		got := strings.TrimSpace(r.Header.Get(webhookSecretHeader))

		// Dev only: fallback to query param ?token= for WAHA compatibility
		if got == "" && config.App.AppEnv == "dev" {
			got = strings.TrimSpace(r.URL.Query().Get("token"))
		}

		if got == "" {
			utils.LogSecurity("webhook unauthorized: header missing")
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		if got != expected {
			utils.LogSecurity("webhook unauthorized: secret mismatch")
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		utils.LogSecurity("webhook authorized")
		next.ServeHTTP(w, r)
	})
}
