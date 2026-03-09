package i18n

// En holds English reply templates.
// Placeholders: {{task}}, {{time}}, {{count}}
var En = map[string]string{
	"create_reminder_success": "Okay, I'll remind you {{time}}.",
	"list_reminders_success":  "You have {{count}} reminder(s).",
	"store_memory_success":    "Noted, I'll remember that.",
	"unknown":                 "Sorry, I didn't understand that.",
}
