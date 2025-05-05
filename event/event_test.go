package event

import (
	"testing"
	"time"
)

func TestParseEvent(t *testing.T) {
	tests := []struct {
		input    string
		expected *Event
		err      bool
	}{
		{
			input: "[12:34:56.789] 1 1001 start",
			expected: &Event{
				Time:         time.Date(0, 1, 1, 12, 34, 56, 789000000, time.UTC),
				EventID:      1,
				CompetitorID: 1001,
				ExtraParams:  []string{"start"},
			},
			err: false,
		},
		{
			input:    "[invalid time] 1 1001 start",
			expected: nil,
			err:      true,
		},
		{
			input: "[12:34:56.789] 1 1001",
			expected: &Event{
				Time:         time.Date(0, 1, 1, 12, 34, 56, 789000000, time.UTC),
				EventID:      1,
				CompetitorID: 1001,
				ExtraParams:  nil,
			},
			err: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseEvent(tt.input)
			if (err != nil) != tt.err {
				t.Errorf("ParseEvent() error = %v, wantErr %v", err, tt.err)
				return
			}
			if got != nil && !got.Time.Equal(tt.expected.Time) {
				t.Errorf("ParseEvent() Time = %v, want %v", got.Time, tt.expected.Time)
			}
			if got != nil && got.EventID != tt.expected.EventID {
				t.Errorf("ParseEvent() EventID = %v, want %v", got.EventID, tt.expected.EventID)
			}
			if got != nil && got.CompetitorID != tt.expected.CompetitorID {
				t.Errorf("ParseEvent() CompetitorID = %v, want %v", got.CompetitorID, tt.expected.CompetitorID)
			}
			if got != nil && !equalStringSlices(got.ExtraParams, tt.expected.ExtraParams) {
				t.Errorf("ParseEvent() ExtraParams = %v, want %v", got.ExtraParams, tt.expected.ExtraParams)
			}
		})
	}
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
