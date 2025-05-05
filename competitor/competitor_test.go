package competitor

import (
	"testing"
	"time"
)

func TestEnterPenalty(t *testing.T) {
	c := &Competitor{}
	startTime := time.Now()

	c.EnterPenalty(startTime)

	if c.CurPenaltyStart != startTime {
		t.Errorf("Expected penalty start time %v, got %v", startTime, c.CurPenaltyStart)
	}
}

func TestExitPenalty(t *testing.T) {
	c := &Competitor{}
	startTime := time.Now()
	endTime := startTime.Add(10 * time.Second)

	c.EnterPenalty(startTime)
	c.ExitPenalty(endTime, 2)

	if c.TotalPenaltyLen != 2 {
		t.Errorf("Expected penalty length 2, got %d", c.TotalPenaltyLen)
	}

	if c.TotalPenaltyTime != 10*time.Second {
		t.Errorf("Expected penalty time 10s, got %v", c.TotalPenaltyTime)
	}
}
