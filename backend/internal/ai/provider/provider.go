package provider

import (
	"context"
	"fmt"
	"hash/fnv"
	"time"

	"nemoris/internal/utils"
)

const providerTimeout = 3 * time.Second

func promptHash(body string) string {
	h := fnv.New32a()
	h.Write([]byte(body))
	return fmt.Sprintf("%08x", h.Sum32())
}

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
// Logs audit fields: prompt_hash, latency, fallback_reason. Never logs raw prompt.
func PrepareExternalCall(body string) (Result, error) {
	start := time.Now()
	ph := promptHash(body)

	ctx, cancel := context.WithTimeout(context.Background(), providerTimeout)
	defer cancel()

	resultCh := make(chan Result, 1)
	go func() {
		// Stub: no real I/O yet. Future: HTTP call, parse into Result.
		resultCh <- ZeroResult()
	}()

	select {
	case r := <-resultCh:
		utils.LogAI(fmt.Sprintf("provider fallback invoked prompt_hash=%s latency=%s fallback_reason=stub", ph, time.Since(start)))
		return r, nil
	case <-ctx.Done():
		utils.LogAI(fmt.Sprintf("provider fallback invoked prompt_hash=%s latency=%s fallback_reason=timeout", ph, time.Since(start)))
		return ZeroResult(), ctx.Err()
	}
}
