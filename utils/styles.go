package utils

import "github.com/charmbracelet/lipgloss"

const (
	Normal          = "#dddddd"
	Gray            = "#626262"
	LightBlue       = "#add8e6"
	DarkLightBlue   = "#87ceeb"
	PastelRed       = "#ff6666"
	DarkYellow      = "#ffcc00"
	normalColor     = lipgloss.Color(Normal)
	grayColor       = lipgloss.Color(Gray)
	lightBlueColor  = lipgloss.Color(LightBlue)
	pastelRedColor  = lipgloss.Color(PastelRed)
	darkYellowColor = lipgloss.Color(DarkYellow)
)

var (
	TitleStyle = lipgloss.NewStyle().Bold(true).Align(lipgloss.Center).
			Background(lightBlueColor).Foreground(grayColor).MarginBottom(2)
	ActionStyle  = lipgloss.NewStyle().Foreground(lightBlueColor)
	DimTextStyle = lipgloss.NewStyle().Foreground(grayColor)
	TextStyle    = lipgloss.NewStyle().Foreground(normalColor)
	WarningStyle = lipgloss.NewStyle().Foreground(darkYellowColor)
	ErrorStyle   = lipgloss.NewStyle().Foreground(pastelRedColor)
	CursorStyle  = lipgloss.NewStyle().Background(lightBlueColor)
	EmptyStyle   = lipgloss.NewStyle()
)
