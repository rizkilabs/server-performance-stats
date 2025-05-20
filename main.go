package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"
)

func main() {
	interval := flag.Int("interval", 0, "Interval in seconds to refresh stats (0 = run once)")
	asJSON := flag.Bool("json", false, "Output stats in JSON format")
	flag.Parse()

	for {
		stats, summary, err := CollectStats()
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			if *asJSON {
				data := struct {
					Stats   Stats  `json:"stats"`
					Summary string `json:"summary"`
				}{stats, summary}

				output, err := json.MarshalIndent(data, "", "  ")
				if err != nil {
					log.Printf("Error encoding JSON: %v", err)
				} else {
					fmt.Println(string(output))
				}
			} else {
				fmt.Println(FormatStats(stats, summary))
			}
		}

		if *interval <= 0 {
			break
		}
		time.Sleep(time.Duration(*interval) * time.Second)
	}
}
