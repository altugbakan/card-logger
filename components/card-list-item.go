package components

import (
	"fmt"
	"io"

	"github.com/altugbakan/card-logger/db"
	"github.com/altugbakan/card-logger/utils"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type CardListItem struct {
	CardID   int
	Number   int
	Name     string
	Patterns []db.Pattern
}

func (c CardListItem) FilterValue() string {
	return c.Name
}

type CardListItemDelegate struct {
	MaxNameLength int
	SelectedIndex int
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
		selectedPatternText := utils.ActionStyle.Render(d.getSelectedPatternText(item.Patterns))

		display = lipgloss.JoinHorizontal(lipgloss.Center, number, itemName, selectedPatternText)
	} else {
		number := utils.TextStyle.Width(3).MarginRight(1).Render(fmt.Sprintf("%d", item.Number))
		itemName := utils.TextStyle.Width(d.MaxNameLength).MarginRight(1).Render(item.Name)
		patternText := utils.TextStyle.Render(d.getPatternText(item.Patterns))

		display = lipgloss.JoinHorizontal(lipgloss.Center, number, itemName, patternText)
	}

	fmt.Fprint(w, display)
}

func (d *CardListItemDelegate) getSelectedPatternText(patterns []db.Pattern) string {
	selectedIndex := min(d.SelectedIndex, len(patterns)-1)
	text := ""
	for i, pattern := range patterns {
		if i == selectedIndex {
			text += fmt.Sprintf(" > %s:%d", pattern.Name, pattern.Quantity)
		} else {
			text += fmt.Sprintf("   %s:%d", pattern.Name, pattern.Quantity)
		}
	}

	return text
}

func (d *CardListItemDelegate) getPatternText(patterns []db.Pattern) string {
	text := ""
	for _, pattern := range patterns {
		text += fmt.Sprintf("   %s:%d", pattern.Name, pattern.Quantity)
	}

	return text
}
