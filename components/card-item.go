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

type CardItem struct {
	CardID   int
	Number   int
	Name     string
	Patterns []db.Pattern
}

func (c CardItem) FilterValue() string {
	return c.Name
}

type CardItemDelegate struct {
	MaxNameLength int
	SelectedIndex int
}

func (d CardItemDelegate) Height() int                               { return 1 }
func (d CardItemDelegate) Spacing() int                              { return 0 }
func (d CardItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d CardItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(CardItem)
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

func (d *CardItemDelegate) getSelectedPatternText(patterns []db.Pattern) string {
	text := ""
	for i, pattern := range patterns {
		if i == d.SelectedIndex {
			text += fmt.Sprintf(" > %s:%d", pattern.Name, pattern.Quantity)
		} else {
			text += fmt.Sprintf("   %s:%d", pattern.Name, pattern.Quantity)
		}
	}

	return text
}

func (d *CardItemDelegate) getPatternText(patterns []db.Pattern) string {
	text := ""
	for _, pattern := range patterns {
		text += fmt.Sprintf("   %s:%d", pattern.Name, pattern.Quantity)
	}

	return text
}
