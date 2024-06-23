package utils

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"time"
)

func SaveBackup() (string, error) {
	fileName := fmt.Sprintf("cards_%s.db", time.Now().Format("2006-01-02_15-04-05"))
	destinationFilePath := BackupDirectory + "/" + fileName

	sourceFile, err := os.Open(DatabaseFilePath)
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

func RestoreBackup(fileName string) error {
	sourceFilePath := BackupDirectory + "/" + fileName
	destinationFile, err := os.Create(DatabaseFilePath)
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
	files, err := os.ReadDir(BackupDirectory)
	if err != nil {
		return nil, err
	}

	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	sort.Sort(sort.Reverse(sort.StringSlice(fileNames)))

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

	pattern := regexp.MustCompile(`^cards_\d{4}-\d{2}-\d{2}_\d{2}-\d{2}-\d{2}\.db$`)
	for _, file := range files {
		if pattern.MatchString(file) {
			return file, nil
		}
	}

	return "none", nil
}
