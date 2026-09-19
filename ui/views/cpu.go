package views

import (
	"fmt"
	"strings"
	"time"

	"allopopot.com/proc-man/config"
	"charm.land/lipgloss/v2"
	"github.com/shirou/gopsutil/cpu"
)

var (
	lastOverall []float64
	lastPerCore []float64
	lastCPUTime time.Time
)

func RenderCPUInfo() string {

	var sb strings.Builder

	sb.WriteString(lipgloss.NewStyle().Foreground(config.ColorAccent).Render("CPU Information"))
	sb.WriteString("\n")

	info, _ := cpu.Info()

	if len(info) > 0 && info[0].ModelName != "" {
		pp, _ := cpu.Counts(false)
		lp, _ := cpu.Counts(true)
		meta := fmt.Sprintf("  %s   Physical: %d   Logical: %d   Base: %.0f MHz",
			info[0].ModelName,
			pp,
			lp,
			info[0].Mhz,
		)
		sb.WriteString(lipgloss.NewStyle().Foreground(config.ColorMuted).Render(meta))
		sb.WriteString("\n")
	}
	sb.WriteString("\n")

	now := time.Now()
	if lastOverall == nil || now.Sub(lastCPUTime) >= 500*time.Millisecond {
		lastOverall, _ = cpu.Percent(0, false)
		lastPerCore, _ = cpu.Percent(0, true)
		lastCPUTime = now
	}

	if len(lastOverall) > 0 {
		sb.WriteString(fmt.Sprintf("Total %-2s", ""))
		sb.WriteString(storageBar(lastOverall[0], 60))
		sb.WriteString("\n")
		sb.WriteString(strings.Repeat("─", 80))
		sb.WriteString("\n")
	}

	for i, px := range lastPerCore {
		sb.WriteString(fmt.Sprintf("Core %d %-1s", i, ""))
		sb.WriteString(storageBar(px, 60))
		sb.WriteString("\n")
	}

	return sb.String()
}
