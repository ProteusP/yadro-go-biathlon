package competitor

import "time"

type Competitor struct {
	ID           int
	NotStarted   bool
	NotFinished  bool
	StartTime    time.Time
	FinishTime   time.Time
	PlannedStart time.Time
	ActualStart  time.Time
	TotalHits    int
	CurrentHits  int

	CurPenaltyStart  time.Time
	CurPenaltyEnd    time.Time
	TotalPenaltyTime time.Duration
	TotalPenaltyLen  int

	CurLapStart   time.Time
	CurLapEnd     time.Time
	LapDurations  []time.Duration
	TotalDuration time.Duration
}

func (c *Competitor) EnterPenalty(t time.Time) {
	c.CurPenaltyStart = t

}

func (c *Competitor) ExitPenalty(t time.Time, pLen int) {
	c.CurPenaltyEnd = t
	penaltyDuration := t.Sub(c.CurPenaltyStart)

	c.TotalPenaltyLen += pLen
	c.TotalPenaltyTime += penaltyDuration
}

func (c *Competitor) EndLap(t time.Time) {
	c.CurLapEnd = t
	duration := t.Sub(c.CurLapStart)
	c.LapDurations = append(c.LapDurations, duration)
	c.TotalDuration += duration
	c.CurLapStart = t
}
