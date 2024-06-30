package displays

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Displayer interface {
	Update(tea.Msg) (Displayer, tea.Cmd)
	View() string
	Help() string
}
