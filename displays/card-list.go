package displays

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

type CardList struct {
	keyMap       keymaps.Set
	list         list.Model
	name         string
	itemDelegate items.CardDelegate
}

func NewCardListScreen(abbr string) (CardList, error) {
	keyMap := keymaps.NewSetKeyMap()

	set, err := db.GetSet(abbr)
	if err != nil {
		utils.LogError("error getting set from db: %v", err)
		return CardList{}, err
	}

	cardItems, maxNameLength, maxPatternLength, err := getCardItems(abbr)
	if err != nil {
		utils.LogError("error getting card items: %v", err)
		return CardList{}, err
	}

	itemDelegate := items.CardDelegate{
		MaxNameLength:    maxNameLength,
		MaxPatternLength: maxPatternLength,
		SelectedIndex:    0,
	}
	list := utils.NewList(cardItems, itemDelegate, "card")

	return CardList{
		keyMap:       keyMap,
		list:         list,
		name:         set.Name,
		itemDelegate: itemDelegate,
	}, nil
}

func (s CardList) Update(msg tea.Msg) (Displayer, tea.Cmd) {
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
			listScreen, err := NewSetListScreen()
			if err != nil {
				utils.LogError("Error creating list screen: %v", err)
				return s, nil
			}
			return listScreen, nil
		case key.Matches(msg, s.keyMap.Left):
			return s.handleLeft()
		case key.Matches(msg, s.keyMap.Right):
			return s.handleRight()
		case key.Matches(msg, s.keyMap.Add):
			s.handleAdd()
			return s, nil
		case key.Matches(msg, s.keyMap.Remove):
			s.handleRemove()
			return s, nil
		}
	case tea.WindowSizeMsg:
		utils.SetListSize(&s.list, msg.Width, msg.Height)
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
	selectedItem, ok := s.list.SelectedItem().(items.Card)
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
	selectedItem, ok := s.list.SelectedItem().(items.Card)
	if !ok {
		utils.LogError("error casting selected item to CardItem")
		return
	}

	selectedPattern := selectedItem.Patterns[s.itemDelegate.SelectedIndex].Name

	if selectedItem.Patterns[s.itemDelegate.SelectedIndex].Quantity < 1 {
		return
	}

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

func (s *CardList) handleLeft() (Displayer, tea.Cmd) {
	if s.itemDelegate.SelectedIndex == 0 {
		prevSelectedItem, ok := s.list.SelectedItem().(items.Card)
		if !ok {
			utils.LogError("error casting selected item to CardItem")
			return s, nil
		}

		var cmd tea.Cmd
		s.list, cmd = s.list.Update(tea.KeyMsg{Type: tea.KeyLeft})

		selectedItem, ok := s.list.SelectedItem().(items.Card)
		if !ok {
			utils.LogError("error casting selected item to CardItem")
			return s, nil
		}

		if selectedItem.CardID == prevSelectedItem.CardID {
			return s, cmd
		}

		s.itemDelegate.SelectedIndex = len(selectedItem.Patterns) - 1
		s.list.SetDelegate(s.itemDelegate)

		return s, cmd
	}

	s.itemDelegate.SelectedIndex--
	s.list.SetDelegate(s.itemDelegate)
	return s, nil
}

func (s *CardList) handleRight() (Displayer, tea.Cmd) {
	selectedItem, ok := s.list.SelectedItem().(items.Card)
	if !ok {
		utils.LogError("error casting selected item to CardItem")
		return s, nil
	}

	if s.itemDelegate.SelectedIndex == len(selectedItem.Patterns)-1 {
		var cmd tea.Cmd
		s.list, cmd = s.list.Update(tea.KeyMsg{Type: tea.KeyRight})

		nextSelectedItem, ok := s.list.SelectedItem().(items.Card)
		if !ok {
			utils.LogError("error casting selected item to CardItem")
			return s, nil
		}

		if nextSelectedItem.CardID == selectedItem.CardID {
			return s, cmd
		}

		s.itemDelegate.SelectedIndex = 0
		s.list.SetDelegate(s.itemDelegate)

		return s, cmd
	}

	s.itemDelegate.SelectedIndex++
	s.list.SetDelegate(s.itemDelegate)
	return s, nil
}

func getCardItems(abbr string) ([]list.Item, int, int, error) {
	userCards, err := db.GetAllSetCardsOfUser(abbr)
	if err != nil {
		return nil, 0, 0, err
	}

	cardItems := make([]list.Item, len(userCards))
	maxNameLength := 0
	maxPatternLength := 0
	for i, card := range userCards {
		item := items.Card{
			CardID:   card.CardID,
			Number:   card.Number,
			Name:     card.Name,
			Patterns: card.Patterns,
		}
		maxNameLength = max(maxNameLength, len(card.Name))
		maxPatternLength = max(maxPatternLength, len(utils.GetPatternItemText(card.Patterns)))
		cardItems[i] = item
	}

	return cardItems, maxNameLength, maxPatternLength, nil
}
