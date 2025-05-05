package main

import (
	"biathlon/config"
	"biathlon/event"
	"biathlon/processor"
	"bufio"
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

	if len(os.Args) == 4 {
		fmt.Println("Usage: <config_path> <events_path> [output_logs_path] [results_path]")
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

	// Write output log & resulting table to files
	if len(os.Args) == 5 {
		writeInFiles(os.Args[3], os.Args[4], logs, results)
	}
}

func writeInFiles(logPath, resPath string, logs, results []string) {
	if err := writeLinesToFile(logPath, logs); err != nil {
		fmt.Printf("Error writing logs: %v\n", err)
		return
	}

	if err := writeLinesToFile(resPath, results); err != nil {
		fmt.Printf("Error writing results: %v\n", err)
		return
	}

	fmt.Println("\n***Logs and results successfully written to files***")
}

func writeLinesToFile(filename string, lines []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("failed to write line: %w", err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush buffer: %w", err)
	}

	return nil
}
