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
	backupHeightMargin = 5
	backupWidthMargin  = 2
)

type Backup struct {
	keyMap       keymaps.Backup
	msg          string
	latestBackup string
	list         list.Model
	showBackups  bool
}

func NewBackupScreen() Backup {
	keyMap := keymaps.NewBackupKeyMap()

	latestBackup, err := db.GetLatestBackup()
	if err != nil {
		utils.LogError("could not get latest backup: %v", err)
	}

	backups, maxNameLength, err := getBackupItems()
	if err != nil {
		utils.LogError("could not get all backups: %v", err)
	}

	width, height := utils.GetWindowSize()

	width -= listWidthMargin * 2
	height = min(getListHeight(len(backups)), height-listHeightMargin*2-utils.TotalHelpWidth)
	utils.LogInfo("initializing backup with size %d x %d", width, height)

	list := list.New(backups, components.BackupItemDelegate{MaxNameLength: maxNameLength}, width, height)
	list.SetShowHelp(false)
	list.SetShowTitle(false)
	list.FilterInput.Cursor.Style = utils.CursorStyle
	list.FilterInput.PromptStyle = utils.ActionStyle
	list.DisableQuitKeybindings()

	return Backup{
		keyMap:       keyMap,
		latestBackup: latestBackup,
		list:         list,
		showBackups:  false,
	}
}

func (s Backup) Update(msg tea.Msg) (Screen, tea.Cmd) {
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
			if s.showBackups {
				s.setShowBackups(false)
				return s, nil
			}
			return NewTitleScreen(), nil
		case key.Matches(msg, s.keyMap.Save):
			res, err := s.saveBackup()
			s.msg = res.Render()
			if err != nil {
				utils.LogError("could not save backup: %v", err)
			}
		case key.Matches(msg, s.keyMap.Load):
			s.msg = s.setShowBackups(true).Render()
			return s, nil
		case key.Matches(msg, s.keyMap.Select):
			res, err := s.restoreBackup()
			s.msg = res.Render()
			if err != nil {
				utils.LogError("could not restore backup: %v", err)
			} else {
				utils.LogInfo("backup restored")
			}
			s.setShowBackups(false)
			return s, nil
		}
	case tea.WindowSizeMsg:
		s.list.SetSize(msg.Width-backupWidthMargin*2,
			min(getListHeight(len(s.list.Items())), msg.Height-backupHeightMargin*2))
		return s, nil
	}

	var cmd tea.Cmd
	s.list, cmd = s.list.Update(msg)
	return s, cmd
}

func (s Backup) View() string {
	if s.showBackups {
		title := utils.TitleStyle.MarginBottom(1).Render("Load Backup")
		list := s.list.View()
		return lipgloss.JoinVertical(lipgloss.Center, title, list)
	} else {
		title := utils.TitleStyle.Render("Backup")
		latestBackup := utils.TextStyle.Render("latest backup: " + s.latestBackup)
		titleAndBackup := lipgloss.JoinVertical(lipgloss.Center, title, latestBackup)
		msg := utils.DimTextStyle.MarginTop(1).Render(s.msg)
		return lipgloss.JoinVertical(lipgloss.Center, titleAndBackup, msg)
	}
}

func (s Backup) Help() string {
	return s.keyMap.Help()
}

func (s *Backup) saveBackup() (utils.Message, error) {
	var err error
	s.latestBackup, err = db.SaveBackup()
	if err != nil {
		return utils.NewErrorMessage("could not save backup"), err
	}
	utils.LogInfo("saved backup %s", s.latestBackup)
	return utils.NewInfoMessage("backup saved"), nil
}

func (s *Backup) restoreBackup() (utils.Message, error) {
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

func (s *Backup) setShowBackups(show bool) utils.Message {
	if show {
		items, _, err := getBackupItems()
		if err != nil {
			utils.LogError("could not load backups: %v", err)
			return utils.NewErrorMessage("could not load backups")
		}
		s.list.SetItems(items)
		if len(items) == 0 {
			s.list.SetShowStatusBar(false)
		}
		s.keyMap.Load.SetEnabled(false)
		s.keyMap.Save.SetEnabled(false)
		s.keyMap.Up.SetEnabled(true)
		s.keyMap.Down.SetEnabled(true)
		s.keyMap.Select.SetEnabled(true)
		s.keyMap.Search.SetEnabled(true)
		s.showBackups = true
	} else {
		s.keyMap.Load.SetEnabled(true)
		s.keyMap.Save.SetEnabled(true)
		s.keyMap.Up.SetEnabled(false)
		s.keyMap.Down.SetEnabled(false)
		s.keyMap.Select.SetEnabled(false)
		s.keyMap.Search.SetEnabled(false)
		s.showBackups = false
		s.list.SetShowStatusBar(true)
	}

	return utils.NewInfoMessage("")
}

func getListHeight(itemCount int) int {
	if itemCount == 0 {
		return 1
	}
	return itemCount*components.BackupItemHeight + backupHeightMargin
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
