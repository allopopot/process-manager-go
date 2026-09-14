package views

import (
	"fmt"
	"strings"
	"time"

	"allopopot.com/proc-man/config"
	"charm.land/lipgloss/v2"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
)

func RenderHostInfo() string {
	return lipgloss.JoinHorizontal(lipgloss.Top, renderSysInfo(), renderLoadAveragesAndLoggedInUsers())

}

func renderSysInfo() string {
	layout := lipgloss.NewStyle().
		Width(70)

	s := ""

	s += lipgloss.NewStyle().Foreground(config.ColorAccent).Render("System Information")
	s += "\n"
	s += strings.Repeat("─", 50)
	s += "\n"

	h, _ := host.Info()

	s += kv("Hostname", h.Hostname)
	s += kv("Operating System", h.OS)
	s += kv("Platform", fmt.Sprintf("%s %s", h.Platform, h.PlatformVersion))
	s += kv("Kernel", h.KernelVersion)
	s += kv("Architecture", h.KernelArch)
	s += kv("Uptime", formatUptime(h.Uptime))
	s += kv("Boot Time", formatBootTime(h.BootTime))
	s += kv("Running Processes", fmt.Sprintf("%d", h.Procs))
	s += kv("Virtualization", fmt.Sprintf("%s (%s)", h.VirtualizationSystem, h.VirtualizationRole))

	return layout.Render(s)

}

func renderLoadAveragesAndLoggedInUsers() string {
	layout := lipgloss.JoinVertical(lipgloss.Left, renderLoggedinUsers(), renderLoadAverages())
	return layout
}

func renderLoggedinUsers() string {
	layout := lipgloss.NewStyle().
		Width(50)
	s := ""

	s += lipgloss.NewStyle().Foreground(config.ColorAccent).Render("Logged-In Users")
	s += "\n"
	s += strings.Repeat("─", 50)
	s += "\n"

	users, _ := host.Users()

	s += lipgloss.NewStyle().Foreground(config.ColorPrimary).Bold(true).Render(fmt.Sprintf("%-14s %-12s %s", "User", "Terminal", "Host"))
	s += "\n"

	for _, u := range users {
		s += fmt.Sprintf("%-14s %-12s %s", u.User, u.Terminal, u.Host)
		s += "\n"
	}

	return layout.Render(s)

}

func renderLoadAverages() string {
	layout := lipgloss.NewStyle().
		Width(50)
	s := ""

	s += lipgloss.NewStyle().Foreground(config.ColorAccent).Render("Load Averages")
	s += "\n"
	s += strings.Repeat("─", 50)
	s += "\n"

	ld, _ := load.Avg()

	s += kv("1 minute", fmt.Sprintf("%.2f", ld.Load1))
	s += kv("5 minutes", fmt.Sprintf("%.2f", ld.Load5))
	s += kv("15 minutes", fmt.Sprintf("%.2f", ld.Load15))

	return layout.Render(s)

}

func formatUptime(seconds uint64) string {
	d := time.Duration(seconds) * time.Second
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	secs := int(d.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %02dh %02dm %02ds", days, hours, minutes, secs)
	}
	return fmt.Sprintf("%02dh %02dm %02ds", hours, minutes, secs)
}

func formatBootTime(ts uint64) string {
	t := time.Unix(int64(ts), 0)
	return t.Format("Mon, 02 Jan 2006  15:04:05")
}

func kv(key string, val string) string {
	key = lipgloss.NewStyle().Width(22).Render(key + ":")
	val = lipgloss.NewStyle().Render(val)
	return key + val + "\n"
}
