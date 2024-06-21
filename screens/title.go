package screens

import (
	"github.com/altugbakan/card-logger/utils"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TitleScreen struct {
	keyMap utils.KeyMap
}

func NewTitleModel() TitleScreen {
	keyMap := utils.NewKeyMap(
		key.NewBinding(
			key.WithKeys("a", "A"),
			key.WithHelp("a/A", "Add cards"),
		),
		key.NewBinding(
			key.WithKeys("l", "L"),
			key.WithHelp("l/L", "List cards"),
		),
		key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "Quit"),
		),
	)

	return TitleScreen{keyMap: keyMap}
}

func (s TitleScreen) Update(msg tea.KeyMsg) (Screen, tea.Cmd) {

	input := msg.String()
	switch input {
	case "q":
		return s, tea.Quit
	case "a", "A":
		return NewAddScreen(), textinput.Blink
	case "l", "L":
		return NewListScreen(), nil
	}

	return s, nil
}

func (s TitleScreen) View() string {
	header := utils.HeaderStyle.Render("Card Logger")
	options := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Left, utils.ActionStyle.Render("[A]"),
			utils.TextStyle.PaddingRight(4).Render("dd Cards")),
		lipgloss.JoinHorizontal(lipgloss.Left, utils.ActionStyle.Render("[L]"), "ist Cards"),
	)
	return lipgloss.JoinVertical(lipgloss.Center, header, options)
}

func (s TitleScreen) Help() string {
	return s.keyMap.Help()
}
