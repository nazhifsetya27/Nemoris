package i18n

import "strings"

// intentToKey maps canonical intent to reply template key.
// Unknown intents fall back to "unknown".
var intentToKey = map[string]string{
	"create_reminder": "create_reminder_success",
	"list_reminders":  "list_reminders_success",
	"store_memory":    "store_memory_success",
	"unknown":         "unknown",
}

// BuildFromIntent returns localized reply for the given lang and canonical intent.
// Maps intent to template key internally; unknown intents use "unknown".
// Delegates to Build for actual lookup.
func BuildFromIntent(lang string, intent string, data map[string]string) string {
	key := intentToKey[intent]
	if key == "" {
		key = "unknown"
	}
	return Build(lang, key, data)
}

// Build returns localized reply text for the given lang and key.
// Fallback chain: lang template key → english same key → english unknown.
// No empty string may escape from this package.
// Replaces {{task}}, {{time}}, {{count}} when data contains those keys.
// Missing placeholders are left as-is (template stays readable).
func Build(lang string, key string, data map[string]string) string {
	t := getTemplates(lang)
	msg := t[key]
	if msg == "" {
		msg = En[key]
	}
	if msg == "" {
		msg = En["unknown"]
	}
	for k, v := range data {
		msg = strings.ReplaceAll(msg, "{{"+k+"}}", v)
	}
	if msg == "" {
		return En["unknown"]
	}
	return msg
}

// getTemplates returns templates for lang; unknown lang falls back to En.
func getTemplates(lang string) map[string]string {
	switch lang {
	case "id":
		return Id
	case "en":
		return En
	default:
		return En
	}
}
