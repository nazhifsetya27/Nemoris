package service

import (
	"strings"

	"github.com/google/uuid"

	"nemoris/internal/i18n"
	"nemoris/internal/model"
	"nemoris/internal/repository"
)

// ListMemories returns memories filtered by sender, newest first.
func ListMemories(from string) ([]model.Memory, error) {
	return repository.ListMemoriesByFrom(from)
}

// RetrieveMemoryReply returns localized reply for memory retrieval.
// Detects latest vs keyword query, uses repository, preserves sender scope and lang.
func RetrieveMemoryReply(from string, body string, lang string) string {
	lower := strings.ToLower(strings.TrimSpace(body))

	isLatest := strings.Contains(lower, "what did i tell") || strings.Contains(lower, "what did i say") ||
		strings.Contains(lower, "apa yang saya bilang") || strings.Contains(lower, "apa yang aku bilang") ||
		strings.Contains(lower, "apa yang saya katakan") || strings.Contains(lower, "apa yang aku katakan")

	if isLatest {
		m, err := repository.GetLatestMemoryBySender(from)
		if err != nil {
			return i18n.Build(lang, "retrieve_memory_empty", nil)
		}
		if m.ID == uuid.Nil {
			return i18n.Build(lang, "retrieve_memory_empty", nil)
		}
		return i18n.Build(lang, "retrieve_memory_success", map[string]string{"content": m.Content})
	}

	keyword := extractRetrieveKeyword(lower)
	if keyword == "" {
		m, err := repository.GetLatestMemoryBySender(from)
		if err != nil {
			return i18n.Build(lang, "retrieve_memory_empty", nil)
		}
		if m.ID == uuid.Nil {
			return i18n.Build(lang, "retrieve_memory_empty", nil)
		}
		return i18n.Build(lang, "retrieve_memory_success", map[string]string{"content": m.Content})
	}

	memories, err := repository.FindMemoriesByKeyword(from, keyword)
	if err != nil {
		return i18n.Build(lang, "retrieve_memory_empty", nil)
	}
	if len(memories) == 0 {
		return i18n.Build(lang, "retrieve_memory_empty", nil)
	}
	return i18n.Build(lang, "retrieve_memory_success", map[string]string{"content": memories[0].Content})
}

func extractRetrieveKeyword(lower string) string {
	if idx := strings.Index(lower, "what is my "); idx >= 0 {
		return strings.TrimSpace(lower[idx+11:])
	}
	if strings.Contains(lower, "berapa nomor rekening") {
		return "rekening"
	}
	if idx := strings.Index(lower, "what do you remember about "); idx >= 0 {
		return strings.TrimSpace(lower[idx+27:])
	}
	if idx := strings.Index(lower, "what did you remember about "); idx >= 0 {
		return strings.TrimSpace(lower[idx+28:])
	}
	if idx := strings.Index(lower, "what do you remember"); idx >= 0 {
		return strings.TrimSpace(lower[idx+20:])
	}
	if idx := strings.Index(lower, "what did you remember"); idx >= 0 {
		return strings.TrimSpace(lower[idx+21:])
	}
	return ""
}
