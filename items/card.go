package items

import (
	"fmt"
	"io"

	"github.com/altugbakan/card-logger/utils"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Card struct {
	CardID   int
	Number   int
	Name     string
	Patterns []utils.Pattern
}

func (c Card) FilterValue() string {
	return c.Name
}

type CardDelegate struct {
	MaxNameLength    int
	MaxPatternLength int
	SelectedIndex    int
}

func (d CardDelegate) Height() int                               { return 1 }
func (d CardDelegate) Spacing() int                              { return 0 }
func (d CardDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d CardDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(Card)
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

	io.WriteString(w, display)
}
