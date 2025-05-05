package main

import (
	"biathlon/config"
	"biathlon/event"
	"biathlon/processor"
	"fmt"
	"os"
)

func runApp(cfgPath, evsPath string) ([]string, []string, error) {
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		return nil, nil, fmt.Errorf("error loading config: %v", err)
	}

	evs, err := event.LoadEvents(evsPath)
	if err != nil {
		return nil, nil, fmt.Errorf("error loading events: %v", err)
	}

	proc := processor.NewProcessor(cfg, evs)
	proc.ProcessEvents()

	var logs []string
	logs = append(logs, proc.Logs...)

	var results []string
	results = append(results, proc.GenerateResults()...)

	return logs, results, nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: <config_path> <events_path>")
		os.Exit(1)
	}

	logs, results, err := runApp(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("===Output log===")
	for _, log := range logs {
		fmt.Println(log)
	}

	fmt.Println("\n===Resulting table===")
	for _, row := range results {
		fmt.Println(row)
	}
}
