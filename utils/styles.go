package utils

import "github.com/charmbracelet/lipgloss"

const (
	normal    = lipgloss.Color("#dddddd")
	gray      = lipgloss.Color("#626262")
	lightBlue = lipgloss.Color("#add8e6")
	pastelRed = lipgloss.Color("#ff6666")
)

var (
	HeaderStyle = lipgloss.NewStyle().Bold(true).Align(lipgloss.Center).
			Background(lightBlue).Foreground(gray).MarginBottom(2)
	ActionStyle  = lipgloss.NewStyle().Foreground(lightBlue)
	DimTextStyle = lipgloss.NewStyle().Foreground(gray)
	TextStyle    = lipgloss.NewStyle().Foreground(normal)
	ErrorStyle   = lipgloss.NewStyle().Foreground(pastelRed)
)