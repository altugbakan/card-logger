package screens

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type AddModel struct{}

func NewAddModel() AddModel {
	return AddModel{}
}

func (m AddModel) Init() tea.Cmd {
	return nil
}

func (m AddModel) Process(input string) (Screen, tea.Cmd) {
	switch input {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "esc":
		return NewTitleModel(), nil
	}

	return m, nil
}

func (m AddModel) View() string {
	return "Hello World from Add"
}

func (m AddModel) KeyBindings() map[string]key.Binding {
	return make(map[string]key.Binding)
}
