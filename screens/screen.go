package screens

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Screen interface {
	Update(tea.KeyMsg) (Screen, tea.Cmd)
	View() string
	Help() string
}
