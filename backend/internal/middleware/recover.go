package middleware

import (
	"fmt"
	"net/http"

	"nemoris/internal/utils"
)

// Recover wraps the next handler and recovers any panic in the request flow.
// Logs the panic via utils.LogSystem, returns 500 with plain text "internal server error".
// Does not expose stack trace. Keeps the server process alive.
//
// Standard Go middleware signature: func(next http.Handler) http.Handler.
// Composable: use independently from WebhookAuth, e.g. Recover(WebhookAuth(handler)).
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				utils.LogSystem(fmt.Sprintf("panic recovered: %v", err))
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("internal server error"))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
