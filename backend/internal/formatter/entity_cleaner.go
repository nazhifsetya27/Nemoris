package formatter

import (
	"regexp"
	"strings"
)

var repeatedSpaceRe = regexp.MustCompile(`\s+`)

// Clean trims unsafe repeated spaces and normalizes obvious text noise.
// Keeps semantics unchanged. Deterministic.
func Clean(text string) string {
	s := strings.TrimSpace(text)
	if s == "" {
		return s
	}
	s = repeatedSpaceRe.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}
