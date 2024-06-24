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
			key.WithHelp("a", "add card"),
		),
		Remove: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "remove card"),
		),
	}
}

func (k Set) Help() string {
	return k.help.View(k)
}

func (k Set) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Back,
	}
}

func (k Set) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back, k.Add, k.Navigate, k.Remove},
	}
}
