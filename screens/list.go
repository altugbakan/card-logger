package screens

import (
	"github.com/altugbakan/card-logger/components"
	"github.com/altugbakan/card-logger/db"
	"github.com/altugbakan/card-logger/utils"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ListScreen struct {
	keyMap utils.KeyMap
	list   list.Model
}

func NewListScreen() (ListScreen, error) {
	keyMap := utils.NewKeyMap(
		key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "go back"),
		),
		key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),
	)

	sets, err := db.GetAllSets()
	if err != nil {
		return ListScreen{}, err
	}

	userCardCounts, err := db.GetUserCardCountsBySet()
	if err != nil {
		return ListScreen{}, err
	}

	items := []list.Item{}
	for _, set := range sets {
		item := components.SetItem{
			Name:  set.Name,
			Total: set.TotalCards,
			Owned: userCardCounts[set.Abbr],
		}
		items = append(items, item)
	}

	maxNameLength := 0
	for _, item := range items {
		if len(item.FilterValue()) > maxNameLength {
			maxNameLength = len(item.FilterValue())
		}
	}

	list := list.New(items, components.SetItemDelegate{MaxNameLength: maxNameLength}, 20, 20)
	list.SetShowHelp(false)
	list.SetShowTitle(false)
	list.FilterInput.Cursor.Style = utils.CursorStyle
	list.FilterInput.PromptStyle = utils.ActionStyle
	list.KeyMap.Quit.SetEnabled(false)
	list.KeyMap.ForceQuit.SetEnabled(false)

	return ListScreen{
		keyMap: keyMap,
		list:   list,
	}, nil
}

func (s ListScreen) Update(msg tea.Msg) (Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if s.list.SettingFilter() {
				s.list.ResetFilter()
				return s, nil
			}
			return NewTitleModel(), nil
		case "enter":
			// TODO: open set screen
		}
	}

	var cmd tea.Cmd
	s.list, cmd = s.list.Update(msg)
	return s, cmd
}

func (s ListScreen) View() string {
	return s.list.View()
}

func (s ListScreen) Help() string {
	return s.keyMap.Help()
}
