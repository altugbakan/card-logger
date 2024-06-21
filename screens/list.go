package screens

import (
	"github.com/altugbakan/card-logger/utils"
	tea "github.com/charmbracelet/bubbletea"
)

type ListScreen struct {
	keyMap utils.KeyMap
}

func NewListScreen() ListScreen {
	return ListScreen{}
}

func (s ListScreen) Update(msg tea.KeyMsg) (Screen, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return s, tea.Quit
	case "esc":
		return NewTitleModel(), nil
	}

	return s, nil
}

func (h ListScreen) View() string {
	return "Hello World from List"
}

func (s ListScreen) Help() string {
	return s.keyMap.Help()
}
