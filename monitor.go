package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

type Stats struct {
	OS            string  `json:"os"`
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryUsed    string  `json:"memory_used"`
	MemoryTotal   string  `json:"memory_total"`
	MemoryUsedPct float64 `json:"memory_used_percent"`
	DiskUsed      string  `json:"disk_used"`
	DiskTotal     string  `json:"disk_total"`
	DiskUsedPct   float64 `json:"disk_used_percent"`
	Load1         float64 `json:"load_1,omitempty"`
	Load5         float64 `json:"load_5,omitempty"`
	Load15        float64 `json:"load_15,omitempty"`
	LoadMsg       string  `json:"load_msg,omitempty"`
}

func CollectStats() (Stats, error) {
	var s Stats
	s.OS = runtime.GOOS

	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return s, fmt.Errorf("CPU usage: %w", err)
	}
	s.CPUPercent = cpuPercent[0]

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return s, fmt.Errorf("Memory usage: %w", err)
	}
	s.MemoryUsedPct = vmStat.UsedPercent
	s.MemoryUsed = formatBytes(vmStat.Used)
	s.MemoryTotal = formatBytes(vmStat.Total)

	diskStat, err := disk.Usage("/")
	if err != nil {
		return s, fmt.Errorf("Disk usage: %w", err)
	}
	s.DiskUsedPct = diskStat.UsedPercent
	s.DiskUsed = formatBytes(diskStat.Used)
	s.DiskTotal = formatBytes(diskStat.Total)

	if s.OS == "linux" || s.OS == "darwin" {
		loadStat, err := load.Avg()
		if err == nil {
			s.Load1 = loadStat.Load1
			s.Load5 = loadStat.Load5
			s.Load15 = loadStat.Load15
		} else {
			s.LoadMsg = "Error retrieving load average"
		}
	} else {
		s.LoadMsg = "Load average not supported on this OS"
	}

	return s, nil
}

// Evaluates thresholds and returns a summary string
func EvaluateThresholds(s Stats) string {
	switch {
	case s.CPUPercent > cpuThreshold:
		return "⚠️  High CPU usage detected!"
	case s.MemoryUsedPct > memThreshold:
		return "⚠️  High memory usage detected!"
	case s.DiskUsedPct > diskThreshold:
		return "⚠️  Disk almost full!"
	default:
		return "✅ System status: Normal"
	}
}

func FormatStats(s Stats, summary string) string {
	loadLine := ""
	if s.LoadMsg != "" {
		loadLine = fmt.Sprintf("Load Average  : %s", s.LoadMsg)
	} else {
		loadLine = fmt.Sprintf("Load Average  : %.2f / %.2f / %.2f (1m / 5m / 15m)", s.Load1, s.Load5, s.Load15)
	}

	return fmt.Sprintf(`
==============================
 Server Performance Snapshot
==============================
Operating System: %s
CPU Usage       : %.2f%%
Memory Usage    : %.2f%% (%s / %s)
Disk Usage      : %.2f%% (%s / %s)
%s
------------------------------
%s
`, s.OS,
		s.CPUPercent,
		s.MemoryUsedPct, s.MemoryUsed, s.MemoryTotal,
		s.DiskUsedPct, s.DiskUsed, s.DiskTotal,
		loadLine,
		summary,
	)
}

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
