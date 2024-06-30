package keymaps

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type SetList struct {
	Back     key.Binding
	Navigate key.Binding
	Select   key.Binding
	Missing  key.Binding
	Search   key.Binding
	help     help.Model
}

func NewSetListKeyMap() SetList {
	return SetList{
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "go back"),
		),
		Navigate: key.NewBinding(
			key.WithKeys("up", "down"),
			key.WithHelp("↑/↓", "navigate"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Missing: key.NewBinding(
			key.WithKeys("m"),
			key.WithHelp("m", "show missing"),
		),
		Search: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),
		help: help.New(),
	}
}

func (k SetList) Help() string {
	return k.help.View(k)
}

func (k SetList) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Back,
		k.Navigate,
		k.Select,
		k.Missing,
	}
}

func (k SetList) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back, k.Navigate, k.Select, k.Missing},
	}
}
