// main.go
package main

import (
	"fmt"
	"os"

	"github.com/altugbakan/card-logger/db"
	"github.com/altugbakan/card-logger/screens"
	"github.com/altugbakan/card-logger/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
)

type model struct {
	currentScreen screens.Screen
	width         int
	height        int
}

func main() {
	db.InitDB()
	defer db.CloseDB()

	if os.Getenv("DEBUG") == "true" {
		f, err := tea.LogToFile("debug.log", "[debug]")
		if err != nil {
			fmt.Printf("could not log to file: %v", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	width, height, err := term.GetSize(0)
	if err != nil {
		utils.LogWarning("could not get terminal size, using defaults: %v", err)
		width = utils.DefaultWidth
		height = utils.DefaultHeight
	}
	utils.LogInfo("initial terminal size: %d x %d", width, height)

	initialModel := model{
		currentScreen: screens.NewTitleScreen(),
		width:         width,
		height:        height,
	}

	utils.LogInfo("starting the program...")
	p := tea.NewProgram(initialModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		utils.LogError("could not run the program: %v", err)
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
			utils.LogInfo("ctrl+c detected, exiting the program...")
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Remove help width from height
		msg.Height -= utils.TotalHelpWidth
	}

	m.currentScreen, cmd = m.currentScreen.Update(msg)
	return m, cmd
}

func (m model) View() string {
	help := utils.EmptyStyle.Margin(utils.HelpMargin).Render(m.currentScreen.Help())
	view := lipgloss.Place(m.width, m.height-lipgloss.Height(help), lipgloss.Center, lipgloss.Center,
		m.currentScreen.View())

	return lipgloss.JoinVertical(lipgloss.Left, view, help)
}
