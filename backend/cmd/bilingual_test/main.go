package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"nemoris/internal/config"
	"nemoris/internal/database"
	"nemoris/internal/service"
)

type testCase struct {
	ID          int    `json:"id"`
	Body        string `json:"body"`
	Description string `json:"description"`
}

type testResult struct {
	ID          int    `json:"id"`
	Request     string `json:"request"`
	Response    string `json:"response"`
	Description string `json:"description"`
}

func main() {
	config.Load()
	database.Connect()
	database.Migrate()

	from := "test-bilingual-001"
	if config.App.BotNumber == from {
		from = "test-bilingual-002"
	}

	cases := []struct {
		id          int
		body        string
		description string
	}{
		{1, "remind me to pay electricity tomorrow 8pm", "English reminder with time"},
		{2, "ingatkan saya untuk bayar listrik besok jam 8 malam", "Indonesian reminder with time"},
		{3, "list", "English list reminders"},
		{4, "daftar", "Indonesian list reminders"},
		{5, "ping", "Ping"},
	}

	var results []testResult
	for _, c := range cases {
		response := service.ProcessMessage(from, c.body)
		results = append(results, testResult{
			ID:          c.id,
			Request:     c.body,
			Response:    response,
			Description: c.description,
		})
	}

	outPath := "test/results.json"
	if wd, err := os.Getwd(); err == nil && filepath.Base(wd) == "backend" {
		outPath = filepath.Join("..", "test", "results.json")
	}

	data, _ := json.MarshalIndent(results, "", "  ")
	_ = os.WriteFile(outPath, data, 0644)
}
