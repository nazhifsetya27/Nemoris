package formatter

// FormatOutput cleans user-facing text before send.
// Preserves existing i18n wording. Improves readability only, not rewrite meaning.
func FormatOutput(text string) string {
	return Clean(text)
}
