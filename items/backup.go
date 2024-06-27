package items

import (
	"io"

	"github.com/altugbakan/card-logger/utils"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	BackupListItemHeight = 1
)

type Backup struct {
	Name string
}

func NewBackupItem(name string) Backup {
	return Backup{Name: name}
}

func (i Backup) FilterValue() string {
	return i.Name
}

type BackupDelegate struct {
	MaxNameLength int
}

func (d BackupDelegate) Height() int                               { return utils.ItemHeight }
func (d BackupDelegate) Spacing() int                              { return utils.ItemSpacing }
func (d BackupDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d BackupDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(Backup)
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
