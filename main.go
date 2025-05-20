package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

func main() {
	stats, err := getFormattedStats()
	if err != nil {
		log.Fatalf("Error collecting stats: %v", err)
	}
	fmt.Println(stats)
}

// GetFormattedStats gathers and returns system stats as a human-readable string.
func getFormattedStats() (string, error) {
	// CPU usage
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return "", fmt.Errorf("CPU usage: %w", err)
	}

	// Memory usage
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return "", fmt.Errorf("Memory usage: %w", err)
	}

	// Disk usage
	diskStat, err := disk.Usage("/")
	if err != nil {
		return "", fmt.Errorf("Disk usage: %w", err)
	}

	// Load average
	loadStat, err := load.Avg()
	if err != nil {
		return "", fmt.Errorf("Load average: %w", err)
	}

	// Format output
	stats := fmt.Sprintf(`
==============================
 Server Performance Snapshot
==============================
CPU Usage     : %.2f%%
Memory Usage  : %.2f%% (%v / %v)
Disk Usage    : %.2f%% (%v / %v)
Load Average  : %.2f / %.2f / %.2f (1m / 5m / 15m)
`,
		cpuPercent[0],
		vmStat.UsedPercent, formatBytes(vmStat.Used), formatBytes(vmStat.Total),
		diskStat.UsedPercent, formatBytes(diskStat.Used), formatBytes(diskStat.Total),
		loadStat.Load1, loadStat.Load5, loadStat.Load15,
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
