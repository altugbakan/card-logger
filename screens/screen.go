package screens

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Screen interface {
	Process(string) (Screen, tea.Cmd)
	View() string
	KeyBindings() map[string]key.Binding
}
