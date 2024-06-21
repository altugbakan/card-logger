package utils

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Bindings []key.Binding
	help     help.Model
}

func NewKeyMap(bindings ...key.Binding) KeyMap {
	return KeyMap{Bindings: bindings, help: help.New()}
}

func (k KeyMap) Help() string {
	return k.help.View(k)
}

func (k KeyMap) ShortHelp() []key.Binding {
	return k.Bindings
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.Bindings,
	}
}
