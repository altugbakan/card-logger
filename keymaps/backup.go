package keymaps

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type Backup struct {
	Back key.Binding
	Save key.Binding
	Load key.Binding
	help help.Model
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
	}
}

func (k Backup) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back, k.Save, k.Load},
	}
}
