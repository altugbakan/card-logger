package db

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/altugbakan/card-logger/utils"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	downloadURL = "https://github.com/altugbakan/card-logger/releases/latest/download/"
)

type DownloadCompleteMsg struct{}
type DownloadFailedMsg struct{}

func FetchLatestRelease() tea.Msg {
	dbURL := downloadURL + fmt.Sprintf("cards_%s.zip", utils.GetConfig().Type)
	zipData, err := downloadFile(dbURL)
	if err != nil {
		utils.LogError("failed to download file: %v", err)
		return DownloadFailedMsg{}
	}
	utils.LogInfo("downloaded file from %s", dbURL)

	if err := unzipAndSave(zipData, getDatabasePath()); err != nil {
		utils.LogError("failed to unzip and save file: %v", err)
		return DownloadFailedMsg{}
	}

	return DownloadCompleteMsg{}
}

func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func unzipAndSave(data []byte, outputPath string) error {
	filename := "cards.db"
	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}

	for _, f := range r.File {
		if f.Name == filename {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			fileData, err := io.ReadAll(rc)
			if err != nil {
				return err
			}

			err = os.MkdirAll(getDatabaseDirectory(), 0755)
			if err != nil {
				return err
			}

			return os.WriteFile(outputPath, fileData, 0644)
		}
	}

	return fmt.Errorf("file %s not found in zip", filename)
}
