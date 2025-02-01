package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(initialDayView(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("The following error occured: %v\n", err)
		os.Exit(1)
	}
}
