package db

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/altugbakan/card-logger/utils"
)

var hasChanges bool = false

const dateTimeFormat = "2006-01-02_15-04-05"

type backupFile struct {
	name     string
	dateTime time.Time
}

type byDateTime []backupFile

func (s byDateTime) Len() int {
	return len(s)
}

func (s byDateTime) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byDateTime) Less(i, j int) bool {
	return s[i].dateTime.After(s[j].dateTime)
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

func SaveBackup() (string, error) {
	fileName := fmt.Sprintf("cards_%s.db", time.Now().Format(dateTimeFormat))
	destinationFilePath := backupDirectory + "/" + fileName

	return saveBackup(destinationFilePath, fileName)
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
	sourceFilePath := backupDirectory + "/" + fileName
	destinationFile, err := os.Create(databaseFilePath)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

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

	sort.Sort(byDateTime(backups))

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

func saveBackup(destinationFilePath string, fileName string) (string, error) {
	sourceFile, err := os.Open(databaseFilePath)
	if err != nil {
		return "", err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(destinationFilePath)
	if err != nil {
		return "", err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func saveAutoBackup() (string, error) {
	if !hasChanges {
		return "", nil
	}

	fileName := fmt.Sprintf("cards_auto_%s.db", time.Now().Format(dateTimeFormat))
	destinationFilePath := backupDirectory + "/" + fileName

	fileName, err := saveBackup(destinationFilePath, fileName)
	if err != nil {
		return "", err
	}

	files, err := os.ReadDir(backupDirectory)
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
		oldestAutoBackup := fileNames[0]
		err := os.Remove(backupDirectory + "/" + oldestAutoBackup)
		if err != nil {
			utils.LogWarning("could not remove oldest auto backup %s: %v",
				oldestAutoBackup, err)
			return "", err
		}
		utils.LogInfo("removed oldest auto backup %s", oldestAutoBackup)
	}

	return fileName, nil
}
