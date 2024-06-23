package keymaps

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type Title struct {
	Add    key.Binding
	List   key.Binding
	Backup key.Binding
	Quit   key.Binding
	help   help.Model
}

func NewTitleKeyMap(bindings ...key.Binding) Title {
	return Title{
		Add: key.NewBinding(
			key.WithKeys("a", "A"),
			key.WithHelp("a", "add cards"),
		),
		List: key.NewBinding(
			key.WithKeys("l", "L"),
			key.WithHelp("l", "list cards"),
		),
		Backup: key.NewBinding(
			key.WithKeys("b", "B"),
			key.WithHelp("b", "backup"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "Q"),
			key.WithHelp("q", "quit"),
		),
		help: help.New(),
	}
}

func (k Title) Help() string {
	return k.help.View(k)
}

func (k Title) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Add,
		k.List,
		k.Backup,
		k.Quit,
	}
}

func (k Title) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Add, k.List, k.Backup, k.Quit},
	}
}
