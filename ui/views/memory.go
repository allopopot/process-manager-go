package views

import (
	"fmt"
	"strings"

	"allopopot.com/proc-man/config"
	"charm.land/lipgloss/v2"
	"github.com/shirou/gopsutil/v4/mem"
)

func RenderMemoryInfo() string {

	vm, _ := mem.VirtualMemory()

	s := ""

	s += lipgloss.NewStyle().Foreground(config.ColorAccent).Render("Memory")
	s += "\n"
	s += mutedStyle.Render(strings.Repeat("─", 50))
	s += "\n"

	s += storageBar(vm.UsedPercent, 50)
	s += "\n\n"

	s += fmt.Sprintf("%s: %s  ", headerStyle.Render("Total"), mutedStyle.Render(formatBytes(vm.Total)))
	s += fmt.Sprintf("%s: %s  ", headerStyle.Render("Used"), mutedStyle.Render(formatBytes(vm.Used)))
	s += fmt.Sprintf("%s: %s  ", headerStyle.Render("Free"), mutedStyle.Render(formatBytes(vm.Free)))
	s += fmt.Sprintf("%s: %s  ", headerStyle.Render("Available"), mutedStyle.Render(formatBytes(vm.Available)))
	s += "\n"
	s += fmt.Sprintf("%s: %s  ", headerStyle.Render("Cached"), mutedStyle.Render(formatBytes(vm.Cached)))
	s += fmt.Sprintf("%s: %s  ", headerStyle.Render("Buffers"), mutedStyle.Render(formatBytes(vm.Buffers)))
	s += "\n"

	return s
}
