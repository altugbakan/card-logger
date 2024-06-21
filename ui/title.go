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
		case "a", "A":
			return NewAddModel(), nil
		case "l", "L":
			return NewListModel(), nil
		}
	}

	return m, nil
}

func (m TitleModel) View() string {
	header := HeaderStyle.PaddingBottom(2).Render("Card Logger")
	options := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Left, ActionStyle.Render("[A]"),
			TextStyle.PaddingRight(4).Render("dd Cards")),
		lipgloss.JoinHorizontal(lipgloss.Left, ActionStyle.Render("[L]"), "ist Cards"),
	)
	return lipgloss.JoinVertical(lipgloss.Center, header, options)
}
