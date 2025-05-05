package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const TIME_FORMAT_NO_MS = "15:04:05"
const TIME_FORMAT_WITH_MS = "15:04:05.000"

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

	startTime, err := rawCfg.parseStartTime()
	if err != nil {
		return nil, err

	}

	startDelta, err := rawCfg.parseDelta()
	if err != nil {
		return nil, err
	}

	return &Config{
		Laps:        rawCfg.Laps,
		LapLen:      rawCfg.LapLen,
		PenaltyLen:  rawCfg.PenaltyLen,
		FiringLines: rawCfg.FiringLines,
		Start:       startTime,
		StartDelta:  startDelta,
	}, nil
}

func (rawCfg *ConfigRaw) parseStartTime() (time.Time, error) {
	stTime, err := time.Parse(TIME_FORMAT_NO_MS, rawCfg.Start)
	if err != nil {
		stTime, err = time.Parse(TIME_FORMAT_WITH_MS, rawCfg.Start)
		if err != nil {
			return stTime, fmt.Errorf("error while formatting start time: %v", err)
		}
	}
	return stTime, nil
}

func (rawCfg *ConfigRaw) parseDelta() (time.Duration, error) {
	pDelta, err := time.Parse(TIME_FORMAT_NO_MS, rawCfg.StartDelta)
	if err != nil {
		pDelta, err = time.Parse(TIME_FORMAT_WITH_MS, rawCfg.StartDelta)
		if err != nil {
			return 0, fmt.Errorf("error while formatting start delta: %v", err)
		}
	}

	formattedDelta := time.Duration(
		pDelta.Hour()*int(time.Hour) + pDelta.Minute()*int(time.Minute) + pDelta.Second()*int(time.Second),
	)
	return formattedDelta, nil
}
