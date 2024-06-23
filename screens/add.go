package screens

import (
	"errors"
	"strconv"
	"strings"

	"github.com/altugbakan/card-logger/db"
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

type AddScreen struct {
	keyMap  utils.KeyMap
	input   textinput.Model
	set     string
	msg     string
	history []addCardArgs
}

func NewAddScreen() AddScreen {
	keyBindings := utils.NewKeyMap(
		key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "go back"),
		),
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "add card"),
		),
		key.NewBinding(
			key.WithKeys("ctrl+z"),
			key.WithHelp("ctrl+z", "undo"),
		),
	)

	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = inputCharLimit
	ti.Width = inputWidth
	ti.Placeholder = placeholder
	ti.PromptStyle = utils.ActionStyle

	msg := utils.DimTextStyle.Render(format)
	return AddScreen{
		keyMap: keyBindings,
		input:  ti,
		msg:    msg,
	}
}

func (s AddScreen) Update(msg tea.Msg) (Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if s.input.Value() == "" {
				return s, nil
			}
			s.msg = s.handleEnterKeyPress().Render()
		case "esc":
			return s.handleEscKeyPress()
		case "ctrl+z":
			s.msg = s.handleUndoKeyPress().Render()
		}
	}

	var cmd tea.Cmd
	s.input, cmd = s.input.Update(msg)

	return s, cmd
}

func (s AddScreen) View() string {
	title := utils.HeaderStyle.MarginBottom(1).Render("Add Card")
	input := lipgloss.JoinHorizontal(lipgloss.Center, utils.ActionStyle.Width(3).Render(s.set), s.input.View())
	titleAndInput := lipgloss.JoinVertical(lipgloss.Center, title, input)
	msg := "\n" + s.msg
	return lipgloss.JoinVertical(lipgloss.Center, titleAndInput, msg)
}

func (s AddScreen) Help() string {
	return s.keyMap.Help()
}

func (s *AddScreen) handleEnterKeyPress() utils.Message {
	submitResult, err := s.submit(s.input.Value())
	if err != nil {
		return utils.NewErrorMessage(err.Error())
	}

	s.input.SetValue("")
	switch result := submitResult.(type) {
	case addCardResult:
		return utils.NewInfoMessage("added " + result.name + " - " +
			utils.GetPatternText(result.rarity, result.patternAmounts))
	case changeSetResult:
		return s.changeSet(result.setName)
	}

	return utils.NewErrorMessage("invalid input")
}

func (s *AddScreen) handleEscKeyPress() (Screen, tea.Cmd) {
	if s.set != "" {
		s.resetSet()
		return s, nil
	}
	return NewTitleModel(), nil
}

func (s *AddScreen) handleUndoKeyPress() utils.Message {
	result, err := s.undoLastAddition()
	if err != nil {
		return utils.NewErrorMessage(err.Error())
	} else {
		return utils.NewInfoMessage("removed " + result.name + " - " +
			utils.GetPatternText(result.rarity, result.patternAmounts))
	}
}

func (s *AddScreen) undoLastAddition() (addCardResult, error) {
	if len(s.history) == 0 {
		return addCardResult{}, errors.New("no cards to undo")
	}

	lastAddition := s.history[len(s.history)-1]
	result, err := removeCard(lastAddition.set, lastAddition.number, lastAddition.pattern)
	if err != nil {
		return addCardResult{}, err
	}

	s.history = s.history[:len(s.history)-1]
	s.msg = utils.DimTextStyle.Render("removed " + result.name + " - " +
		utils.GetPatternText(result.rarity, result.patternAmounts))
	return result, nil
}

