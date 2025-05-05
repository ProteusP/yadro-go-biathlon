package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type ConfigRaw struct {
	Laps        int    `json:"laps"`
	LapLen      int    `json:"lapLen"`
	PenaltyLen  int    `json:"penaltyLen"`
	FiringLines int    `json:"firingLines"`
	Start       string `json:"start"`
	StartDelta  string `json:"startDelta"`
}

type Config struct {
	Laps        int
	LapLen      int
	PenaltyLen  int
	FiringLines int
	Start       time.Time
	StartDelta  time.Duration
}

func LoadConfig(path string) (*Config, error) {

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error while reading file")
	}

	var rawCfg ConfigRaw
	err = json.Unmarshal(raw, &rawCfg)
	if err != nil {
		return nil, fmt.Errorf("error while parsing json")
	}

	startTime, err := time.Parse("15:04:05", rawCfg.Start)
	if err != nil {
		startTime, err = time.Parse("15:04:05.000", rawCfg.Start)
		if err != nil {
			return nil, fmt.Errorf("error while formatting start time: %v", err)
		}
	}

	parsedDelta, err := time.Parse("15:04:05", rawCfg.StartDelta)
	if err != nil {
		parsedDelta, err = time.Parse("15:04:05.000", rawCfg.StartDelta)
		if err != nil {
			return nil, fmt.Errorf("error while formatting start delta: %v", err)
		}
	}

	startDelta := time.Duration(
		parsedDelta.Hour()*int(time.Hour) + parsedDelta.Minute()*int(time.Minute) + parsedDelta.Second()*int(time.Second),
	)

	return &Config{
		Laps:        rawCfg.Laps,
		LapLen:      rawCfg.LapLen,
		PenaltyLen:  rawCfg.PenaltyLen,
		FiringLines: rawCfg.FiringLines,
		Start:       startTime,
		StartDelta:  startDelta,
	}, nil
}
