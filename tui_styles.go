package main

import (
	"github.com/charmbracelet/lipgloss"
)

var ContentHeaderStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#999999")).
	PaddingTop(1).
	PaddingBottom(2).
	Align(lipgloss.Center)
