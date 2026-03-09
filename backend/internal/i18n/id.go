package i18n

// Id holds Indonesian reply templates.
// Placeholders: {{task}}, {{time}}, {{count}}
var Id = map[string]string{
	"create_reminder_success":  "Siap, saya ingatkan {{time}}.",
	"list_reminders_success":   "Anda punya {{count}} pengingat.",
	"store_memory_success":     "Baik, saya simpan.",
	"retrieve_memory_success":  "Anda bilang: {{content}}",
	"retrieve_memory_empty":    "Saya tidak punya catatan itu.",
	"unknown":                  "Maaf, saya belum paham.",
}
