package screens

import (
	"github.com/altugbakan/card-logger/components"
	"github.com/altugbakan/card-logger/db"
	"github.com/altugbakan/card-logger/keymaps"
	"github.com/altugbakan/card-logger/utils"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SetList struct {
	keyMap keymaps.SetList
	list   list.Model
	sets   []db.Set
}

func NewSetListScreen() (SetList, error) {
	keyMap := keymaps.NewListKeyMap()

	sets, err := db.GetAllSets()
	if err != nil {
		return SetList{}, err
	}

	userCardCounts, err := db.GetUserCardCountsBySet()
	if err != nil {
		return SetList{}, err
	}

	items := []list.Item{}
	for _, set := range sets {
		item := components.SetListItem{
			Abbr:  set.Abbr,
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

	list := utils.NewList(items, components.SetListItemDelegate{MaxNameLength: maxNameLength}, "set")

	return SetList{
		keyMap: keyMap,
		list:   list,
		sets:   sets,
	}, nil
}

func (s SetList) Update(msg tea.Msg) (Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s.list.SettingFilter() {
			break
		}
		switch {
		case key.Matches(msg, s.keyMap.Back):
			if s.list.IsFiltered() {
				s.list.ResetFilter()
				return s, nil
			}
			return NewTitleScreen(), nil
		case key.Matches(msg, s.keyMap.Select):
			set, ok := s.list.SelectedItem().(components.SetListItem)
			if !ok {
				utils.LogError("error casting selected item to SetItem")
				return s, nil
			}
			setScreen, err := NewCardListScreen(set.Abbr)
			if err != nil {
				utils.LogError("error creating set screen: %v", err)
				return s, nil
			}

			return setScreen, nil
		}
	case tea.WindowSizeMsg:
		utils.SetListSize(&s.list, msg.Width, msg.Height)
		return s, nil
	}

	var cmd tea.Cmd
	s.list, cmd = s.list.Update(msg)
	return s, cmd
}

func (s SetList) View() string {
	title := utils.TitleStyle.MarginBottom(1).Render("Expansion Sets")
	list := s.list.View()
	return lipgloss.JoinVertical(lipgloss.Center, title, list)
}

func (s SetList) Help() string {
	return s.keyMap.Help()
}
