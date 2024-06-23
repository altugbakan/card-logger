package screens

import (
	"github.com/altugbakan/card-logger/keymaps"
	"github.com/altugbakan/card-logger/utils"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Backup struct {
	keyMap       keymaps.Backup
	msg          string
	latestBackup string
}

func NewBackupScreen() Backup {
	keyMap := keymaps.NewBackupKeyMap()

	latestBackup, err := utils.GetLatestBackup()
	if err != nil {
		utils.LogWarning("could not get latest backup: %v", err)
	}

	return Backup{
		keyMap:       keyMap,
		latestBackup: latestBackup,
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
				utils.LogWarning("could not save backup: %v", err)
			} else {
				utils.LogInfo("backup saved")
			}
		case key.Matches(msg, s.keyMap.Load):
			//TODO: implement load backup
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
	s.latestBackup, err = utils.SaveBackup()
	if err != nil {
		return utils.NewErrorMessage("could not save backup"), err
	}
	return utils.NewInfoMessage("backup saved"), nil
}
