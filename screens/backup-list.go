package screens

import (
	"github.com/altugbakan/card-logger/components"
	"github.com/altugbakan/card-logger/db"
	"github.com/altugbakan/card-logger/keymaps"
	"github.com/altugbakan/card-logger/utils"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type BackupList struct {
	keyMap         keymaps.SetList
	list           list.Model
	previousScreen Screen
}

func NewBackupListScreen(previousScreen Screen) BackupList {
	keyMap := keymaps.NewListKeyMap()

	backups, maxNameLength, err := getBackupItems()
	if err != nil {
		utils.LogError("could not get all backups: %v", err)
	}

	list := utils.NewList(backups, components.BackupListItemDelegate{MaxNameLength: maxNameLength}, "backup")
	utils.LogInfo("filter input width of the list: %d", list.FilterInput.Width)

	return BackupList{
		keyMap:         keyMap,
		list:           list,
		previousScreen: previousScreen,
	}
}

func (s BackupList) Update(msg tea.Msg) (Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s.list.SettingFilter() {
			break
		}
		switch {
		case key.Matches(msg, s.keyMap.Back):
			if s.list.IsFiltered() {
				s.list.ResetFilter()
				return s, nil
			}
			return s.previousScreen, nil
		case key.Matches(msg, s.keyMap.Select):
			res, err := s.restoreBackup()
			if err != nil {
				utils.LogError("could not restore backup: %v", err)
			} else {
				utils.LogInfo("backup restored")
			}
			db.Reinit()
			return NewBackupScreen(WithMessage(res.Render())), nil
		}
	case tea.WindowSizeMsg:
		utils.SetListSize(&s.list, msg.Width, msg.Height)
		return s, nil
	}
	var cmd tea.Cmd
	s.list, cmd = s.list.Update(msg)
	return s, cmd
}

func (s BackupList) View() string {
	title := utils.TitleStyle.MarginBottom(1).Render("Load Backup")
	list := s.list.View()
	return lipgloss.JoinVertical(lipgloss.Center, title, list)
}

func (s BackupList) Help() string {
	return s.keyMap.Help()
}

func getBackupItems() ([]list.Item, int, error) {
	allBackups, err := db.ListBackups()
	if err != nil {
		return nil, 0, err
	}

	backups := []list.Item{}
	maxNameLength := 0
	for _, backup := range allBackups {
		item := components.NewBackupItem(backup)
		maxNameLength = max(maxNameLength, len(backup))
		backups = append(backups, item)
	}

	return backups, maxNameLength, nil
}

func (s *BackupList) restoreBackup() (utils.Message, error) {
	i, ok := s.list.SelectedItem().(components.BackupListItem)
	if ok {
		err := db.RestoreBackup(i.Name)
		if err != nil {
			return utils.NewErrorMessage("could not restore backup"), err
		} else {
			return utils.NewInfoMessage("backup restored"), nil
		}
	}
	return utils.NewErrorMessage("could not restore backup"), nil
}
