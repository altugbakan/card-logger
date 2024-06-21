// main.go
package main

import (
	"log"

	"github.com/altugbakan/card-logger/screens"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	currentScreen screens.Screen
	width         int
	height        int
}

func main() {
	initialModel := model{
		currentScreen: screens.NewTitleModel(),
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
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		input := msg.String()
		switch input {
		case "ctrl+c":
			return m, tea.Quit
		}
		m.currentScreen, cmd = m.currentScreen.Process(input)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, cmd
}

func (m model) View() string {
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		m.currentScreen.View())
}
