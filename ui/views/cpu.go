package views

import (
	"fmt"
	"strings"

	"allopopot.com/proc-man/config"
	"charm.land/lipgloss/v2"
	"github.com/shirou/gopsutil/cpu"
)

func RenderCPUInfo() string {

	var sb strings.Builder

	sb.WriteString(lipgloss.NewStyle().Foreground(config.ColorAccent).Render("CPU Information"))
	sb.WriteString("\n")

	info, _ := cpu.Info()

	if info[0].ModelName != "" {
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

	overallCPUPercent, _ := cpu.Percent(0, false)
	sb.WriteString(fmt.Sprintf("Total %-2s", ""))
	sb.WriteString(storageBar(overallCPUPercent[0], 60))
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("─", 80))
	sb.WriteString("\n")

	perCoreCPUPercent, _ := cpu.Percent(0, true)
	for i, px := range perCoreCPUPercent {
		sb.WriteString(fmt.Sprintf("Core %d %-1s", i, ""))
		sb.WriteString(storageBar(px, 60))
		sb.WriteString("\n")
	}

	return sb.String()
}
