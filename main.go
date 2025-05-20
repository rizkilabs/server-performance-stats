package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

func main() {
	// Define CLI flag
	interval := flag.Int("interval", 0, "Interval in seconds to refresh stats (0 = run once)")
	flag.Parse()

	// Run loop
	for {
		stats, err := getFormattedStats()
		if err != nil {
			log.Printf("Error collecting stats: %v", err)
		} else {
			fmt.Println(stats)
		}

		if *interval <= 0 {
			break
		}

		time.Sleep(time.Duration(*interval) * time.Second)
	}
}

// GetFormattedStats gathers system stats and formats them.
func getFormattedStats() (string, error) {
	osName := runtime.GOOS

	// CPU
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return "", fmt.Errorf("CPU usage: %w", err)
	}
	cpuVal := cpuPercent[0]

	// Memory
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return "", fmt.Errorf("Memory usage: %w", err)
	}

	// Disk
	diskStat, err := disk.Usage("/")
	if err != nil {
		return "", fmt.Errorf("Disk usage: %w", err)
	}

	// Load Average (Unix only)
	var loadLine string
	if osName == "linux" || osName == "darwin" {
		loadStat, err := load.Avg()
		if err != nil {
			loadLine = "Load Average  : unavailable (error)"
		} else {
			loadLine = fmt.Sprintf("Load Average  : %.2f / %.2f / %.2f (1m / 5m / 15m)", loadStat.Load1, loadStat.Load5, loadStat.Load15)
		}
	} else {
		loadLine = "Load Average  : not supported on this OS"
	}

	// Determine status summary
	summary := "System status: Normal"
	switch {
	case cpuVal > 80:
		summary = "⚠️  High CPU usage detected!"
	case vmStat.UsedPercent > 90:
		summary = "⚠️  High memory usage detected!"
	case diskStat.UsedPercent > 90:
		summary = "⚠️  Disk almost full!"
	}

	// Format output
	stats := fmt.Sprintf(`
==============================
 Server Performance Snapshot
==============================
Operating System: %s
CPU Usage       : %.2f%%
Memory Usage    : %.2f%% (%v / %v)
Disk Usage      : %.2f%% (%v / %v)
%s
------------------------------
%s
`, osName,
		cpuVal,
		vmStat.UsedPercent, formatBytes(vmStat.Used), formatBytes(vmStat.Total),
		diskStat.UsedPercent, formatBytes(diskStat.Used), formatBytes(diskStat.Total),
		loadLine,
		summary,
	)

	return stats, nil
}

// formatBytes converts bytes to a human-readable format.
func formatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
