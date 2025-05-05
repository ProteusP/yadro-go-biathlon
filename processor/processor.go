package processor

import (
	"biathlon/competitor"
	"biathlon/config"
	"biathlon/event"
	"fmt"
	"sort"
	"time"
)

const SHOTS_PER_FIRING_LINE = 5

type Processor struct {
	Config      *config.Config
	Competitors map[int]*competitor.Competitor
	Logs        []string
	Events      []*event.Event
}

func NewProcessor(cfg *config.Config, events []*event.Event) *Processor {
	return &Processor{
		Config:      cfg,
		Competitors: make(map[int]*competitor.Competitor),
		Events:      events,
	}
}

func (p *Processor) AddLog(t time.Time, msg string) {
	log := fmt.Sprintf("[%s] %s", t.Format("15:04:05.000"), msg)
	p.Logs = append(p.Logs, log)
}

func (p *Processor) getOrCreateCompetitor(id int) *competitor.Competitor {
	if c, exists := p.Competitors[id]; exists {
		return c
	}

	c := &competitor.Competitor{
		ID:          id,
		NotStarted:  true,
		NotFinished: true,
	}
	p.Competitors[id] = c
	return c
}

func (p *Processor) ProcessEvents() {

	for _, e := range p.Events {
		comp := p.getOrCreateCompetitor(e.CompetitorID)

		switch e.EventID {
		case 1: // Registration
			p.handleRegistration(e, comp)

		case 2: // Start time
			p.handleStartTime(e, comp)

		case 3: // On the start line
			p.handleOnStartLine(e, comp)

		case 4: // Started
			p.handleStarted(e, comp)

		case 5: // On the firing range
			p.handleOnFiringRange(e, comp)

		case 6: // Hit
			p.handleHit(e, comp)

		case 7: // Left the firing range
			p.handleLeftFiringRange(e, comp)

		case 8: // Entered the penalty lap(s)
			p.handleEnteredPLaps(e, comp)

		case 9: // Left the penalty lap(s)
			p.handleLeftPLaps(e, comp)

		case 10: // Ended the main lap
			p.handleEndedMainLap(e, comp)

		case 11: // Cant continue
			p.handleCantContinue(e, comp)
		}
	}

}

func (p *Processor) GenerateResults() []string {

	comps := []*competitor.Competitor{}

	for _, c := range p.Competitors {
		comps = append(comps, c)
	}

	sort.SliceStable(comps, func(i, j int) bool {
		ci, cj := comps[i], comps[j]

		if ci.NotStarted || ci.NotFinished {
			return false
		}

		if cj.NotStarted || cj.NotFinished {
			return true
		}
		return ci.TotalDuration < cj.TotalDuration
	})

	results := []string{}

	for _, c := range comps {
		entry := ""

		entry += p.parseTimeAndStatus(c)

		entry += p.parseID(c)

		entry += p.parseMainLaps(c)

		entry += p.parsePLaps(c)

		entry += p.parseHitsAndShots(c)

		results = append(results, entry)
	}
	return results
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	millis := int(d.Milliseconds()) % 1000

	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, millis)
}

func (p *Processor) handleRegistration(e *event.Event, _ *competitor.Competitor) {
	log := fmt.Sprintf("The competitor(%d) registered", e.CompetitorID)
	p.AddLog(e.Time, log)
}

func (p *Processor) handleStartTime(e *event.Event, comp *competitor.Competitor) {
	if len(e.ExtraParams) == 1 {
		startTime := e.ExtraParams[0]
		parsedStartTime, err := time.Parse("15:04:05.000", startTime)
		comp.CurLapStart = parsedStartTime
		if err == nil {
			comp.PlannedStart = parsedStartTime
			log := fmt.Sprintf("The start time of competitor(%d) was set by a draw to %s", e.CompetitorID, startTime)
			p.AddLog(e.Time, log)
		}
	}
}

func (p *Processor) handleOnStartLine(e *event.Event, _ *competitor.Competitor) {
	log := fmt.Sprintf("The competitor(%d) is on the start line", e.CompetitorID)
	p.AddLog(e.Time, log)
}

func (p *Processor) handleStarted(e *event.Event, comp *competitor.Competitor) {
	comp.ActualStart = e.Time

	startWindow := comp.PlannedStart.Add(p.Config.StartDelta)

	if comp.ActualStart.After(startWindow) {
		log := fmt.Sprintf("The competitor(%d) is disqualified", e.CompetitorID)
		p.AddLog(e.Time, log)
	} else {
		log := fmt.Sprintf("The competitor(%d) has started", e.CompetitorID)
		comp.NotStarted = false
		p.AddLog(e.Time, log)
	}
}

