package screens

import (
	"strings"

	"github.com/altugbakan/card-logger/db"
	"github.com/altugbakan/card-logger/items"
	"github.com/altugbakan/card-logger/keymaps"
	"github.com/altugbakan/card-logger/utils"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type viewMode int

type noPatternValues struct {
	items    []list.Item
	delegate items.NoPatternDelegate
}

type incompleteValues struct {
	items    []list.Item
	delegate items.IncompleteDelegate
}

const (
	viewNoPatterns viewMode = iota
	viewIncompletePatterns
)

type MissingList struct {
	keyMap           keymaps.MissingList
	mode             viewMode
	noPatternValues  noPatternValues
	incompleteValues incompleteValues
	list             list.Model
}

func NewMissingListScreen(abbr string) (MissingList, error) {
	keyMap := keymaps.NewMissingKeyMap()

	noPatternValues, err := getNoPatternItems(abbr)
	if err != nil {
		utils.LogError("could not get no pattern cards: %v", err)
		return MissingList{}, err
	}

	incompleteValues, err := getIncompleteItems(abbr)
	if err != nil {
		utils.LogError("could not get incomplete cards: %v", err)
		return MissingList{}, err
	}

	list := utils.NewList(noPatternValues.items, noPatternValues.delegate, "no pattern")

	return MissingList{
		keyMap:           keyMap,
		mode:             viewNoPatterns,
		noPatternValues:  noPatternValues,
		incompleteValues: incompleteValues,
		list:             list,
	}, nil
}

func (s MissingList) Update(msg tea.Msg) (Screen, tea.Cmd) {
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
			setList, err := NewSetListScreen()
			if err != nil {
				utils.LogError("error creating set list screen: %v", err)
				return s, nil
			}
			return setList, nil
		case key.Matches(msg, s.keyMap.Toggle):
			if s.mode == viewNoPatterns {
				s.list = utils.NewList(s.incompleteValues.items, s.incompleteValues.delegate, "incomplete")
				s.mode = viewIncompletePatterns
			} else {
				s.list = utils.NewList(s.noPatternValues.items, s.noPatternValues.delegate, "no pattern")
				s.mode = viewNoPatterns
			}
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

func (s MissingList) View() string {
	var title string
	if s.mode == viewNoPatterns {
		title = "Cards with No Patterns"
	} else {
		title = "Cards with Incomplete Patterns"
	}

	return lipgloss.JoinVertical(lipgloss.Center,
		utils.TitleStyle.MarginBottom(1).Render(title),
		s.list.View())
}

func (s MissingList) Help() string {
	return s.keyMap.Help()
}

func getNoPatternItems(set string) (noPatternValues, error) {
	missingCards, err := db.GetUserCardsWithNoPatternsForSet(set)
	if err != nil {
		return noPatternValues{}, err
	}

	maxNameLength := 0
	maxRarityLength := 0
	cardItems := make([]list.Item, len(missingCards))
	for i, card := range missingCards {
		cardItems[i] = items.NoPattern{
			Number: card.Number,
			Name:   card.Name,
			Rarity: card.Rarity,
		}
		maxNameLength = max(maxNameLength, len(card.Name))
		maxRarityLength = max(maxRarityLength, len(card.Rarity))
	}

	return noPatternValues{
		items: cardItems,
		delegate: items.NoPatternDelegate{MaxNameLength: maxNameLength,
			MaxRarityLength: maxRarityLength},
	}, nil
}

func getIncompleteItems(set string) (incompleteValues, error) {
	missingCards, err := db.GetUserCardsWithIncompletePatternsForSet(set)
	if err != nil {
		return incompleteValues{}, err
	}

	maxNameLength := 0
	maxPatternLength := 0
	maxRarityLength := 0
	cardItems := make([]list.Item, len(missingCards))
	for i, card := range missingCards {
		cardItems[i] = items.Incomplete{
			Number:   card.Number,
			Name:     card.Name,
			Rarity:   card.Rarity,
			Patterns: card.Patterns,
		}
		maxNameLength = max(maxNameLength, len(card.Name))
		maxPatternLength = max(maxPatternLength, len(strings.Join(card.Patterns, ", ")))
		maxRarityLength = max(maxRarityLength, len(card.Rarity))
	}

	return incompleteValues{
		items: cardItems,
		delegate: items.IncompleteDelegate{MaxNameLength: maxNameLength,
			MaxPatternLength: maxPatternLength,
			MaxRarityLength:  maxRarityLength},
	}, nil
}
