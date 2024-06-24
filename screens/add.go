package screens

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/altugbakan/card-logger/db"
	"github.com/altugbakan/card-logger/keymaps"
	"github.com/altugbakan/card-logger/utils"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	placeholder           = "e.g. TEF 1 RH"
	placeholderWithoutSet = "e.g. 1 RH"
	format                = "format: set number pattern"
	formatWithoutSet      = "format: number pattern"
	inputWidth            = 15
	inputCharLimit        = 15
)

type submitResult interface{}

type addCardResult struct {
	name           string
	set            string
	rarity         string
	patternAmounts map[string]int
}

type addCardArgs struct {
	set     string
	number  int
	pattern string
}

type changeSetResult struct {
	setName string
}

type emptySubmitResult struct{}

type Add struct {
	keyMap  keymaps.Add
	input   textinput.Model
	set     string
	msg     string
	history []addCardArgs
}

func NewAddScreen() Add {
	keyMap := keymaps.NewAddKeyMap()

	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = inputCharLimit
	ti.Width = inputWidth
	ti.Placeholder = placeholder
	ti.PromptStyle = utils.ActionStyle

	msg := utils.DimTextStyle.Render(format)
	return Add{
		keyMap: keyMap,
		input:  ti,
		msg:    msg,
	}
}

func (s Add) Update(msg tea.Msg) (Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keyMap.Add):
			if s.input.Value() == "" {
				return s, nil
			}
			s.msg = s.handleAddKeyPress().Render()
		case key.Matches(msg, s.keyMap.Back):
			return s.handleBackKeyPress()
		case key.Matches(msg, s.keyMap.Undo):
			s.msg = s.handleUndoKeyPress().Render()
		}
	}

	var cmd tea.Cmd
	s.input, cmd = s.input.Update(msg)

	return s, cmd
}

func (s Add) View() string {
	title := utils.TitleStyle.MarginBottom(1).Render("Add Card")
	input := lipgloss.JoinHorizontal(lipgloss.Center, utils.ActionStyle.Width(3).Render(s.set), s.input.View())
	titleAndInput := lipgloss.JoinVertical(lipgloss.Center, title, input)
	msg := utils.EmptyStyle.MarginTop(1).Render(s.msg)
	return lipgloss.JoinVertical(lipgloss.Center, titleAndInput, msg)
}

func (s Add) Help() string {
	return s.keyMap.Help()
}

func (s *Add) handleAddKeyPress() utils.Message {
	submitResult, err := s.submit(s.input.Value())
	if err != nil {
		return utils.NewErrorMessage(err.Error())
	}

	s.input.SetValue("")
	switch result := submitResult.(type) {
	case addCardResult:
		possiblePatterns, err := db.GetPatternsForRarity(result.set, result.rarity)
		if err != nil {
			utils.LogError("could not get patterns for rarity %s: %v", result.rarity, err)
			return utils.NewErrorMessage("could not get patterns for rarity")
		}
		message := fmt.Sprintf("added %s - %s", result.name,
			utils.GetPatternText(result.rarity, possiblePatterns, result.patternAmounts))
		utils.LogInfo(message)
		return utils.NewInfoMessage(message)
	case changeSetResult:
		return s.changeSet(result.setName)
	}

	return utils.NewErrorMessage("invalid input")
}

func (s *Add) handleBackKeyPress() (Screen, tea.Cmd) {
	if s.set != "" {
		s.resetSet()
		return s, nil
	}
	return NewTitleScreen(), nil
}

func (s *Add) handleUndoKeyPress() utils.Message {
	result, err := s.undoLastAddition()
	if err != nil {
		return utils.NewErrorMessage(err.Error())
	} else {
		possiblePatterns, err := db.GetPatternsForRarity(result.set, result.rarity)
		if err != nil {
			utils.LogError("could not get patterns for rarity %s: %v", result.rarity, err)
			return utils.NewErrorMessage("could not get patterns for rarity")
		}
		message := fmt.Sprintf("removed %s - %s", result.name,
			utils.GetPatternText(result.rarity, possiblePatterns, result.patternAmounts))
		utils.LogInfo(message)
		return utils.NewInfoMessage(message)
	}
}

func (s *Add) undoLastAddition() (addCardResult, error) {
	if len(s.history) == 0 {
		return addCardResult{}, errors.New("no cards to undo")
	}

	lastAddition := s.history[len(s.history)-1]
	result, err := removeCard(lastAddition.set, lastAddition.number, lastAddition.pattern)
	if err != nil {
		utils.LogError("could not remove card: %v", err)
		return addCardResult{}, err
	}

	s.history = s.history[:len(s.history)-1]
	return result, nil
}

func (s *Add) submit(input string) (submitResult, error) {
	args := strings.Fields(input)

	if len(args) == 0 {
		return emptySubmitResult{}, nil
	}

	if len(args) > 3 {
		return emptySubmitResult{}, errors.New("too many arguments")
	}

	if len(args) == 1 {
		return s.handleOneArgument(args)
	} else if len(args) == 2 {
		return s.handleTwoArguments(args)
	} else {
		cardNumber, err := strconv.Atoi(args[1])
		if err != nil {
			return addCardResult{}, err
		}
		return s.addCard(args[0], cardNumber, args[2])
	}
}

