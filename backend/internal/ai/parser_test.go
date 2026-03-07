package ai

import (
	"strings"
	"testing"
)

func TestParse_EnglishReminder(t *testing.T) {
	// "remind me to pay electricity tomorrow 8pm" -> Task, Time filled
	result := Parse("remind me to pay electricity tomorrow 8pm")

	if result.Intent != "create_reminder" {
		t.Errorf("Intent = %q; want create_reminder", result.Intent)
	}
	if result.Lang != "en" {
		t.Errorf("Lang = %q; want en", result.Lang)
	}
	if result.Task != "pay electricity" {
		t.Errorf("Task = %q; want pay electricity", result.Task)
	}
	if result.Time == "" {
		t.Error("Time should be RFC3339, got empty")
	}
	if !strings.Contains(result.Time, "T20:00:00") {
		t.Errorf("Time should contain 20:00 (8pm), got %q", result.Time)
	}
}

func TestParse_IndonesianReminder(t *testing.T) {
	// "ingatkan saya untuk bayar listrik besok jam 8 malam" -> same canonical structure
	result := Parse("ingatkan saya untuk bayar listrik besok jam 8 malam")

	if result.Intent != "create_reminder" {
		t.Errorf("Intent = %q; want create_reminder", result.Intent)
	}
	if result.Lang != "id" {
		t.Errorf("Lang = %q; want id", result.Lang)
	}
	if result.Task != "bayar listrik" {
		t.Errorf("Task = %q; want bayar listrik", result.Task)
	}
	if result.Time == "" {
		t.Error("Time should be RFC3339, got empty")
	}
	if !strings.Contains(result.Time, "T20:00:00") {
		t.Errorf("Time should contain 20:00 (8pm), got %q", result.Time)
	}
}

func TestParse_NonReminderIntent(t *testing.T) {
	result := Parse("show pending tasks")
	if result.Intent != "list_reminders" {
		t.Errorf("Intent = %q; want list_reminders", result.Intent)
	}
	if result.Task != "" || result.Time != "" {
		t.Errorf("Task and Time should be empty for non-create_reminder; got Task=%q Time=%q", result.Task, result.Time)
	}
}

func TestParse_ReminderParseFails(t *testing.T) {
	// "remind" triggers create_reminder but body doesn't match ParseReminder pattern
	result := Parse("remind")
	if result.Intent != "create_reminder" {
		t.Errorf("Intent = %q; want create_reminder", result.Intent)
	}
	if result.Task != "" || result.Time != "" {
		t.Errorf("Task and Time should stay empty when ParseReminder fails; got Task=%q Time=%q", result.Task, result.Time)
	}
}
