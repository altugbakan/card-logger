package ui

import tea "github.com/charmbracelet/bubbletea"

type Screen interface {
	tea.Model
	HandleCommand(string) (tea.Model, tea.Cmd)
}
