package i18n

// En holds English reply templates.
// Placeholders: {{task}}, {{time}}, {{count}}, {{content}}
var En = map[string]string{
	"create_reminder_success":  "Okay, I'll remind you {{time}}.",
	"list_reminders_success":   "You have {{count}} reminder(s).",
	"store_memory_success":     "Noted, I'll remember that.",
	"retrieve_memory_success":   "I found this: {{content}}",
	"retrieve_memory_empty":     "I couldn't find anything yet.",
	"retrieve_memory_ambiguous": "I found several matches, can you be more specific?",
	"unknown":                  "Sorry, I didn't understand that.",
}
