package mails

import (
	"log"
	"os"
	"path/filepath"
	"wazuhIssues/internal/regex"
)

func FindFiles(dir string, reg string) ([]string, error) {
	filePaths := make([]string, 0)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err == nil && regex.UseRegex(path, reg) {
			filePaths = append(filePaths, path)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return filePaths, nil
}

func ReadMails(files []string) ([]string, error) {
	mails := make([]string, 0)
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}
		mails = append(mails, string(data))
	}
	return mails, nil
}
