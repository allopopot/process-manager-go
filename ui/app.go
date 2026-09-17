package ui

import (
	"fmt"
	"image/color"
	"strconv"
	"time"

	"allopopot.com/proc-man/config"
	"allopopot.com/proc-man/ui/views"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
)

type App struct {
	activeTab int
	tabList   []string
	width     int
	height    int
}

type tickMsg struct{}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m App) Init() tea.Cmd {
	return tick()
}

func (m App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		return m, tick()

	case tea.KeyPressMsg:

		switch msg.String() {

		case "1", "2", "3", "4", "5", "6":
			n, _ := strconv.Atoi(msg.String())
			if n >= 0 && n < len(m.tabList)+1 {
				m.activeTab = n - 1
			}

		case "ctrl+c", "q":
			return m, tea.Quit

		}

	}

	return m, nil
}

func (m App) View() tea.View {

	s := ""
	s += lipgloss.NewStyle().SetString("  Process Manager (Proc Man)  ").Background(config.ColorAccent).Render()
	s += "\n"

	for i, x := range m.tabList {
		var a color.Color
		if i == m.activeTab {
			a = config.ColorPrimary
		} else {
			a = config.ColorBase
		}
		s += lipgloss.NewStyle().Background(a).Render(fmt.Sprintf("  %v %s  ", i+1, x))
	}
	s += "\n\n"

	switch m.activeTab {
	case 0:
		s += views.RenderHostInfo()
	case 1:
		s += views.RenderStorageInfo()
	case 2:
		s += views.RenderCPUInfo()
	case 3:
		s += views.RenderMemoryInfo()
	}

	return tea.NewView(
		lipgloss.NewStyle().
			Width(m.width).
			Height(m.height).
			Render(s),
	)

}

func AppInit() App {
	return App{
		activeTab: 3,
		tabList:   []string{"Host Info", "Storage", "CPU", "Memory", "Network", "Tasks"},
	}
}
