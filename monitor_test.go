package main

import (
	"strings"
	"testing"
)

// TestEvaluateThresholds checks the warning messages based on thresholds.
func TestEvaluateThresholds(t *testing.T) {
	cpuThreshold = 80
	memThreshold = 90
	diskThreshold = 90

	tests := []struct {
		name     string
		stats    Stats
		expected string
	}{
		{
			name: "All below thresholds",
			stats: Stats{
				CPUPercent:    30.0,
				MemoryUsedPct: 40.0,
				DiskUsedPct:   50.0,
			},
			expected: "All systems normal",
		},
		{
			name: "High CPU",
			stats: Stats{
				CPUPercent:    95.0,
				MemoryUsedPct: 40.0,
				DiskUsedPct:   50.0,
			},
			expected: "High CPU usage",
		},
		{
			name: "High Mem",
			stats: Stats{
				CPUPercent:    50.0,
				MemoryUsedPct: 95.0,
				DiskUsedPct:   50.0,
			},
			expected: "High memory usage",
		},
		{
			name: "High Disk",
			stats: Stats{
				CPUPercent:    50.0,
				MemoryUsedPct: 50.0,
				DiskUsedPct:   95.0,
			},
			expected: "High disk usage",
		},
		{
			name: "Multiple alerts",
			stats: Stats{
				CPUPercent:    91.0,
				MemoryUsedPct: 91.0,
				DiskUsedPct:   91.0,
			},
			expected: "High CPU usage",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := EvaluateThresholds(test.stats)
			if !strings.Contains(result, test.expected) {
				t.Errorf("Expected to contain %q, got %q", test.expected, result)
			}
		})
	}
}

// TestFormatStats ensures readable output contains key values.
func TestFormatStats(t *testing.T) {
	stats := Stats{
		OS:            "linux",
		CPUPercent:    50.5,
		MemoryUsedPct: 60.6,
		DiskUsedPct:   70.7,
		Load1:         1.1,
		Load5:         1.2,
		Load15:        1.3,
	}

	output := FormatStats(stats, "Test summary")

	expectedKeywords := []string{
		"linux",
		"50.5", "60.6", "70.7",
		"1.1", "1.2", "1.3",
		"Test summary",
	}

	for _, keyword := range expectedKeywords {
		if !strings.Contains(output, keyword) {
			t.Errorf("Expected output to contain %q", keyword)
		}
	}
}
