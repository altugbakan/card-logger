package screens

import (
	"github.com/altugbakan/card-logger/db"
	"github.com/altugbakan/card-logger/keymaps"
	"github.com/altugbakan/card-logger/utils"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Initialize struct {
	keyMap        keymaps.Initialize
	isDownloading bool
	spinner       spinner.Model
}

func NewInitializeScreen() Initialize {
	keyMap := keymaps.NewInitializeKeyMap()
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = utils.ActionStyle

	return Initialize{
		keyMap:  keyMap,
		spinner: s,
	}
}

func (s Initialize) Update(msg tea.Msg) (Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keyMap.Quit):
			return s, tea.Quit
		case key.Matches(msg, s.keyMap.Load):
			utils.LogInfo("initializing by loading backup")
			return NewBackupListScreen(s), nil
		case key.Matches(msg, s.keyMap.Download):
			if s.isDownloading {
				return s, nil
			}
			utils.LogInfo("initializing by downloading database")
			s.startDownloading()
			return s, tea.Batch(utils.FetchLatestRelease, s.spinner.Tick)
		}
	case utils.DownloadCompleteMsg:
		utils.LogInfo("database download complete")
		db.Init()
		return NewTitleScreen(), nil
	}

	var cmd tea.Cmd
	s.spinner, cmd = s.spinner.Update(msg)
	return s, cmd
}

func (s Initialize) View() string {
	if s.isDownloading {
		s.spinner.Tick()
		title := utils.TitleStyle.Render("Downloading Database")
		spinner := s.spinner.View()
		message := utils.TextStyle.MarginLeft(1).Render("Please wait...")

		spinnerAndMessage := lipgloss.JoinHorizontal(lipgloss.Center, spinner, message)

		return lipgloss.JoinVertical(lipgloss.Center, title, spinnerAndMessage)
	}

	title := utils.TitleStyle.Render("Card Logger")
	message := utils.TextStyle.MarginBottom(1).Render("No database found. Choose an action:")
	download := lipgloss.JoinHorizontal(lipgloss.Center, utils.ActionStyle.Render("[D]"),
		utils.TextStyle.Render("ownload database"))
	load := lipgloss.JoinHorizontal(lipgloss.Center, utils.ActionStyle.Render("[L]"),
		utils.TextStyle.Render("oad saved backup"))

	downloadAndLoad := lipgloss.JoinVertical(lipgloss.Left, download, load)

	return lipgloss.JoinVertical(lipgloss.Center, title, message, downloadAndLoad)
}

func (s Initialize) Help() string {
	return s.keyMap.Help()
}

func (s *Initialize) startDownloading() {
	s.isDownloading = true
	s.keyMap.Download.SetEnabled(false)
	s.keyMap.Load.SetEnabled(false)
	s.keyMap.Quit.SetEnabled(false)
}