func (s *Add) changeSet(set string) utils.Message {
	s.set = set
	s.input.Placeholder = placeholderWithoutSet
	s.input.SetValue("")
	return utils.NewInfoMessage(formatWithoutSet)
}

func (s *Add) resetSet() utils.Message {
	s.set = ""
	s.input.Placeholder = placeholder
	return utils.NewInfoMessage(format)
}

func (s *Add) handleOneArgument(args []string) (submitResult, error) {
	num, err := strconv.Atoi(args[0])
	if err == nil {
		return s.addCardDefault(s.set, num)
	}

	setName := strings.ToUpper(args[0])
	if checkSetExists(setName) {
		return changeSetResult{setName: setName}, nil
	} else {
		utils.LogWarning("non-existing set %s specified while adding card", setName)
		return changeSetResult{}, fmt.Errorf("set %s does not exist", setName)
	}
}

func (s *Add) handleTwoArguments(args []string) (submitResult, error) {
	if s.set == "" {
		_, err := strconv.Atoi(args[0])
		if err == nil {
			return addCardResult{}, errors.New("set not specified")
		}

		cardNumber, err := strconv.Atoi(args[1])
		if err != nil {
			utils.LogWarning("could not convert card number to integer while adding card: %v", err)
			return addCardResult{}, fmt.Errorf("invalid card number %s", args[1])
		}
		return s.addCardDefault(args[0], cardNumber)
	} else {
		cardNumber, err := strconv.Atoi(args[0])
		if err != nil {
			utils.LogWarning("could not convert card number to integer while adding card: %v", err)
			return addCardResult{}, fmt.Errorf("invalid card number %s", args[0])
		}
		return s.addCard(s.set, cardNumber, args[1])
	}
}

func (s *Add) addCardDefault(set string, number int) (addCardResult, error) {
	return s.addCard(set, number, "")
}

func (s *Add) addCard(set string, number int, pattern string) (addCardResult, error) {
	card, pattern, err := validateInput(set, number, pattern, "adding")
	if err != nil {
		return addCardResult{}, err
	}

	err = db.AddUserCard(card.ID, pattern)
	if err != nil {
		utils.LogError("could not add card to database: %v", err)
		return addCardResult{}, fmt.Errorf("could not add card %s", card.Name)
	}

	patternAmounts, err := db.GetAllUserPatternAmounts(card.ID)
	if err != nil {
		utils.LogError("could not get pattern amounts from database while adding card: %v", err)
		return addCardResult{}, fmt.Errorf("could not get pattern amounts for card %s", card.Name)
	}

	s.history = append(s.history, addCardArgs{set: set, number: number, pattern: pattern})

	return addCardResult{
		name:           card.Name,
		set:            card.Set,
		rarity:         card.Rarity,
		patternAmounts: patternAmounts,
	}, nil
}

func removeCard(set string, number int, pattern string) (addCardResult, error) {
	card, pattern, err := validateInput(set, number, pattern, "removing")
	if err != nil {
		return addCardResult{}, err
	}

	err = db.RemoveUserCard(card.ID, pattern)
	if err != nil {
		utils.LogError("could not remove card from database: %v", err)
		return addCardResult{}, err
	}

	patternAmounts, err := db.GetAllUserPatternAmounts(card.ID)
	if err != nil {
		utils.LogError("could not get pattern amounts from database while removing card: %v", err)
		return addCardResult{}, fmt.Errorf("could not get pattern amounts for card %s", card.Name)
	}

	return addCardResult{
		name:           card.Name,
		set:            card.Set,
		rarity:         card.Rarity,
		patternAmounts: patternAmounts,
	}, nil
}

func checkSetExists(set string) bool {
	_, err := db.GetSet(set)
	return err == nil
}

func validateInput(set string, number int, pattern string, operation string) (db.Card, string, error) {
	set = strings.ToUpper(set)
	pattern = utils.CorrectPattern(pattern)

	if !checkSetExists(set) {
		utils.LogWarning("non-existing set %s specified while %s card", set, operation)
		return db.Card{}, "", fmt.Errorf("set %s does not exist", set)
	}

	card, err := db.GetCard(set, number)
	if err != nil {
		utils.LogWarning("could not get card from database while %s card: %v", operation, err)
		return db.Card{}, "", err
	}

	if pattern == "" {
		patterns, err := db.GetPatternsForRarity(card.Set, card.Rarity)
		if err != nil {
			utils.LogError("could not get patterns for rarity %s while %s card: %v", card.Rarity, operation, err)
			return db.Card{}, "", fmt.Errorf("could not get patterns for rarity %s", card.Rarity)
		}
		pattern = patterns[0]
	} else {
		valid, err := db.IsPatternValidForRarity(card.Set, card.Rarity, pattern)
		if err != nil {
			utils.LogError("could not check if pattern is valid for rarity %s while %s card: %v", card.Rarity, operation, err)
			return db.Card{}, "", fmt.Errorf("could not check if pattern is valid for rarity %s", card.Rarity)
		}
		if !valid {
			utils.LogWarning("invalid pattern %s specified for rarity %s while %s card", pattern, card.Rarity, operation)
			return db.Card{}, "", fmt.Errorf("pattern %s is not valid for rarity %s", pattern, card.Rarity)
		}
	}

	return card, pattern, nil
}
