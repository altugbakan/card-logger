package ui

import tea "github.com/charmbracelet/bubbletea"

type AddModel struct{}

func NewAddModel() AddModel {
	return AddModel{}
}

func (m AddModel) Init() tea.Cmd {
	return nil
}

func (m AddModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m AddModel) View() string {
	return "Hello World from Add"
}
