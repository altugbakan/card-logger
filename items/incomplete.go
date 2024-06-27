package items

import (
	"io"
	"strconv"
	"strings"

	"github.com/altugbakan/card-logger/utils"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Incomplete struct {
	Number   int
	Name     string
	Rarity   string
	Patterns []string
}

func (i Incomplete) FilterValue() string {
	return i.Name
}

type IncompleteDelegate struct {
	MaxNameLength    int
	MaxRarityLength  int
	MaxPatternLength int
}

func (d IncompleteDelegate) Height() int                               { return utils.ItemHeight }
func (d IncompleteDelegate) Spacing() int                              { return utils.ItemSpacing }
func (d IncompleteDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d IncompleteDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(Incomplete)
	if !ok {
		return
	}

	var display string
	if index == m.Index() {
		number := utils.ActionStyle.Width(3).MarginRight(1).Render(strconv.Itoa(item.Number))
		itemName := utils.ActionStyle.Width(d.MaxNameLength).MarginRight(1).Render(item.Name)
		rarity := utils.ActionStyle.Width(d.MaxRarityLength).MarginRight(1).Render(item.Rarity)
		selectedPatternText := utils.ActionStyle.Width(d.MaxPatternLength).Render(strings.Join(item.Patterns, ", "))

		display = lipgloss.JoinHorizontal(lipgloss.Center, number, itemName, rarity, selectedPatternText)
	} else {
		number := utils.TextStyle.Width(3).MarginRight(1).Render(strconv.Itoa(item.Number))
		itemName := utils.TextStyle.Width(d.MaxNameLength).MarginRight(1).Render(item.Name)
		rarity := utils.TextStyle.Width(d.MaxRarityLength).MarginRight(1).Render(item.Rarity)
		patternText := utils.TextStyle.Width(d.MaxPatternLength).Render(strings.Join(item.Patterns, ", "))

		display = lipgloss.JoinHorizontal(lipgloss.Center, number, itemName, rarity, patternText)
	}

	io.WriteString(w, display)
}
