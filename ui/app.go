package ui

import (
	"fmt"
	"image/color"
	"strconv"

	"allopopot.com/proc-man/config"
	"allopopot.com/proc-man/ui/views"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
)

type App struct {
	activeTab int
	tabList   []string
}

func (m App) Init() tea.Cmd {
	return nil
}

func (m App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

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
	}

	return tea.NewView(s)

}

func AppInit() App {
	return App{
		activeTab: 1,
		tabList:   []string{"Host Info", "Storage", "CPU", "Memory", "Network", "Tasks"},
	}
}
