package ui

import "github.com/charmbracelet/lipgloss"

const (
	normal    = lipgloss.Color("#dddddd")
	gray      = lipgloss.Color("#626262")
	lightBlue = lipgloss.Color("#add8e6")
)

var (
	HeaderStyle  = lipgloss.NewStyle().Bold(true).Align(lipgloss.Center).Foreground(normal)
	ActionStyle  = lipgloss.NewStyle().Foreground(lightBlue)
	DimTextStyle = lipgloss.NewStyle().Foreground(gray)
	TextStyle    = lipgloss.NewStyle().Foreground(normal)
)
