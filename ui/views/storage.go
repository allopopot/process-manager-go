package views

import (
	"fmt"
	"strings"

	"allopopot.com/proc-man/config"
	"charm.land/lipgloss/v2"
	"github.com/shirou/gopsutil/disk"
)

// column widths (characters)
const (
	colDevice     = 22
	colMountpoint = 16
	colFS         = 10
	colTotal      = 11
	colUsed       = 11
	colFree       = 11
	colUsage      = 29
)

var (
	headerStyle = lipgloss.NewStyle().Foreground(config.ColorPrimary).Bold(true)
	mutedStyle  = lipgloss.NewStyle().Foreground(config.ColorMuted)
	textStyle   = lipgloss.NewStyle().Foreground(config.ColorText)
)

func RenderStorageInfo() string {
	partitions, _ := disk.Partitions(false)

	s := ""

	// Header
	s += lipgloss.NewStyle().Foreground(config.ColorAccent).Render("Storage")
	s += "\n"
	s += mutedStyle.Render(strings.Repeat("─", colDevice+colMountpoint+colFS+colTotal+colUsed+colFree+colUsage))
	s += "\n"
	s += headerStyle.Render(
		fmt.Sprintf("%-*s%-*s%-*s%-*s%-*s%-*s%s",
			colDevice, "Device",
			colMountpoint, "Mount",
			colFS, "FS",
			colTotal, "Total",
			colUsed, "Used",
			colFree, "Free",
			"Usage",
		),
	)
	s += "\n"
	s += mutedStyle.Render(strings.Repeat("─", colDevice+colMountpoint+colFS+colTotal+colUsed+colFree+colUsage))
	s += "\n"

	// Rows
	for _, p := range partitions {
		usage, err := disk.Usage(p.Mountpoint)
		if err != nil {
			continue
		}
		row := fmt.Sprintf("%-*s%-*s%-*s%-*s%-*s%-*s",
			colDevice, truncate(p.Device, colDevice-1),
			colMountpoint, truncate(p.Mountpoint, colMountpoint-1),
			colFS, truncate(p.Fstype, colFS-1),
			colTotal, formatGiB(usage.Total),
			colUsed, formatGiB(usage.Used),
			colFree, formatGiB(usage.Free),
		)
		s += textStyle.Render(row) + storageBar(usage.UsedPercent, 20)
		s += "\n"
	}

	return s
}

// truncate shortens a string to max chars, adding "…" if needed.
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}

// formatGiB formats bytes as a GiB string.
func formatGiB(bytes uint64) string {
	gib := float64(bytes) / (1 << 30)
	return fmt.Sprintf("%.1f GiB", gib)
}

// storageBar renders a fixed-width progress bar with a percentage label.
func storageBar(pct float64, width int) string {
	filled := int(pct / 100.0 * float64(width))
	if filled > width {
		filled = width
	}

	barColor := config.ColorGood
	switch {
	case pct >= 90:
		barColor = config.ColorCritical
	case pct >= 70:
		barColor = config.ColorWarning
	}

	bar := lipgloss.NewStyle().Foreground(barColor).Render(strings.Repeat("█", filled)) +
		lipgloss.NewStyle().Foreground(config.ColorMuted).Render(strings.Repeat("░", width-filled))

	return fmt.Sprintf("%s %5.1f%%", bar, pct)
}
