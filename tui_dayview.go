package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
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
	return getContentsForDate(dv.Date)
}

func (dv DayView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ErrMsg:
		dv.Content = msg.err.Error()
	case FileLoadedMsg:
		dv.Content = msg.content
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return dv, tea.Quit
		}
	}

	return dv, nil
}

func (dv DayView) View() string {
	s := fmt.Sprintf("\n# %s\n\n", dv.Date.Format("2006-01-02"))
	s = fmt.Sprintf("%s%s\n", s, dv.Content)
	return s
}
}

type ErrMsg struct {
	err error
}

type FileLoadedMsg struct {
	content string
}

func getContentsForDate(date time.Time) tea.Cmd {
	return func() tea.Msg {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return ErrMsg{err}
		}

		filename := fmt.Sprintf("%s/gonotes/%s.md", homeDir, date.Format("2006-01-02"))
		filePath, err := filepath.Abs(filename)
		if err != nil {
			return ErrMsg{err}
		}

		_, err = os.Stat(filePath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				// If the file does not yet exists, we don't want to create it right now.
				// If the user is browsing, creating any missing files could lead to a
				// directory full of meaningless, empty files.
				return FileLoadedMsg{""}
			}

			return ErrMsg{err}
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			return ErrMsg{err}
		}

		return FileLoadedMsg{string(content)}
	}
}
