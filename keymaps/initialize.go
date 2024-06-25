package keymaps

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type Initialize struct {
	Quit     key.Binding
	Load     key.Binding
	Download key.Binding
	help     help.Model
}

func NewInitializeKeyMap() Initialize {
	return Initialize{
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
		Load: key.NewBinding(
			key.WithKeys("l"),
			key.WithHelp("l", "load backup"),
		),
		Download: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "download database"),
		),
		help: help.New(),
	}
}

func (k Initialize) Help() string {
	return k.help.View(k)
}

func (k Initialize) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Quit,
		k.Load,
		k.Download,
	}
}

func (k Initialize) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit, k.Load, k.Download},
	}
}
