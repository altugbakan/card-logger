package keymaps

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type List struct {
	Back   key.Binding
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Search key.Binding
	help   help.Model
}

func NewListKeyMap() List {
	return List{
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "go back"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Search: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),
		help: help.New(),
	}
}

func (k List) Help() string {
	return k.help.View(k)
}

func (k List) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Back,
		k.Up,
		k.Down,
		k.Select,
		k.Search,
	}
}

func (k List) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back, k.Up, k.Down, k.Select, k.Search},
	}
}
