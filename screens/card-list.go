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

type CardList struct {
	keyMap       keymaps.Set
	list         list.Model
	name         string
	itemDelegate components.CardListItemDelegate
}

func NewCardListScreen(abbr string) (CardList, error) {
	keyMap := keymaps.NewSetKeyMap()

	set, err := db.GetSet(abbr)
	if err != nil {
		utils.LogError("Error getting set from db: %v", err)
		return CardList{}, err
	}

	userCards, err := db.GetAllSetCardsOfUser(abbr)
	if err != nil {
		utils.LogError("Error getting all set cards of user from db: %v", err)
		return CardList{}, err
	}

	items := []list.Item{}
	maxNameLength := 0
	for _, card := range userCards {
		item := components.CardListItem{
			CardID:   card.CardID,
			Number:   card.Number,
			Name:     card.Name,
			Patterns: card.Patterns,
		}
		maxNameLength = max(maxNameLength, len(card.Name))
		items = append(items, item)
	}

	itemDelegate := components.CardListItemDelegate{MaxNameLength: maxNameLength, SelectedIndex: 0}
	list := utils.NewList(items, itemDelegate, "card")

	return CardList{
		keyMap:       keyMap,
		list:         list,
		name:         set.Name,
		itemDelegate: itemDelegate,
	}, nil
}

func (s CardList) Update(msg tea.Msg) (Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s.list.SettingFilter() {
			break
		}
		switch {
		case key.Matches(msg, s.keyMap.Back):
			listScreen, err := NewSetListScreen()
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
			selectedItem, ok := s.list.SelectedItem().(components.CardListItem)
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
		width, height := utils.GetListSize(len(s.list.Items()), msg.Width, msg.Height)
		utils.LogInfo("resizing card list to %d x %d", width, height)
		s.list.SetSize(width, height)
		return s, nil
	}

	var cmd tea.Cmd
	s.list, cmd = s.list.Update(msg)
	return s, cmd
}

func (s CardList) View() string {
	title := utils.TitleStyle.MarginBottom(1).Render(s.name)
	return lipgloss.JoinVertical(lipgloss.Center, title, s.list.View())
}

func (s CardList) Help() string {
	return s.keyMap.Help()
}

func (s *CardList) handleAdd() {
	selectedItem, ok := s.list.SelectedItem().(components.CardListItem)
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

func (s *CardList) handleRemove() {
	selectedItem, ok := s.list.SelectedItem().(components.CardListItem)
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
