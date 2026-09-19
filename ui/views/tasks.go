package views

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"allopopot.com/proc-man/config"
	"charm.land/bubbles/v2/table"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/shirou/gopsutil/v4/process"
)

type ProcessInfo struct {
	Pid    int32
	Name   string
	Status string
	CPU    float64
	Mem    float32
	RSS    uint64
}

type TasksModel struct {
	table     table.Model
	search    textinput.Model
	procs     []ProcessInfo
	sortCol   string
	sortDesc  bool
	searching bool
	width     int
	height    int
}

func NewTasksModel() TasksModel {
	ti := textinput.New()
	ti.Placeholder = "Search by name..."
	ti.CharLimit = 50
	ti.SetWidth(30)

	cols := []table.Column{
		{Title: "PID", Width: 8},
		{Title: "Name", Width: 32},
		{Title: "Status", Width: 10},
		{Title: "CPU%", Width: 8},
		{Title: "MEM%", Width: 8},
		{Title: "RSS", Width: 12},
	}

	tbl := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(config.ColorBorder).
		BorderBottom(true).
		Foreground(config.ColorPrimary).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#1e1e2e")).
		Background(config.ColorPrimary).
		Bold(true)
	s.Cell = s.Cell.Foreground(config.ColorText)
	tbl.SetStyles(s)

	return TasksModel{
		table:    tbl,
		search:   ti,
		sortCol:  "CPU%",
		sortDesc: true,
	}
}

type TasksUpdateMsg []ProcessInfo

func FetchTasks() tea.Msg {
	procs, err := process.Processes()
	if err != nil {
		return TasksUpdateMsg{}
	}

	var infos []ProcessInfo
	for _, p := range procs {
		name, _ := p.Name()
		status, _ := p.Status()
		st := ""
		if len(status) > 0 {
			st = status[0]
		}
		cpu, _ := p.CPUPercent()
		mem, _ := p.MemoryPercent()
		memInfo, _ := p.MemoryInfo()
		rss := uint64(0)
		if memInfo != nil {
			rss = memInfo.RSS
		}

		infos = append(infos, ProcessInfo{
			Pid:    p.Pid,
			Name:   name,
			Status: st,
			CPU:    cpu,
			Mem:    mem,
			RSS:    rss,
		})
	}
	return TasksUpdateMsg(infos)
}

func (m TasksModel) SetSize(w, h int) TasksModel {
	m.width = w
	m.height = h
	m.table.SetWidth(w)
	m.table.SetHeight(h - 10)
	return m
}

func formatSize(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func (m *TasksModel) updateTable() {
	var filtered []ProcessInfo
	term := strings.ToLower(m.search.Value())

	for _, p := range m.procs {
		if term == "" || strings.Contains(strings.ToLower(p.Name), term) {
			filtered = append(filtered, p)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		a := filtered[i]
		b := filtered[j]
		var less bool

		switch m.sortCol {
		case "PID":
			less = a.Pid < b.Pid
		case "Name":
			less = strings.ToLower(a.Name) < strings.ToLower(b.Name)
		case "Status":
			less = a.Status < b.Status
		case "CPU%":
			less = a.CPU < b.CPU
		case "MEM%":
			less = a.Mem < b.Mem
		case "RSS":
			less = a.RSS < b.RSS
		default:
			less = a.Pid < b.Pid
		}

		if m.sortDesc {
			return !less
		}
		return less
	})

	var rows []table.Row
	for _, p := range filtered {
		rows = append(rows, table.Row{
			strconv.Itoa(int(p.Pid)),
			p.Name,
			p.Status,
			fmt.Sprintf("%.1f", p.CPU),
			fmt.Sprintf("%.1f", p.Mem),
			formatSize(p.RSS),
		})
	}

	// Try to maintain selected row
	cursor := m.table.Cursor()
	m.table.SetRows(rows)

	if cursor >= len(rows) {
		cursor = len(rows) - 1
		if cursor < 0 {
			cursor = 0
		}
	}
	m.table.SetCursor(cursor)
}

func (m TasksModel) IsSearching() bool {
	return m.searching
}

func (m TasksModel) Update(msg tea.Msg) (TasksModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case TasksUpdateMsg:
		m.procs = msg
		m.updateTable()
		return m, nil

	case tea.KeyPressMsg:
		if m.searching {
			switch msg.String() {
			case "enter", "esc":
				m.searching = false
				m.search.Blur()
			default:
				m.search, cmd = m.search.Update(msg)
				m.updateTable()
				cmds = append(cmds, cmd)
			}
			return m, tea.Batch(cmds...)
		}

		switch msg.String() {
		case "/":
			m.searching = true
			return m, m.search.Focus()
		case "k":
			row := m.table.SelectedRow()
			if row != nil {
				pid, _ := strconv.Atoi(row[0])
				p, err := process.NewProcess(int32(pid))
				if err == nil {
					_ = p.Kill()
				}
			}
		case "c":
			m.sortCol = "CPU%"
			m.sortDesc = true
			m.updateTable()
		case "m":
			m.sortCol = "MEM%"
			m.sortDesc = true
			m.updateTable()
		case "p":
			m.sortCol = "PID"
			m.sortDesc = false
			m.updateTable()
		case "n":
			m.sortCol = "Name"
			m.sortDesc = false
			m.updateTable()
		case "s":
			m.sortDesc = !m.sortDesc
			m.updateTable()
		}
	}

	m.search, cmd = m.search.Update(msg)
	cmds = append(cmds, cmd)

	if !m.searching {
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m TasksModel) View() string {
	s := ""

	s += lipgloss.NewStyle().Foreground(config.ColorAccent).Render("Task Manager")
	s += "\n"
	s += mutedStyle.Render(strings.Repeat("─", 100))
	s += "\n"

	headerInfo := " [/] Search  [k] Kill  [c] Sort CPU  [m] Sort Mem  [p] Sort PID  [n] Sort Name  [s] Toggle ASC/DESC"
	s += lipgloss.NewStyle().Foreground(config.ColorText).Render(headerInfo)
	s += "\n"
	s += mutedStyle.Render(strings.Repeat("─", 100))
	s += "\n"

	if m.searching {
		s += m.search.View() + "\n\n"
	} else {
		s += "\n\n"
	}

	s += m.table.View()

	return s
}
