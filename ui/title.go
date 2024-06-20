// ui/title.go
package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TitleModel struct{}

func NewTitleModel() TitleModel {
	return TitleModel{}
}

func (m TitleModel) Init() tea.Cmd {
	return nil
}

func (m TitleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "a":
			return NewAddModel(), nil
		case "l":
			return NewListModel(), nil
		}
	}

	return m, nil
}

func (m TitleModel) View() string {
	header := HeaderStyle.Render("Card Logger")
	options := lipgloss.JoinVertical(lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Left, ActionStyle.Render("[A]"), "dd Cards"),
		lipgloss.JoinHorizontal(lipgloss.Left, ActionStyle.Render("[L]"), "ist Cards"),
	)
	return lipgloss.JoinVertical(lipgloss.Center, header, options)
}
