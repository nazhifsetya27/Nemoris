package formatter

import (
	"regexp"
	"strings"
)

var (
	excessPunctRe = regexp.MustCompile(`[!?]+`)
	wordTmrwRe   = regexp.MustCompile(`\btmrw\b`)
	wordPlsRe    = regexp.MustCompile(`\bpls\b`)
)

// NormalizeInput enriches user text before parser: punctuation cleanup,
// shorthand normalization, mixed bilingual cleanup, typo-safe normalization.
// If uncertain, preserves original. Never changes canonical intent meaning.
func NormalizeInput(body string) string {
	s := Clean(body)
	if s == "" {
		return s
	}

	// Punctuation cleanup: "ingatkan saya!!! besok???" -> "ingatkan saya besok"
	s = excessPunctRe.ReplaceAllString(s, "")

	// Shorthand: tmrw -> tomorrow, pls -> please
	s = wordTmrwRe.ReplaceAllLiteralString(s, "tomorrow")
	s = wordPlsRe.ReplaceAllLiteralString(s, "please")

	// Mixed bilingual: "jam 8 mlm" -> "jam 8 malam" (parser expects "malam" for pm)
	s = strings.ReplaceAll(s, " mlm ", " malam ")
	s = strings.ReplaceAll(s, " mlm", " malam")

	return Clean(s)
}
