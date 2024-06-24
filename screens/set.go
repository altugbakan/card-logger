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

const (
	setHeightMargin = 7
	setWidthMargin  = 2
)

type Set struct {
	keyMap       keymaps.Set
	list         list.Model
	name         string
	itemDelegate components.CardItemDelegate
}

func NewSetScreen(abbr string) (Set, error) {
	keyMap := keymaps.NewSetKeyMap()

	set, err := db.GetSet(abbr)
	if err != nil {
		utils.LogError("Error getting set from db: %v", err)
		return Set{}, err
	}

	userCards, err := db.GetAllSetCardsOfUser(abbr)
	if err != nil {
		utils.LogError("Error getting all set cards of user from db: %v", err)
		return Set{}, err
	}

	items := []list.Item{}
	maxNameLength := 0
	for _, card := range userCards {
		item := components.CardItem{
			CardID:   card.CardID,
			Number:   card.Number,
			Name:     card.Name,
			Patterns: card.Patterns,
		}
		maxNameLength = max(maxNameLength, len(card.Name))
		items = append(items, item)
	}

	width, height := utils.GetWindowSize()
	width -= setWidthMargin * 2
	height -= setHeightMargin*2 - utils.TotalHelpWidth
	utils.LogInfo("initializing set with size %d x %d", width, height)

	itemDelegate := components.CardItemDelegate{MaxNameLength: maxNameLength, SelectedIndex: 0}
	list := list.New(items, itemDelegate, width, height)
	list.SetShowHelp(false)
	list.SetShowTitle(false)
	list.FilterInput.Cursor.Style = utils.CursorStyle
	list.FilterInput.PromptStyle = utils.ActionStyle
	list.DisableQuitKeybindings()

	return Set{
		keyMap:       keyMap,
		list:         list,
		name:         set.Name,
		itemDelegate: itemDelegate,
	}, nil
}

func (s Set) Update(msg tea.Msg) (Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s.list.SettingFilter() {
			break
		}
		switch {
		case key.Matches(msg, s.keyMap.Back):
			listScreen, err := NewListScreen()
			if err != nil {
				utils.LogError("Error creating list screen: %v", err)
				return s, nil
			}
			return listScreen, nil
		case key.Matches(msg, s.keyMap.Left):
			s.itemDelegate.SelectedIndex = max(0, s.itemDelegate.SelectedIndex-1)
			s.list.SetDelegate(s.itemDelegate)
			return s, nil
		case key.Matches(msg, s.keyMap.Right):
			selectedItem, ok := s.list.SelectedItem().(components.CardItem)
			if !ok {
				utils.LogError("error casting selected item to CardItem")
				return s, nil
			}
			s.itemDelegate.SelectedIndex = min(s.itemDelegate.SelectedIndex+1,
				len(selectedItem.Patterns)-1)
			s.list.SetDelegate(s.itemDelegate)
			return s, nil
		case key.Matches(msg, s.keyMap.Add):
			s.handleAdd()
			return s, nil
		case key.Matches(msg, s.keyMap.Remove):
			s.handleRemove()
			return s, nil
		}
	case tea.WindowSizeMsg:
		s.list.SetSize(msg.Width-setWidthMargin*2, msg.Height-setHeightMargin*2)
		return s, nil
	}

	var cmd tea.Cmd
	s.list, cmd = s.list.Update(msg)
	return s, cmd
}

func (s Set) View() string {
	title := utils.TitleStyle.MarginBottom(1).Render(s.name)
	return lipgloss.JoinVertical(lipgloss.Center, title, s.list.View())
}

func (s Set) Help() string {
	return s.keyMap.Help()
}

func (s *Set) handleAdd() {
	selectedItem, ok := s.list.SelectedItem().(components.CardItem)
	if !ok {
		utils.LogError("error casting selected item to CardItem")
		return
	}

	selectedPattern := selectedItem.Patterns[s.itemDelegate.SelectedIndex].Name
	err := db.AddUserCard(selectedItem.CardID, selectedPattern)
	if err != nil {
		utils.LogError("error adding user card: %v", err)
		return
	} else {
		utils.LogInfo("added user card %s with pattern %s", selectedItem.Name, selectedPattern)
	}
	selectedItem.Patterns[s.itemDelegate.SelectedIndex].Quantity++
	s.list.SetItem(s.list.Index(), selectedItem)
}

func (s *Set) handleRemove() {
	selectedItem, ok := s.list.SelectedItem().(components.CardItem)
	if !ok {
		utils.LogError("error casting selected item to CardItem")
		return
	}

	selectedPattern := selectedItem.Patterns[s.itemDelegate.SelectedIndex].Name
	err := db.RemoveUserCard(selectedItem.CardID, selectedPattern)
	if err != nil {
		utils.LogError("error removing user card: %v", err)
		return
	} else {
		utils.LogInfo("removed user card %s with pattern %s", selectedItem.Name, selectedPattern)
	}
	selectedItem.Patterns[s.itemDelegate.SelectedIndex].Quantity--
	s.list.SetItem(s.list.Index(), selectedItem)
}
