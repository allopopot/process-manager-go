package main

import (
	"fmt"
	"os"

	"allopopot.com/proc-man/ui"
	tea "charm.land/bubbletea/v2"
)

func main() {
	p := tea.NewProgram(ui.AppInit())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
