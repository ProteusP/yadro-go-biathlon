package event

import (
	"biathlon/config"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Event struct {
	Time         time.Time
	EventID      int
	CompetitorID int
	ExtraParams  []string
}

func ParseEvent(line string) (*Event, error) {

	endIdx := strings.Index(line, "]")
	timeStr := line[1:endIdx]
	eventPart := strings.TrimSpace(line[endIdx+1:])

	parsedTime, err := time.Parse(config.TIME_FORMAT_WITH_MS, timeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid time format: %v", err)
	}

	tokens := strings.Fields(eventPart)
	if len(tokens) < 2 {
		return nil, fmt.Errorf("not enough tokens: %v", err)
	}

	eventID, err := strconv.Atoi(tokens[0])
	if err != nil {
		return nil, fmt.Errorf("invalid event ID: %s", tokens[0])
	}

	competitorID, err := strconv.Atoi(tokens[1])
	if err != nil {
		return nil, fmt.Errorf("invalid competitor ID: %s", tokens[1])
	}

	extra := []string{}
	if len(tokens) > 2 {
		extra = tokens[2:]
	}

	return &Event{
		Time:         parsedTime,
		EventID:      eventID,
		CompetitorID: competitorID,
		ExtraParams:  extra,
	}, nil

}

func LoadEvents(filename string) ([]*Event, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var events []*Event
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		event, err := ParseEvent(line)
		if err != nil {
			return nil, fmt.Errorf("failed to parse line '%s': %v", line, err)
		}
		events = append(events, event)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
