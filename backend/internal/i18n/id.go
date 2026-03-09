package i18n

// Id holds Indonesian reply templates.
// Placeholders: {{task}}, {{time}}, {{count}}, {{content}}
var Id = map[string]string{
	"create_reminder_success":  "Siap, saya ingatkan {{time}}.",
	"list_reminders_success":   "Anda punya {{count}} pengingat.",
	"store_memory_success":     "Baik, saya simpan.",
	"retrieve_memory_success":   "Saya menemukan ini: {{content}}",
	"retrieve_memory_empty":     "Saya belum menemukan data yang cocok.",
	"retrieve_memory_ambiguous": "Saya menemukan beberapa data, bisa lebih spesifik?",
	"unknown":                  "Maaf, saya belum paham.",
}
