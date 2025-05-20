package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

var (
	interval      int
	enableJSON    bool
	enableLog     bool
	cpuThreshold  float64
	memThreshold  float64
	diskThreshold float64
	exportPath    string
	logger        *log.Logger
)

func main() {
	flag.IntVar(&interval, "interval", 0, "Interval in seconds to refresh stats (0 = run once)")
	flag.BoolVar(&enableJSON, "json", false, "Output stats in JSON format")
	flag.BoolVar(&enableLog, "log", false, "Log stats to monitor.log file")
	flag.Float64Var(&cpuThreshold, "cpu-threshold", 80, "CPU usage threshold (%)")
	flag.Float64Var(&memThreshold, "mem-threshold", 90, "Memory usage threshold (%)")
	flag.Float64Var(&diskThreshold, "disk-threshold", 90, "Disk usage threshold (%)")
	flag.StringVar(&exportPath, "export", "", "Export stats to CSV file (e.g., stats.csv)")
	flag.Parse()

	if enableLog {
		logFile, err := os.OpenFile("monitor.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		defer logFile.Close()
		logger = log.New(logFile, "", log.LstdFlags)
	}

	for {
		stats, err := CollectStats()
		if err != nil {
			log.Printf("Error: %v", err)
			if logger != nil {
				logger.Printf("Error collecting stats: %v", err)
			}
		} else {
			summary := EvaluateThresholds(stats)

			// Text or JSON output
			if enableJSON {
				data := struct {
					Timestamp string `json:"timestamp"`
					Stats     Stats  `json:"stats"`
					Summary   string `json:"summary"`
				}{
					Timestamp: time.Now().Format(time.RFC3339),
					Stats:     stats,
					Summary:   summary,
				}
				output, _ := json.MarshalIndent(data, "", "  ")
				fmt.Println(string(output))
				if logger != nil {
					logger.Println(string(output))
				}
			} else {
				output := FormatStats(stats, summary)
				fmt.Println(output)
				if logger != nil {
					logger.Printf("[%s] %s", time.Now().Format(time.RFC3339), summary)
				}
			}

			// CSV export
			if exportPath != "" {
				err := ExportToCSV(stats, exportPath)
				if err != nil {
					log.Printf("CSV export error: %v", err)
					if logger != nil {
						logger.Printf("CSV export error: %v", err)
					}
				}
			}
		}

		if interval <= 0 {
			break
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
