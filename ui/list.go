package ui

import tea "github.com/charmbracelet/bubbletea"

type ListModel struct{}

func NewListModel() ListModel {
	return ListModel{}
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			return NewTitleModel(), nil
		}
	}

	return m, nil
}

func (m ListModel) View() string {
	return "Hello World from List"
}
