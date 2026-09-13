package ui

import (
	"fmt"
	"image/color"
	"strconv"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
)

type App struct {
	activeTab int
	tabList   []string
}

var (
	ColorGood     = lipgloss.Color("#00e887")
	ColorWarning  = lipgloss.Color("#ffd700")
	ColorCritical = lipgloss.Color("#ff4444")
	ColorPrimary  = lipgloss.Color("#5f87ff")
	ColorAccent   = lipgloss.Color("#00bcd4")
	ColorMuted    = lipgloss.Color("#6c7086")
	ColorText     = lipgloss.Color("#cdd6f4")
	ColorSubtext  = lipgloss.Color("#a6adc8")
	ColorBorder   = lipgloss.Color("#45475a")
	ColorBase     = lipgloss.Color("#1e1e2e")
)

func (m App) Init() tea.Cmd {
	return nil
}

func (m App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyPressMsg:

		switch msg.String() {

		case "1", "2", "3", "4", "5", "6":
			n, _ := strconv.Atoi(msg.String())
			m.activeTab = n - 1

		case "ctrl+c", "q":
			return m, tea.Quit

		}
	}

	return m, nil
}

func (m App) View() tea.View {

	s := ""
	s += lipgloss.NewStyle().SetString("  Process Manager (Proc Man)  ").Background(lipgloss.Green).Render()
	s += "\n"

	for i, x := range m.tabList {
		var a color.Color
		if i == m.activeTab {
			a = ColorGood
		} else {
			a = ColorBase
		}
		s += lipgloss.NewStyle().Background(a).Render(fmt.Sprintf("  %v %s  ", i+1, x))
	}

	return tea.NewView(s)

}

func AppInit() App {
	return App{
		activeTab: 0,
		tabList:   []string{"Host Info", "Storage", "CPU", "Memory", "Network", "Tasks"},
	}
}
