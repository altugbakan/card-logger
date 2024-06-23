package keymaps

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type Backup struct {
	Back   key.Binding
	Save   key.Binding
	Load   key.Binding
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Search key.Binding
	help   help.Model
}

func NewBackupKeyMap() Backup {
	return Backup{
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "go back"),
		),
		Save: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "save backup"),
		),
		Load: key.NewBinding(
			key.WithKeys("l"),
			key.WithHelp("l", "load backup"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
			key.WithDisabled(),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
			key.WithDisabled(),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
			key.WithDisabled(),
		),
		Search: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
			key.WithDisabled(),
		),
		help: help.New(),
	}
}

func (k Backup) Help() string {
	return k.help.View(k)
}

func (k Backup) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Back,
		k.Save,
		k.Load,
		k.Up,
		k.Down,
		k.Select,
		k.Search,
	}
}

func (k Backup) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back, k.Save, k.Load, k.Up, k.Down, k.Select, k.Search},
	}
}
