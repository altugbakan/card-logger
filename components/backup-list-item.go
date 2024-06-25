package components

import (
	"io"

	"github.com/altugbakan/card-logger/utils"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	BackupListItemHeight = 1
)

type BackupListItem struct {
	Name string
}

func NewBackupItem(name string) BackupListItem {
	return BackupListItem{Name: name}
}

func (i BackupListItem) FilterValue() string {
	return i.Name
}

type BackupListItemDelegate struct {
	MaxNameLength int
}

func (d BackupListItemDelegate) Height() int                               { return utils.ListItemHeight }
func (d BackupListItemDelegate) Spacing() int                              { return utils.ListItemSpacing }
func (d BackupListItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d BackupListItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(BackupListItem)
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
