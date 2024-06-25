package utils

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	downloadURL = "https://github.com/altugbakan/card-logger/releases/latest/download/cards.db.zip"
)

type DownloadCompleteMsg struct{}
type DownloadFailedMsg struct{}

func FetchLatestRelease() tea.Msg {
	zipData, err := downloadFile(downloadURL)
	if err != nil {
		LogError("failed to download file: %v", err)
		return DownloadFailedMsg{}
	}

	if err := unzipAndSave(zipData, DatabaseFilePath); err != nil {
		LogError("failed to unzip and save file: %v", err)
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
	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}

	for _, f := range r.File {
		if f.Name == "cards.db" {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			fileData, err := io.ReadAll(rc)
			if err != nil {
				return err
			}

			return os.WriteFile(outputPath, fileData, 0644)
		}
	}

	return fmt.Errorf("cards.db not found in zip")
}
