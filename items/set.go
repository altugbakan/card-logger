package items

import (
	"io"
	"strconv"

	"github.com/altugbakan/card-logger/utils"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Set struct {
	Abbr  string
	Name  string
	Total int
	Owned int
}

func (i Set) FilterValue() string {
	return i.Name
}

type SetDelegate struct {
	MaxNameLength int
}

const (
	maxOwnedAndTotalLength = 7
	progressBarWidth       = 20
)

func (d SetDelegate) Height() int                               { return utils.ItemHeight }
func (d SetDelegate) Spacing() int                              { return utils.ItemSpacing }
func (d SetDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d SetDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(Set)
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

	io.WriteString(w, display)
}
