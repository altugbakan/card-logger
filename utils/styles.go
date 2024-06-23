package utils

import "github.com/charmbracelet/lipgloss"

const (
	Normal         = "#dddddd"
	Gray           = "#626262"
	LightBlue      = "#add8e6"
	DarkLightBlue  = "#87ceeb"
	PastelRed      = "#ff6666"
	normalColor    = lipgloss.Color(Normal)
	grayColor      = lipgloss.Color(Gray)
	lightBlueColor = lipgloss.Color(LightBlue)
	pastelRedColor = lipgloss.Color(PastelRed)
)

var (
	HeaderStyle = lipgloss.NewStyle().Bold(true).Align(lipgloss.Center).
			Background(lightBlueColor).Foreground(grayColor).MarginBottom(2)
	ActionStyle  = lipgloss.NewStyle().Foreground(lightBlueColor)
	DimTextStyle = lipgloss.NewStyle().Foreground(grayColor)
	TextStyle    = lipgloss.NewStyle().Foreground(normalColor)
	ErrorStyle   = lipgloss.NewStyle().Foreground(pastelRedColor)
	CursorStyle  = lipgloss.NewStyle().Background(lightBlueColor)
	EmptyStyle   = lipgloss.NewStyle()
)
