package components

import (
	"fmt"
	"io"
	"strconv"

	"github.com/altugbakan/card-logger/utils"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SetItem struct {
	Name  string
	Total int
	Owned int
}

func (i SetItem) FilterValue() string {
	return i.Name
}

type SetItemDelegate struct {
	MaxNameLength int
}

const (
	maxOwnedAndTotalLength = 7
	progressBarWidth       = 20
)

func (d SetItemDelegate) Height() int                               { return 1 }
func (d SetItemDelegate) Spacing() int                              { return 0 }
func (d SetItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d SetItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(SetItem)
	if !ok {
		return
	}

	progressBar := progress.New(progress.WithScaledGradient(utils.DarkLightBlue, utils.LightBlue))
	progressBar.Width = progressBarWidth

	percent := float64(item.Owned) / float64(item.Total)
	owned := strconv.Itoa(item.Owned)
	total := strconv.Itoa(item.Total)

	var display string
	if index == m.Index() {
		itemName := utils.ActionStyle.Width(d.MaxNameLength).Render(item.Name)
		ownedAndTotal := utils.ActionStyle.Width(maxOwnedAndTotalLength).
			AlignHorizontal(lipgloss.Right).MarginLeft(1).Render(owned + "/" + total)
		progressBar.PercentageStyle = utils.ActionStyle
		display = utils.ActionStyle.Render("> " + itemName + progressBar.ViewAs(percent) + ownedAndTotal)
	} else {
		itemName := utils.TextStyle.Width(d.MaxNameLength).Render(item.Name)
		ownedAndTotal := utils.TextStyle.Width(maxOwnedAndTotalLength).
			AlignHorizontal(lipgloss.Right).MarginLeft(1).Render(owned + "/" + total)
		progressBar.PercentageStyle = utils.TextStyle
		display = utils.TextStyle.Render("  " + itemName + progressBar.ViewAs(percent) + ownedAndTotal)
	}

	fmt.Fprint(w, display)
}
