package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

func main() {
	// CLI flag
	interval := flag.Int("interval", 0, "Interval in seconds to refresh stats (0 = run once)")
	flag.Parse()

	for {
		stats, err := GetFormattedStats()
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Println(stats)
		}

		if *interval <= 0 {
			break
		}
		time.Sleep(time.Duration(*interval) * time.Second)
	}
}
