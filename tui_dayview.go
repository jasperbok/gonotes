package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/glamour"

	tea "github.com/charmbracelet/bubbletea"
)

type DayView struct {
	Date    time.Time
	Content string
	help    help.Model
}

func initialDayView() DayView {
	return DayView{
		Date:    time.Now(),
		Content: time.Now().Format("2006-01-02"),
		help:    help.New(),
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
	case EditorFinishedMsg:
		if msg.err != nil {
			return dv, tea.Quit
		}
		dv.ReloadContents()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultDayViewKeyMap.Quit):
			return dv, tea.Quit
		case key.Matches(msg, DefaultDayViewKeyMap.Prev):
			return dv.GoToPreviousDay(), nil
		case key.Matches(msg, DefaultDayViewKeyMap.Edit):
			filePath, err := dv.FilePath()
			if err != nil {
				return dv, tea.Quit
			}
			return dv, openEditor(filePath)
		case key.Matches(msg, DefaultDayViewKeyMap.Next):
			return dv.GoToNextDay(), nil
		}
	}

	return dv, nil
}

func (dv DayView) View() string {
	s := fmt.Sprintf("\n# %s\n\n", dv.Date.Format("2006-01-02"))

	contents, err := glamour.Render(dv.Content, "dark")
	if err != nil {
		s = fmt.Sprintf("%s%s\n", s, err.Error())
	} else {
		s = fmt.Sprintf("%s%s\n", s, contents)
	}
	s = fmt.Sprintf("%s\n\n%s\n", s, dv.help.View(DefaultDayViewKeyMap))
	return s
}

// FilePath returns the absolute path to the markdown file for DayView's Date.
func (dv DayView) FilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/gonotes/%s.md", homeDir, dv.Date.Format("2006-01-02")), nil
}

// ReloadContents re-reads the DayView's contents from disk.
func (dv *DayView) ReloadContents() {
	msg := getContentsForDate(dv.Date)()
	switch msg := msg.(type) {
	case ErrMsg:
		dv.Content = msg.err.Error()
	case FileLoadedMsg:
		dv.Content = msg.content
	}
}

func (dv DayView) GoToDate(date time.Time) DayView {
	dv.Date = date

	msg := getContentsForDate(date)()
	switch msg := msg.(type) {
	case ErrMsg:
		dv.Content = msg.err.Error()
	case FileLoadedMsg:
		dv.Content = msg.content
	}

	return dv
}

func (dv DayView) GoToPreviousDay() DayView {
	yesterday := dv.Date.Add(time.Hour * 24 * -1)
	return dv.GoToDate(yesterday)
}

func (dv DayView) GoToNextDay() DayView {
	tomorrow := dv.Date.Add(time.Hour * 24)
	return dv.GoToDate(tomorrow)
}

type ErrMsg struct {
	err error
}

type FileLoadedMsg struct {
	content string
}

type EditorFinishedMsg struct {
	err error
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

func openEditor(filePath string) tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	c := exec.Command(editor, filePath)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return EditorFinishedMsg{err}
	})
}

type DayViewKeyMap struct {
	Next key.Binding
	Prev key.Binding
	Edit key.Binding
	Quit key.Binding
}

func (k DayViewKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Prev, k.Next, k.Edit, k.Quit}
}

func (k DayViewKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Prev, k.Next, k.Edit, k.Quit},
	}
}

var DefaultDayViewKeyMap = DayViewKeyMap{
	Next: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("l", "next day"),
	),
	Prev: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("h", "prev day"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("q", "quit"),
	),
}
