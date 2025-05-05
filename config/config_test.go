package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig_Success(t *testing.T) {

	cfgContent := `{
		"laps": 5,
		"lapLen": 120,
		"penaltyLen": 30,
		"firingLines": 3,
		"start": "10:00:00",
		"startDelta": "00:00:30"
	}`

	tmpFile, err := os.CreateTemp("", "config_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(cfgContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	cfg, err := LoadConfig(tmpFile.Name())

	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.Laps != 5 {
		t.Errorf("Expected Laps = 5, got %d", cfg.Laps)
	}
	if cfg.LapLen != 120 {
		t.Errorf("Expected LapLen = 120, got %d", cfg.LapLen)
	}
	if cfg.PenaltyLen != 30 {
		t.Errorf("Expected PenaltyLen = 30, got %d", cfg.PenaltyLen)
	}
	if cfg.FiringLines != 3 {
		t.Errorf("Expected FiringLines = 3, got %d", cfg.FiringLines)
	}
	expectedStart := time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC)
	if !cfg.Start.Equal(expectedStart) {
		t.Errorf("Expected Start = %v, got %v", expectedStart, cfg.Start)
	}
	expectedDelta := 30 * time.Second
	if cfg.StartDelta != expectedDelta {
		t.Errorf("Expected StartDelta = %v, got %v", expectedDelta, cfg.StartDelta)
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {

	invalidContent := `{
		"laps": "five",
		"lapLen": "one hundred twenty"
	}`

	tmpFile, err := os.CreateTemp("", "config_invalid_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(invalidContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	cfg, err := LoadConfig(tmpFile.Name())

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if cfg != nil {
		t.Errorf("Expected nil, got %v", cfg)
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {

	cfg, err := LoadConfig("nonexistent_config.json")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if cfg != nil {
		t.Errorf("Expected nil, got %v", cfg)
	}
}
