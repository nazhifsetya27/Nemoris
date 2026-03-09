package i18n

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
// Falls back to en if lang is unknown; falls back to unknown message if key is unknown.
// data is reserved for future placeholder substitution (e.g. task, time).
func Build(lang string, key string, data map[string]string) string {
	t := getTemplates(lang)
	msg := t[key]
	if msg == "" {
		msg = t["unknown"]
	}
	if msg == "" {
		msg = En["unknown"]
	}
	return msg
}

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
