package processor

import (
	"biathlon/config"
	"biathlon/event"
	"testing"
	"time"
)

func TestAddLog(t *testing.T) {
	cfg := &config.Config{}
	events := []*event.Event{}
	p := NewProcessor(cfg, events)

	testTime := time.Date(2025, 5, 5, 12, 0, 0, 0, time.UTC)

	p.AddLog(testTime, "Test log message")

	if len(p.Logs) != 1 {
		t.Errorf("Expected 1 log, got %d", len(p.Logs))
	}

	expectedLog := "[12:00:00.000] Test log message"
	if p.Logs[0] != expectedLog {
		t.Errorf("Expected log '%s', got '%s'", expectedLog, p.Logs[0])
	}
}

func TestGetOrCreateCompetitor(t *testing.T) {
	cfg := &config.Config{}
	events := []*event.Event{}
	p := NewProcessor(cfg, events)

	comp := p.getOrCreateCompetitor(1)

	if comp.ID != 1 {
		t.Errorf("Expected Competitor ID 1, got %d", comp.ID)
	}
	if comp.NotStarted != true {
		t.Errorf("Expected Competitor NotStarted to be true, got %v", comp.NotStarted)
	}
	if comp.NotFinished != true {
		t.Errorf("Expected Competitor NotFinished to be true, got %v", comp.NotFinished)
	}
}

func TestProcessEvents(t *testing.T) {
	cfg := &config.Config{}
	events := []*event.Event{
		{CompetitorID: 1, EventID: 1, Time: time.Now()},
		{CompetitorID: 1, EventID: 2, ExtraParams: []string{"15:30:00.000"}, Time: time.Now()},
	}
	p := NewProcessor(cfg, events)

	p.ProcessEvents()

	comp := p.Competitors[1]
	if comp == nil {
		t.Errorf("Expected Competitor with ID 1, got nil")
	} else if comp.PlannedStart.IsZero() {
		t.Errorf("Expected Competitor PlannedStart to be set, got zero value")
	}
}
