package main

import (
	"SlitheringJake/pkg/slitheringjake"
	"encoding/json"
	"log"
	"os"
	"testing"
)

func TestNewSlitheringJake(t *testing.T) {
	configFile := ".local/config.json"

	content, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("failed to load config: %s", err.Error())
	}

	var config slitheringjake.Config
	if err := json.Unmarshal(content, &config); err != nil {
		t.Fatalf("failed to parse config: %s", err.Error())
	}

	jake, err := slitheringjake.NewSlitheringJake(config)
	if err != nil {
		t.Fatalf("Failed to create bot: %s", err.Error())
	}

	if jake == nil {
		t.Fatal("Expected non-nil SlitheringJake instance")
	}

	chain := jake.Bot.GetChain("log")
	defer jake.Bot.PutChain("log")
	if chain == nil {
		t.Fatal("Expected non-nil chain")
	}

	for i := 0; i < 10; i++ {
		sentence, score := chain.Contains("mario")
		log.Printf("[%f] %s", score, sentence)
	}
}
