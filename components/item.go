package components

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/altugbakan/card-logger/utils"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
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

func (d SetItemDelegate) Height() int                             { return 1 }
func (d SetItemDelegate) Spacing() int                            { return 0 }
func (d SetItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d SetItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(SetItem)
	if !ok {
		return
	}

	progressBar := progress.New(progress.WithScaledGradient(utils.Gray, utils.LightBlue))
	progressBar.Width = 20
	percent := float64(item.Owned) / float64(item.Total)
	owned := strconv.Itoa(item.Owned)
	total := strconv.Itoa(item.Total)

	spaces := max(d.MaxNameLength-len(item.Name), 2)

	display := item.Name + strings.Repeat(" ", spaces) + progressBar.ViewAs(percent) + "  " + owned + "/" + total

	if index == m.Index() {
		display = utils.ActionStyle.Render("> " + display)
	} else {
		display = utils.TextStyle.Render("  " + display)
	}

	fmt.Fprint(w, display)
}
