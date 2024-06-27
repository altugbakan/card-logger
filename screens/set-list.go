package screens

import (
	"github.com/altugbakan/card-logger/db"
	"github.com/altugbakan/card-logger/items"
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
		utils.LogError("could not get all sets: %v", err)
		return SetList{}, err
	}

	setItems, maxNameLength, err := getSetItems()
	if err != nil {
		utils.LogError("could not get set items: %v", err)
		return SetList{}, err
	}

	list := utils.NewList(setItems, items.SetDelegate{MaxNameLength: maxNameLength}, "set")

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
			set, ok := s.list.SelectedItem().(items.Set)
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
		case key.Matches(msg, s.keyMap.Missing):
			set, ok := s.list.SelectedItem().(items.Set)
			if !ok {
				utils.LogError("error casting selected item to SetItem")
				return s, nil
			}
			missingScreen, err := NewMissingListScreen(set.Abbr)
			if err != nil {
				utils.LogError("error creating missing screen: %v", err)
				return s, nil
			}
			return missingScreen, nil
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

func getSetItems() ([]list.Item, int, error) {
	sets, err := db.GetAllSets()
	if err != nil {
		return nil, 0, err
	}

	userCardCounts, err := db.GetUserCardCountsBySet()
	if err != nil {
		return nil, 0, err
	}

	setItems := make([]list.Item, len(sets))
	maxNameLength := 0
	for i, set := range sets {
		item := items.Set{
			Abbr:  set.Abbr,
			Name:  set.Name,
			Total: set.TotalCards,
			Owned: userCardCounts[set.Abbr],
		}
		maxNameLength = max(maxNameLength, len(item.FilterValue()))
		setItems[i] = item
	}

	return setItems, maxNameLength, nil
}
