package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
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
	logFile       *os.File
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), `Server Monitor CLI Tool

Usage:
  server-monitor [options]

Options:
`)
		flag.PrintDefaults()
	}
}

func main() {
	// CLI flags
	flag.IntVar(&interval, "interval", 0, "Interval in seconds to refresh stats (0 = run once)")
	flag.BoolVar(&enableJSON, "json", false, "Output stats in JSON format")
	flag.BoolVar(&enableLog, "log", false, "Log stats to monitor.log file")
	flag.Float64Var(&cpuThreshold, "cpu-threshold", 80, "CPU usage threshold (%)")
	flag.Float64Var(&memThreshold, "mem-threshold", 90, "Memory usage threshold (%)")
	flag.Float64Var(&diskThreshold, "disk-threshold", 90, "Disk usage threshold (%)")
	flag.StringVar(&exportPath, "export", "", "Export stats to CSV file (e.g., stats.csv)")
	flag.Parse()

	// Signal handling setup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle Ctrl+C or kill
	go handleSignals(cancel)

	// Logging setup
	if enableLog {
		var err error
		logFile, err = os.OpenFile("monitor.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		defer logFile.Close()
		logger = log.New(logFile, "", log.LstdFlags)
	}

	// Monitoring loop
	for {
		select {
		case <-ctx.Done():
			fmt.Println("🛑 Graceful shutdown triggered.")
			if logger != nil {
				logger.Println("Shutting down gracefully.")
			}
			return

		default:
			stats, err := CollectStats()
			if err != nil {
				log.Printf("Error collecting stats: %v", err)
				if logger != nil {
					logger.Printf("Error collecting stats: %v", err)
				}
			} else {
				summary := EvaluateThresholds(stats)

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
					jsonData, _ := json.MarshalIndent(data, "", "  ")
					fmt.Println(string(jsonData))
					if logger != nil {
						logger.Println(string(jsonData))
					}
				} else {
					output := FormatStats(stats, summary)
					fmt.Println(output)
					if logger != nil {
						logger.Printf("[%s] %s", time.Now().Format(time.RFC3339), summary)
					}
				}

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
				return
			}
			time.Sleep(time.Duration(interval) * time.Second)
		}
	}
}

// handleSignals cancels context on Ctrl+C or SIGTERM
func handleSignals(cancelFunc context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan // wait for interrupt
	cancelFunc()
}