func (s *AddScreen) submit(input string) (submitResult, error) {
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

func (s *AddScreen) changeSet(set string) utils.Message {
	s.set = set
	s.input.Placeholder = placeholderWithoutSet
	s.input.SetValue("")
	return utils.NewInfoMessage(formatWithoutSet)
}

func (s *AddScreen) resetSet() utils.Message {
	s.set = ""
	s.input.Placeholder = placeholder
	return utils.NewInfoMessage(format)
}

func (s *AddScreen) handleOneArgument(args []string) (submitResult, error) {
	num, err := strconv.Atoi(args[0])
	if err == nil {
		return s.addCardDefault(s.set, num)
	}

	setName := strings.ToUpper(args[0])
	if checkSetExists(setName) {
		return changeSetResult{setName: setName}, nil
	} else {
		return changeSetResult{}, errors.New("set does not exist")
	}
}

func (s *AddScreen) handleTwoArguments(args []string) (submitResult, error) {
	if s.set == "" {
		_, err := strconv.Atoi(args[0])
		if err == nil {
			return addCardResult{}, errors.New("set not specified")
		}

		cardNumber, err := strconv.Atoi(args[1])
		if err != nil {
			return addCardResult{}, err
		}
		return s.addCardDefault(args[0], cardNumber)
	} else {
		cardNumber, err := strconv.Atoi(args[0])
		if err != nil {
			return addCardResult{}, err
		}
		return s.addCard(s.set, cardNumber, args[1])
	}
}

func (s *AddScreen) addCardDefault(set string, number int) (addCardResult, error) {
	return s.addCard(set, number, "")
}

func (s *AddScreen) addCard(set string, number int, pattern string) (addCardResult, error) {
	set = strings.ToUpper(set)
	pattern = utils.CorrectPattern(pattern)

	if !checkSetExists(set) {
		return addCardResult{}, errors.New("set does not exist")
	}

	// check if card exists
	card, err := db.GetCard(set, number)
	if err != nil {
		return addCardResult{}, err
	}

	// use first pattern if not provided
	if pattern == "" {
		pattern = utils.GetPatternsForRarity(card.Rarity)[0]
	} else if !utils.IsPatternValidForRarity(pattern, card.Rarity) {
		return addCardResult{}, errors.New("pattern " + pattern + " is not valid for rarity " + card.Rarity)
	}

	// add card to user's collection
	err = db.AddUserCard(card.ID, pattern)
	if err != nil {
		return addCardResult{}, err
	}

	// get all user's pattern amounts
	patternAmounts, err := db.GetAllUserPatternAmounts(card.ID)
	if err != nil {
		return addCardResult{}, err
	}

	// add to history
	s.history = append(s.history, addCardArgs{set: set, number: number, pattern: pattern})

	return addCardResult{
		name:           card.Name,
		rarity:         card.Rarity,
		patternAmounts: patternAmounts,
	}, nil
}

func removeCard(set string, number int, pattern string) (addCardResult, error) {
	set = strings.ToUpper(set)
	pattern = utils.CorrectPattern(pattern)

	if !checkSetExists(set) {
		return addCardResult{}, errors.New("set does not exist")
	}

	// check if card exists
	card, err := db.GetCard(set, number)
	if err != nil {
		return addCardResult{}, err
	}

	// use first pattern if not provided
	if pattern == "" {
		pattern = utils.GetPatternsForRarity(card.Rarity)[0]
	} else if !utils.IsPatternValidForRarity(pattern, card.Rarity) {
		return addCardResult{}, errors.New("pattern " + pattern + " is not valid for rarity " + card.Rarity)
	}

	// remove card from user's collection
	err = db.RemoveUserCard(card.ID, pattern)
	if err != nil {
		return addCardResult{}, err
	}

	// get all user's pattern amounts
	patternAmounts, err := db.GetAllUserPatternAmounts(card.ID)
	if err != nil {
		return addCardResult{}, err
	}

	return addCardResult{
		name:           card.Name,
		rarity:         card.Rarity,
		patternAmounts: patternAmounts,
	}, nil
}

func checkSetExists(set string) bool {
	_, err := db.GetSet(set)
	return err == nil
}
