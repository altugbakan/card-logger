package db

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/altugbakan/card-logger/utils"
)

const dateTimeFormat = "2006-01-02_15-04-05"

var hasChanges bool = false

type backupFile struct {
	name     string
	dateTime time.Time
}

func SaveBackup() (string, error) {
	fileName := fmt.Sprintf("cards_%s.db", time.Now().Format(dateTimeFormat))
	destinationFilePath := filepath.Join(getBackupDirectory(), fileName)

	return fileName, saveBackup(destinationFilePath)
}

func SaveAutoBackup() {
	fileName, err := saveAutoBackup()
	if err != nil {
		utils.LogError("could not save auto backup: %v", err)
	} else if fileName != "" {
		utils.LogInfo("auto backup saved to %s", fileName)
	} else {
		utils.LogInfo("no changes detected, auto backup not saved")
	}
}

func RestoreBackup(fileName string) error {
	err := os.MkdirAll(getDatabaseDirectory(), 0755)
	if err != nil {
		return err
	}

	destinationFile, err := os.Create(getDatabasePath())
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	sourceFilePath := filepath.Join(getBackupDirectory(), fileName)
	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

func ListBackups() ([]string, error) {
	backupDirectory := getBackupDirectory()
	if _, err := os.Stat(backupDirectory); os.IsNotExist(err) {
		return []string{}, nil
	}

	files, err := os.ReadDir(backupDirectory)
	if err != nil {
		return nil, err
	}

	var backups []backupFile
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "cards") {
			dateTime := extractDateTimeFromName(file.Name())
			backups = append(backups, backupFile{name: file.Name(), dateTime: dateTime})
		}
	}

	sort.Slice(backups, func(i, j int) bool {
		return backups[i].dateTime.After(backups[j].dateTime)
	})

	var fileNames []string
	for _, backup := range backups {
		fileNames = append(fileNames, backup.name)
	}

	return fileNames, nil
}

func GetLatestBackup() (string, error) {
	files, err := ListBackups()
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "none", nil
	}

	return files[0], nil
}

func saveBackup(destinationFilePath string) error {
	err := os.MkdirAll(getBackupDirectory(), 0755)
	if err != nil {
		return err
	}

	sourceFile, err := os.Open(getDatabasePath())
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(destinationFilePath)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	hasChanges = false
	return nil
}

func saveAutoBackup() (string, error) {
	if !hasChanges {
		return "", nil
	}

	fileName := fmt.Sprintf("cards_auto_%s.db", time.Now().Format(dateTimeFormat))
	destinationFilePath := filepath.Join(getBackupDirectory(), fileName)

	err := saveBackup(destinationFilePath)
	if err != nil {
		return "", err
	}

	files, err := os.ReadDir(getBackupDirectory())
	if err != nil {
		return "", err
	}

	fileNames := []string{}
	for _, file := range files {
		if regexp.MustCompile(`^cards_auto_\d{4}-\d{2}-\d{2}_\d{2}-\d{2}-\d{2}\.db$`).MatchString(file.Name()) {
			fileNames = append(fileNames, file.Name())
		}
	}

	if len(fileNames) > 10 {
		sort.Strings(fileNames)
		oldestAutoBackup := filepath.Join(getBackupDirectory(), fileNames[0])
		err := os.Remove(oldestAutoBackup)
		if err != nil {
			utils.LogWarning("could not remove oldest auto backup %s: %v",
				oldestAutoBackup, err)
			return "", err
		}
		utils.LogInfo("removed oldest auto backup %s", oldestAutoBackup)
	}

	return fileName, nil
}

func extractDateTimeFromName(name string) time.Time {
	name = strings.TrimPrefix(name, "cards")
	name = strings.TrimPrefix(name, "_auto")
	name = strings.Split(name, ".")[0]
	dateTimePart := strings.TrimLeft(name, "_")

	t, err := time.Parse(dateTimeFormat, dateTimePart)
	if err != nil {
		return time.Time{}
	}
	return t
}

func getBackupDirectory() string {
	config := utils.GetConfig()
	return filepath.Join("backups", config.Type)
}
