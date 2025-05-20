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
	// CPU usage
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		log.Fatalf("Error getting CPU usage: %v", err)
	}

	// Memory usage
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		log.Fatalf("Error getting memory info: %v", err)
	}

	// Disk usage
	diskStat, err := disk.Usage("/")
	if err != nil {
		log.Fatalf("Error getting disk info: %v", err)
	}

	// Load average (Unix-like systems)
	loadStat, err := load.Avg()
	if err != nil {
		log.Fatalf("Error getting load average: %v", err)
	}

	// Output
	fmt.Printf("CPU Usage: %.2f%%\n", cpuPercent[0])
	fmt.Printf("Memory Usage: %.2f%% (%v / %v)\n", vmStat.UsedPercent, formatBytes(vmStat.Used), formatBytes(vmStat.Total))
	fmt.Printf("Disk Usage: %.2f%% (%v / %v)\n", diskStat.UsedPercent, formatBytes(diskStat.Used), formatBytes(diskStat.Total))
	fmt.Printf("Load Average (1/5/15 min): %.2f / %.2f / %.2f\n", loadStat.Load1, loadStat.Load5, loadStat.Load15)
}

// Helper to format bytes to human-readable format
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
