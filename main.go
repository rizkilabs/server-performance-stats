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
	enableJSON bool
	interval   int
	enableLog  bool
	logger     *log.Logger
)

func main() {
	// CLI flags
	flag.IntVar(&interval, "interval", 0, "Interval in seconds to refresh stats (0 = run once)")
	flag.BoolVar(&enableJSON, "json", false, "Output stats in JSON format")
	flag.BoolVar(&enableLog, "log", false, "Log stats to monitor.log file")
	flag.Parse()

	// Setup logger if logging enabled
	if enableLog {
		logFile, err := os.OpenFile("monitor.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		defer logFile.Close()
		logger = log.New(logFile, "", log.LstdFlags)
	}

	for {
		stats, summary, err := CollectStats()
		if err != nil {
			log.Printf("Error: %v", err)
			if logger != nil {
				logger.Printf("Error collecting stats: %v", err)
			}
		} else {
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

				output, err := json.MarshalIndent(data, "", "  ")
				if err != nil {
					log.Printf("Error encoding JSON: %v", err)
					continue
				}
				fmt.Println(string(output))
				if logger != nil {
					logger.Println(string(output))
				}
			} else {
				formatted := FormatStats(stats, summary)
				fmt.Println(formatted)
				if logger != nil {
					logEntry := fmt.Sprintf("[%s] %s\n", time.Now().Format(time.RFC3339), summary)
					logger.Println(logEntry)
				}
			}
		}

		if interval <= 0 {
			break
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
