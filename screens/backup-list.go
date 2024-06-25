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

const (
	backupListHeightMargin = 5
	backupListWidthMargin  = 2
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

	width, height := utils.GetWindowSize()

	width -= setListWidthMargin * 2
	height = min(getListHeight(len(backups)), height-setListHeightMargin*2-utils.TotalHelpWidth)
	utils.LogInfo("initializing backup list with size %d x %d", width, height)

	list := list.New(backups, components.BackupItemDelegate{MaxNameLength: maxNameLength}, width, height)
	list.SetShowHelp(false)
	list.SetShowTitle(false)
	list.FilterInput.Cursor.Style = utils.CursorStyle
	list.FilterInput.PromptStyle = utils.ActionStyle
	list.DisableQuitKeybindings()
	if len(backups) == 0 {
		list.SetShowStatusBar(false)
	}

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
		s.list.SetSize(msg.Width-backupListWidthMargin*2,
			min(getListHeight(len(s.list.Items())), msg.Height-backupListHeightMargin*2))
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

func getListHeight(itemCount int) int {
	if itemCount == 0 {
		return 1
	}
	return itemCount*components.BackupItemHeight + backupListHeightMargin
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
	i, ok := s.list.SelectedItem().(components.BackupItem)
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
