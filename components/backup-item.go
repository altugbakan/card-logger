package components

import (
	"io"

	"github.com/altugbakan/card-logger/utils"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	BackupItemHeight = 1
)

type BackupItem struct {
	Name string
}

func NewBackupItem(name string) BackupItem {
	return BackupItem{Name: name}
}

func (i BackupItem) FilterValue() string {
	return i.Name
}

type BackupItemDelegate struct {
	MaxNameLength int
}

func (d BackupItemDelegate) Height() int                               { return BackupItemHeight }
func (d BackupItemDelegate) Spacing() int                              { return 0 }
func (d BackupItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d BackupItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(BackupItem)
	if !ok {
		return
	}

	display := utils.EmptyStyle.Width(d.MaxNameLength).Render(item.Name)
	if index == m.Index() {
		display = utils.ActionStyle.Render("> " + display)
	} else {
		display = utils.TextStyle.Render("  " + display)
	}

	io.WriteString(w, display)
}
