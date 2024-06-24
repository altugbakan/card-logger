package keymaps

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type Set struct {
	Back     key.Binding
	Navigate key.Binding
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Add      key.Binding
	Remove   key.Binding
	Search   key.Binding
	help     help.Model
}

func NewSetKeyMap(bindings ...key.Binding) Set {
	return Set{
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "go back"),
		),
		Navigate: key.NewBinding(
			key.WithKeys("up", "down", "left", "right"),
			key.WithHelp("↑/↓/←/→", "navigate"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
		),
		Add: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add"),
		),
		Remove: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "remove"),
		),
		Search: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),
		help: help.New(),
	}
}

func (k Set) Help() string {
	return k.help.View(k)
}

func (k Set) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Back,
		k.Navigate,
		k.Add,
		k.Remove,
		k.Search,
	}
}

func (k Set) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back, k.Add, k.Navigate, k.Remove, k.Search},
	}
}
