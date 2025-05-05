package main

import (
	"testing"
)

func TestRunApp(t *testing.T) {
	tests := []struct {
		name        string
		cfgPath     string
		evsPath     string
		wantLogs    []string
		wantResults []string
		wantErr     bool
	}{
		{
			name:        "valid input",
			cfgPath:     "testdata/config.json",
			evsPath:     "testdata/events.txt",
			wantLogs:    []string{"[10:00:00.000] The competitor(1) registered"},
			wantResults: []string{"[NotStarted] 1 [{,}, {,}] {,} 0/10"},
			wantErr:     false,
		},
		{
			name:        "invalid config path",
			cfgPath:     "testdata/invalid_config.json",
			evsPath:     "testdata/events.txt",
			wantLogs:    nil,
			wantResults: nil,
			wantErr:     true,
		},
		{
			name:        "invalid events path",
			cfgPath:     "testdata/config.json",
			evsPath:     "testdata/invalid_events.txt",
			wantLogs:    nil,
			wantResults: nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLogs, gotResults, err := runApp(tt.cfgPath, tt.evsPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("runApp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !equal(gotLogs, tt.wantLogs) {
				t.Errorf("runApp() gotLogs = %v, want %v", gotLogs, tt.wantLogs)
			}
			if !equal(gotResults, tt.wantResults) {
				t.Errorf("runApp() gotResults = %v, want %v", gotResults, tt.wantResults)
			}
		})
	}
}

func equal(a, b []string) bool {
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
