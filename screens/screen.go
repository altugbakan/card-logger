package screens

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Screen interface {
	Update(tea.Msg) (Screen, tea.Cmd)
	View() string
	Help() string
}
