package screens

import (
	"github.com/altugbakan/card-logger/components"
	"github.com/altugbakan/card-logger/db"
	"github.com/altugbakan/card-logger/keymaps"
	"github.com/altugbakan/card-logger/utils"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	listHeightMargin = 5
	listWidthMargin  = 2
)

type List struct {
	keyMap keymaps.List
	list   list.Model
	sets   []db.Set
}

func NewListScreen() (List, error) {
	keyMap := keymaps.NewListKeyMap()

	sets, err := db.GetAllSets()
	if err != nil {
		return List{}, err
	}

	userCardCounts, err := db.GetUserCardCountsBySet()
	if err != nil {
		return List{}, err
	}

	items := []list.Item{}
	for _, set := range sets {
		item := components.SetItem{
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

	width, height := utils.GetWindowSize()
	width -= listWidthMargin * 2
	height -= listHeightMargin*2 - utils.TotalHelpWidth
	utils.LogInfo("initializing list with size %d x %d", width, height)

	list := list.New(items, components.SetItemDelegate{MaxNameLength: maxNameLength},
		width, height)
	list.SetShowHelp(false)
	list.SetShowTitle(false)
	list.FilterInput.Cursor.Style = utils.CursorStyle
	list.FilterInput.PromptStyle = utils.ActionStyle
	list.DisableQuitKeybindings()

	return List{
		keyMap: keyMap,
		list:   list,
		sets:   sets,
	}, nil
}

func (s List) Update(msg tea.Msg) (Screen, tea.Cmd) {
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
			set, ok := s.list.SelectedItem().(components.SetItem)
			if !ok {
				utils.LogError("error casting selected item to SetItem")
				return s, nil
			}
			setScreen, err := NewSetScreen(set.Abbr)
			if err != nil {
				utils.LogError("error creating set screen: %v", err)
				return s, nil
			}

			return setScreen, nil
		}
	case tea.WindowSizeMsg:
		s.list.SetSize(msg.Width-listWidthMargin*2, msg.Height-listHeightMargin*2)
		return s, nil
	}

	var cmd tea.Cmd
	s.list, cmd = s.list.Update(msg)
	return s, cmd
}

func (s List) View() string {
	return s.list.View()
}

func (s List) Help() string {
	return s.keyMap.Help()
}
