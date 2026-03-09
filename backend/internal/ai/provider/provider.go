package provider

import (
	"context"
	"time"
)

const providerTimeout = 3 * time.Second

// Result holds the canonical structured contract for external AI responses.
// Matches the shape expected by the parser layer: intent, task, time, lang.
// Contract: explicit typed struct only. No raw prose. No map[string]interface{}.
type Result struct {
	Intent         string
	Task           string
	Time           string
	Lang           string
	RecurrenceType string
}

// ZeroResult returns the safe zero-value canonical struct for unsupported or empty provider results.
func ZeroResult() Result {
	return Result{
		Intent:         "unknown",
		Lang:           "en",
		Task:           "",
		Time:           "",
		RecurrenceType: "",
	}
}

// PrepareExternalCall prepares an external AI call and returns a stub-safe result.
// Enforces strict timeout; timeout failure returns ZeroResult().
// No real provider integration yet. No HTTP client, no SDK.
func PrepareExternalCall(body string) (Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), providerTimeout)
	defer cancel()

	resultCh := make(chan Result, 1)
	go func() {
		// Stub: no real I/O yet. Future: HTTP call, parse into Result.
		resultCh <- ZeroResult()
	}()

	select {
	case r := <-resultCh:
		return r, nil
	case <-ctx.Done():
		return ZeroResult(), ctx.Err()
	}
}
