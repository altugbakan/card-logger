package keymaps

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type Add struct {
	Back key.Binding
	Add  key.Binding
	Undo key.Binding
	help help.Model
}

func NewAddKeyMap(bindings ...key.Binding) Add {
	return Add{
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "go back"),
		),
		Add: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "add card"),
		),
		Undo: key.NewBinding(
			key.WithKeys("ctrl+z"),
			key.WithHelp("ctrl+z", "undo"),
		),
		help: help.New(),
	}
}

func (k Add) Help() string {
	return k.help.View(k)
}

func (k Add) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Back,
		k.Add,
		k.Undo,
	}
}

func (k Add) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back, k.Add, k.Undo},
	}
}
