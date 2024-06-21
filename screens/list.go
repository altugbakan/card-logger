package screens

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type ListModel struct{}

func NewListModel() ListModel {
	return ListModel{}
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m ListModel) Process(input string) (Screen, tea.Cmd) {

	switch input {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "esc":
		return NewTitleModel(), nil
	}

	return m, nil
}

func (m ListModel) View() string {
	return "Hello World from List"
}

func (m ListModel) KeyBindings() map[string]key.Binding {
	return make(map[string]key.Binding)
}