func (p *Processor) handleOnFiringRange(e *event.Event, comp *competitor.Competitor) {
	comp.CurrentHits = 0
	if len(e.ExtraParams) == 1 {
		firingRange := e.ExtraParams[0]
		log := fmt.Sprintf("The competitor(%d) is on the firing range(%s)", e.CompetitorID, firingRange)
		p.AddLog(e.Time, log)
	}
}

func (p *Processor) handleHit(e *event.Event, comp *competitor.Competitor) {
	comp.TotalHits++
	comp.CurrentHits++
	if len(e.ExtraParams) >= 1 {
		target := e.ExtraParams[0]
		log := fmt.Sprintf("The target(%s) has been hit by competitor(%d)", target, e.CompetitorID)
		p.AddLog(e.Time, log)
	}
}

func (p *Processor) handleLeftFiringRange(e *event.Event, _ *competitor.Competitor) {
	log := fmt.Sprintf("The competitor(%d) left the firing range", e.CompetitorID)
	p.AddLog(e.Time, log)
}

func (p *Processor) handleEnteredPLaps(e *event.Event, comp *competitor.Competitor) {
	comp.EnterPenalty(e.Time)

	log := fmt.Sprintf("The competitor(%d) entered the penalty laps", e.CompetitorID)
	p.AddLog(e.Time, log)
}

func (p *Processor) handleLeftPLaps(e *event.Event, comp *competitor.Competitor) {
	missedShots := SHOTS_PER_FIRING_LINE - comp.CurrentHits
	pLen := missedShots * p.Config.PenaltyLen
	comp.ExitPenalty(e.Time, pLen)

	log := fmt.Sprintf("The competitor(%d) left the penalty laps", e.CompetitorID)
	p.AddLog(e.Time, log)
}

func (p *Processor) handleEndedMainLap(e *event.Event, comp *competitor.Competitor) {
	comp.EndLap(e.Time)

	// Finished
	if len(comp.LapDurations) == p.Config.Laps {
		comp.NotFinished = false
		comp.FinishTime = e.Time
	}

	log := fmt.Sprintf("The competitor(%d) ended the main lap", e.CompetitorID)
	p.AddLog(e.Time, log)
}

func (p *Processor) handleCantContinue(e *event.Event, _ *competitor.Competitor) {
	comment := ""
	if len(e.ExtraParams) > 0 {
		for _, word := range e.ExtraParams {
			comment += word + " "
		}
	}
	log := fmt.Sprintf("The competitor can`t continue: %s", comment)
	p.AddLog(e.Time, log)
}

func (p *Processor) parseMainLaps(c *competitor.Competitor) string {
	res := "["
	if len(c.LapDurations) > 0 {
		d := c.LapDurations[0]
		avgSpeed := float64(p.Config.LapLen) / d.Seconds()
		res += fmt.Sprintf("{%s, %.3f}", formatDuration(d), avgSpeed)

		if len(c.LapDurations) > 1 {
			for _, d := range c.LapDurations[1:] {
				avgSpeed := float64(p.Config.LapLen) / d.Seconds()
				res += fmt.Sprintf(", {%s, %.3f}", formatDuration(d), avgSpeed)
			}
		}
	}

	if len(c.LapDurations) < p.Config.Laps {
		for i := len(c.LapDurations) + 1; i <= p.Config.Laps; i++ {
			if i == 1 {
				res += "{,}"
			} else {
				res += ", {,}"
			}
		}
	}
	res += "] "
	return res
}

func (p *Processor) parseTimeAndStatus(c *competitor.Competitor) string {
	res := ""
	switch {
	case c.NotStarted:
		res += "[NotStarted] "

	case c.NotFinished:
		res += "[NotFinished] "

	default:
		res += fmt.Sprintf("[%v] ", formatDuration(c.TotalDuration))
	}
	return res
}

func (p *Processor) parsePLaps(c *competitor.Competitor) string {
	res := ""
	if c.TotalPenaltyLen > 0 {
		avgPenSpeed := float64(c.TotalPenaltyLen) / c.TotalPenaltyTime.Seconds()
		res += fmt.Sprintf("{%s, %.3f} ", formatDuration(c.TotalPenaltyTime), avgPenSpeed)
	} else {
		res += "{,} "
	}
	return res
}

func (p *Processor) parseHitsAndShots(c *competitor.Competitor) string {
	totalShots := SHOTS_PER_FIRING_LINE * p.Config.FiringLines
	return fmt.Sprintf("%d/%d", c.TotalHits, totalShots)
}

func (p *Processor) parseID(c *competitor.Competitor) string {
	return fmt.Sprintf("%d ", c.ID)
}
