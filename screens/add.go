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

type AddScreen struct {
	keyMap   utils.KeyMap
	input    textinput.Model
	set      string
	errorMsg string
}

func NewAddScreen() AddScreen {
	keyBindings := utils.NewKeyMap(
		key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "Go back"),
		),
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "Add card"),
		),
	)
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 20
	ti.Width = 20
	return AddScreen{
		keyMap: keyBindings,
		input:  ti,
	}
}

func (s AddScreen) Update(msg tea.KeyMsg) (Screen, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "enter":
		err := s.addCard(s.input.Value())
		if err != nil {
			s.errorMsg = err.Error()
		} else {
			s.errorMsg = ""
			s.input.SetValue("")
		}
		return s, nil
	case "esc":
		return NewTitleModel(), nil
	}
	s.input, cmd = s.input.Update(msg)

	return s, cmd
}

func (s AddScreen) View() string {
	title := utils.HeaderStyle.Render("Add Card")
	input := s.input.View()
	if s.errorMsg != "" {
		input = lipgloss.JoinVertical(lipgloss.Left, input,
			utils.ErrorStyle.Render(s.errorMsg))
	}

	return lipgloss.JoinVertical(lipgloss.Center, title, input)
}

func (s AddScreen) Help() string {
	return s.keyMap.Help()
}

func (s AddScreen) addCard(input string) error {
	args := strings.Split(input, " ")

	if len(args) == 0 || len(args) > 3 {
		return errors.New("invalid input")
	}

	if len(args) == 1 {
		//TODO: set set
	} else if len(args) == 2 {
		//TODO: handle
	} else {
		cardNumber, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}
		return addCardToDB(args[0], cardNumber, args[2])
	}

	return nil
}

func addCardToDB(set string, number int, pattern string) error {
	if !checkSetExists(set) {
		return errors.New("set does not exist")
	}

	// check if card exists
	card, err := db.GetCard(set, number)
	if err != nil {
		return err
	}

	// check if pattern is valid
	if !utils.IsPatternValidForRarity(pattern, card.Rarity) {
		return errors.New("invalid pattern for rarity")
	}

	// add card to user's collection
	return db.AddUserCard(card.ID, pattern)
}

func checkSetExists(set string) bool {
	_, err := db.GetSet(set)
	return err == nil
}
