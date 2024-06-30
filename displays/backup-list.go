package displays

import (
	"github.com/altugbakan/card-logger/db"
	"github.com/altugbakan/card-logger/items"
	"github.com/altugbakan/card-logger/keymaps"
	"github.com/altugbakan/card-logger/utils"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type BackupList struct {
	keyMap         keymaps.BackupList
	list           list.Model
	previousScreen Displayer
}

func NewBackupListScreen(previousScreen Displayer) BackupList {
	keyMap := keymaps.NewBackupListKeyMap()

	backups, maxNameLength, err := getBackupItems()
	if err != nil {
		utils.LogError("could not get all backups: %v", err)
	}

	list := utils.NewList(backups, items.BackupDelegate{MaxNameLength: maxNameLength}, "backup")

	if len(backups) == 0 {
		keyMap.Select.SetEnabled(false)
	}

	return BackupList{
		keyMap:         keyMap,
		list:           list,
		previousScreen: previousScreen,
	}
}

func (s BackupList) Update(msg tea.Msg) (Displayer, tea.Cmd) {
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
			res := s.restoreBackup()
			switch res := res.(type) {
			case utils.ErrorMessage:
				utils.LogError(res.Text)
			case utils.InfoMessage:
				utils.LogInfo(res.Text)
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

	backups := make([]list.Item, len(allBackups))
	maxNameLength := 0
	for i, backup := range allBackups {
		item := items.NewBackupItem(backup)
		maxNameLength = max(maxNameLength, len(backup))
		backups[i] = item
	}

	return backups, maxNameLength, nil
}

func (s *BackupList) restoreBackup() utils.Renderer {
	i, ok := s.list.SelectedItem().(items.Backup)
	if ok {
		err := db.RestoreBackup(i.Name)
		if err != nil {
			return utils.NewErrorMessage("could not restore backup: %v", err)
		} else {
			return utils.NewInfoMessage("backup restored")
		}
	}
	return utils.NewErrorMessage("could not restore backup due to an error")
}
