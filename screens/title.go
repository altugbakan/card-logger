package screens

import (
	"github.com/altugbakan/card-logger/keymaps"
	"github.com/altugbakan/card-logger/utils"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Title struct {
	keyMap keymaps.Title
}

func NewTitleScreen() Title {
	keyMap := keymaps.NewTitleKeyMap()

	return Title{keyMap: keyMap}
}

func (s Title) Update(msg tea.Msg) (Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keyMap.Quit):
			utils.LogInfo("quitting the program...")
			return s, tea.Quit
		case key.Matches(msg, s.keyMap.Add):
			return NewAddScreen(), textinput.Blink
		case key.Matches(msg, s.keyMap.List):
			listScreen, err := NewListScreen()
			if err != nil {
				return s, tea.Quit
			}
			return listScreen, nil
		case key.Matches(msg, s.keyMap.Backup):
			return NewBackupScreen(), nil
		}
	}

	return s, nil
}

func (s Title) View() string {
	title := utils.TitleStyle.Render("Card Logger")
	options := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Left, utils.ActionStyle.Render("[A]"), "dd Cards"),
		lipgloss.JoinHorizontal(lipgloss.Left, utils.ActionStyle.MarginLeft(4).Render("[L]"), "ist Cards"),
	)
	backup := lipgloss.JoinHorizontal(lipgloss.Left, utils.ActionStyle.Render("[B]"), "ackup")
	backup = utils.EmptyStyle.MarginTop(1).Render(backup)
	options = lipgloss.JoinVertical(lipgloss.Center, options, backup)
	return lipgloss.JoinVertical(lipgloss.Center, title, options)
}

func (s Title) Help() string {
	return s.keyMap.Help()
}
