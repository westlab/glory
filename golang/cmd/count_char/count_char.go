package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"code.sajari.com/docconv"
)

var NoDocxError = errors.New("no docx file")

func DirWalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, DirWalk(filepath.Join(dir, file.Name()))...)
			continue
		}
		if file.Name()[0] != '.' {
			paths = append(paths, filepath.Join(dir, file.Name()))
		}
	}

	return paths
}

func FetchLatestDocx(dirPath string) (string, error) {

	paths := DirWalk(dirPath)

	latestDocx := ""
	latestTime := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.Local)

	for _, path := range paths {

		fileExt := filepath.Ext(path)
		if fileExt != ".docx" {
			continue
		}

		file, err := os.Stat(path)
		if err != nil {
			return "", fmt.Errorf("failed to get file stat '%s' : %v", path, err)
		}

		modifiedTime := file.ModTime()
		if modifiedTime.After(latestTime) {
			latestDocx = path
			latestTime = modifiedTime
		}

	}

	if latestDocx == "" {
		return "", NoDocxError
	}

	return latestDocx, nil
}

func CountCharsInDocx(filePath string) (int, time.Time, error) {

	if filepath.Ext(filePath) != ".docx" {
		return -1, time.Time{}, fmt.Errorf("%s is not docx file", filePath)
	}

	// get last modified time
	file, err := os.Stat(filePath)
	if err != nil {
		return -1, time.Time{}, err
	}
	modified := file.ModTime()

	// get count char
	f, err := os.Open(filePath)
	if err != nil {
		return -1, time.Time{}, fmt.Errorf("fail to open file: %v", err)
	}
	defer f.Close()
	document, _, err := docconv.ConvertDocx(f)
	if err != nil {
		return -1, time.Time{}, err
	}
	document = strings.ReplaceAll(document, " ", "")
	document = strings.ReplaceAll(document, "\n", "")
	count := strings.Count(document, "")

	return count, modified, nil
}
