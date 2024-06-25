package screens

import (
	"github.com/altugbakan/card-logger/db"
	"github.com/altugbakan/card-logger/keymaps"
	"github.com/altugbakan/card-logger/utils"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type BackupScreenOption func(*Backup)

type Backup struct {
	keyMap       keymaps.Backup
	msg          string
	latestBackup string
}

func NewBackupScreen(opts ...BackupScreenOption) Backup {
	keyMap := keymaps.NewBackupKeyMap()

	latestBackup, err := db.GetLatestBackup()
	if err != nil {
		utils.LogError("could not get latest backup: %v", err)
	}

	backup := Backup{
		keyMap:       keyMap,
		latestBackup: latestBackup,
	}

	for _, opt := range opts {
		opt(&backup)
	}

	return backup
}

func WithMessage(msg string) BackupScreenOption {
	return func(s *Backup) {
		s.msg = msg
	}
}

func (s Backup) Update(msg tea.Msg) (Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keyMap.Back):
			return NewTitleScreen(), nil
		case key.Matches(msg, s.keyMap.Save):
			res, err := s.saveBackup()
			s.msg = res.Render()
			if err != nil {
				utils.LogError("could not save backup: %v", err)
			}
		case key.Matches(msg, s.keyMap.Load):
			return NewBackupListScreen(s), nil
		}
	}
	return s, nil
}

func (s Backup) View() string {
	title := utils.TitleStyle.Render("Backup")
	latestBackup := utils.TextStyle.Render("latest backup: " + s.latestBackup)
	titleAndBackup := lipgloss.JoinVertical(lipgloss.Center, title, latestBackup)
	msg := utils.DimTextStyle.MarginTop(1).Render(s.msg)
	return lipgloss.JoinVertical(lipgloss.Center, titleAndBackup, msg)
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
