// main.go
package main

import (
	"log"

	ui "github.com/altugbakan/card-logger/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	currentScreen tea.Model
	width         int
	height        int
}

func main() {
	initialModel := model{
		currentScreen: ui.NewTitleModel(),
		width:         60,
		height:        80,
	}

	p := tea.NewProgram(initialModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	var cmd tea.Cmd
	m.currentScreen, cmd = m.currentScreen.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		m.currentScreen.View())
}
