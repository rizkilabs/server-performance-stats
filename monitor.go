package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"encoding/csv"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

// Stats struct to hold performance metrics
type Stats struct {
	OS            string
	CPUPercent    float64
	MemoryUsedPct float64
	DiskUsedPct   float64
	Load1         float64
	Load5         float64
	Load15        float64
}

// CollectStats gathers performance metrics
func CollectStats() (Stats, error) {
	var s Stats

	s.OS = runtime.GOOS

	cpuPercents, err := cpu.Percent(0, false)
	if err != nil {
		return s, err
	}
	if len(cpuPercents) > 0 {
		s.CPUPercent = cpuPercents[0]
	}

	vm, err := mem.VirtualMemory()
	if err != nil {
		return s, err
	}
	s.MemoryUsedPct = vm.UsedPercent

	diskStat, err := disk.Usage("/")
	if err != nil {
		return s, err
	}
	s.DiskUsedPct = diskStat.UsedPercent

	if s.OS == "linux" || s.OS == "darwin" {
		loadAvg, err := load.Avg()
		if err == nil {
			s.Load1 = loadAvg.Load1
			s.Load5 = loadAvg.Load5
			s.Load15 = loadAvg.Load15
		}
	}

	return s, nil
}

// EvaluateThresholds provides a warning summary
func EvaluateThresholds(s Stats) string {
	var warnings []string
	if s.CPUPercent > cpuThreshold {
		warnings = append(warnings, fmt.Sprintf("High CPU usage: %.2f%%", s.CPUPercent))
	}
	if s.MemoryUsedPct > memThreshold {
		warnings = append(warnings, fmt.Sprintf("High memory usage: %.2f%%", s.MemoryUsedPct))
	}
	if s.DiskUsedPct > diskThreshold {
		warnings = append(warnings, fmt.Sprintf("High disk usage: %.2f%%", s.DiskUsedPct))
	}

	if len(warnings) > 0 {
		return "⚠️ " + fmt.Sprintf("%s", warnings)
	}
	return "✅ All systems normal"
}

// FormatStats returns a human-readable string of the stats
func FormatStats(s Stats, summary string) string {
	return fmt.Sprintf(`
OS:            %s
CPU Usage:     %.2f%%
Memory Usage:  %.2f%%
Disk Usage:    %.2f%%
Load Average:  %.2f, %.2f, %.2f
Summary:       %s
`, s.OS, s.CPUPercent, s.MemoryUsedPct, s.DiskUsedPct, s.Load1, s.Load5, s.Load15, summary)
}

// ExportToCSV appends stats to a CSV file
func ExportToCSV(s Stats, filePath string) error {
	_, err := os.Stat(filePath)
	fileExists := !os.IsNotExist(err)

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if !fileExists {
		header := []string{"Timestamp", "OS", "CPU (%)", "Memory (%)", "Disk (%)", "Load 1", "Load 5", "Load 15"}
		if err := writer.Write(header); err != nil {
			return err
		}
	}

	record := []string{
		time.Now().Format(time.RFC3339),
		s.OS,
		fmt.Sprintf("%.2f", s.CPUPercent),
		fmt.Sprintf("%.2f", s.MemoryUsedPct),
		fmt.Sprintf("%.2f", s.DiskUsedPct),
		fmt.Sprintf("%.2f", s.Load1),
		fmt.Sprintf("%.2f", s.Load5),
		fmt.Sprintf("%.2f", s.Load15),
	}

	return writer.Write(record)
}
