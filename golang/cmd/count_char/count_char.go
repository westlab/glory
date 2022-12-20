package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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

	cmdStr1 := "unzip -p "
	cmdStr2 := " word/document.xml | xmllint --format - | grep -E '<w:t[\"-~\\ ]*?>'| sed -E 's/<\\/?w:t[\"-=\\?-~\\ ]*>//g'| tr -d '\\r' | tr -d '\\n' | sed -E 's/\\ {2,}//g' | sed -E 's/&(gt|lt);/@/g'| tr -d ' ' |tr -d ' 　' | wc -m | tr -d ' '"
	cmdRes, err := exec.Command("sh", "-c", cmdStr1+filePath+cmdStr2).Output()
	if err != nil {
		return -1, time.Time{}, err
	}

	countStr := string(cmdRes)
	countStr = strings.Replace(countStr, "\n", "", -1)
	count, err := strconv.Atoi(countStr)
	if err != nil {
		return -1, time.Time{}, err
	}

	file, err := os.Stat(filePath)
	if err != nil {
		return -1, time.Time{}, err
	}
	modified := file.ModTime()

	return count, modified, nil
}
