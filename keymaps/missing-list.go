package keymaps

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type MissingList struct {
	Back     key.Binding
	Navigate key.Binding
	Toggle   key.Binding
	Search   key.Binding
	help     help.Model
}

func NewMissingKeyMap() MissingList {
	return MissingList{
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "go back"),
		),
		Navigate: key.NewBinding(
			key.WithKeys("up", "down"),
			key.WithHelp("↑/↓", "navigate"),
		),
		Toggle: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "toggle"),
		),
		Search: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),
		help: help.New(),
	}
}

func (k MissingList) Help() string {
	return k.help.View(k)
}

func (k MissingList) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Back,
		k.Navigate,
		k.Toggle,
		k.Search,
	}
}

func (k MissingList) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back, k.Navigate, k.Toggle, k.Search},
	}
}
