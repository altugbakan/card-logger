package items

import (
	"fmt"
	"io"

	"github.com/altugbakan/card-logger/utils"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type NoPattern struct {
	Number int
	Name   string
	Rarity string
}

func (m NoPattern) FilterValue() string {
	return m.Name
}

type NoPatternDelegate struct {
	MaxNameLength   int
	MaxRarityLength int
}

func (d NoPatternDelegate) Height() int                               { return utils.ItemHeight }
func (d NoPatternDelegate) Spacing() int                              { return utils.ItemSpacing }
func (d NoPatternDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d NoPatternDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(NoPattern)
	if !ok {
		return
	}

	var display string
	if index == m.Index() {
		number := utils.ActionStyle.Width(3).MarginRight(1).Render(fmt.Sprintf("%d", item.Number))
		itemName := utils.ActionStyle.Width(d.MaxNameLength).MarginRight(1).Render(item.Name)
		rarity := utils.ActionStyle.Width(d.MaxRarityLength).Render(item.Rarity)
		display = utils.ActionStyle.Render("> " + number + itemName + rarity)
	} else {
		number := utils.TextStyle.Width(3).MarginRight(1).Render(fmt.Sprintf("%d", item.Number))
		itemName := utils.TextStyle.Width(d.MaxNameLength).MarginRight(1).Render(item.Name)
		rarity := utils.TextStyle.Width(d.MaxRarityLength).Render(item.Rarity)
		display = utils.TextStyle.Render("  " + number + itemName + rarity)
	}

	io.WriteString(w, display)
}
