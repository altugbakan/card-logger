// main.go
package main

import (
	"log"

	ui "github.com/altugbakan/card-logger/ui"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	currentScreen tea.Model
}

func main() {
	initialModel := model{
		currentScreen: ui.NewTitleModel(),
	}

	p := tea.NewProgram(initialModel.currentScreen, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.currentScreen, cmd = m.currentScreen.Update(msg)
	return m.currentScreen, cmd
}

func (m model) View() string {
	return m.currentScreen.View()
}
