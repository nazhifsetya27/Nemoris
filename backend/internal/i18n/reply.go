package i18n

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
