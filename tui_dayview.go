package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type DayView struct {
	Date    time.Time
	Content string
}

func initialDayView() DayView {
	return DayView{
		Date:    time.Now(),
		Content: time.Now().Format("2006-01-02"),
	}
}

func (dv DayView) Init() tea.Cmd {
	return nil
}

func (dv DayView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return dv, tea.Quit
		}
	}

	return dv, nil
}

func (dv DayView) View() string {
	return fmt.Sprintf("%s\n", dv.Content)
}
