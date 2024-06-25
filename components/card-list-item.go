package components

import (
	"fmt"
	"io"

	"github.com/altugbakan/card-logger/utils"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type CardListItem struct {
	CardID   int
	Number   int
	Name     string
	Patterns []utils.Pattern
}

func (c CardListItem) FilterValue() string {
	return c.Name
}

type CardListItemDelegate struct {
	MaxNameLength    int
	MaxPatternLength int
	SelectedIndex    int
}

func (d CardListItemDelegate) Height() int                               { return 1 }
func (d CardListItemDelegate) Spacing() int                              { return 0 }
func (d CardListItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d CardListItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(CardListItem)
	if !ok {
		return
	}

	var display string
	if index == m.Index() {
		number := utils.ActionStyle.Width(3).MarginRight(1).Render(fmt.Sprintf("%d", item.Number))
		itemName := utils.ActionStyle.Width(d.MaxNameLength).MarginRight(1).Render(item.Name)
		selectedPatternText := utils.ActionStyle.Width(d.MaxPatternLength).Render(utils.GetSelectedPatternItemText(d.SelectedIndex, item.Patterns))

		display = lipgloss.JoinHorizontal(lipgloss.Center, number, itemName, selectedPatternText)
	} else {
		number := utils.TextStyle.Width(3).MarginRight(1).Render(fmt.Sprintf("%d", item.Number))
		itemName := utils.TextStyle.Width(d.MaxNameLength).MarginRight(1).Render(item.Name)
		patternText := utils.TextStyle.Width(d.MaxPatternLength).Render(utils.GetPatternItemText(item.Patterns))

		display = lipgloss.JoinHorizontal(lipgloss.Center, number, itemName, patternText)
	}

	fmt.Fprint(w, display)
}
