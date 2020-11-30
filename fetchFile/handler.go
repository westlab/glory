package fetchFile

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Songmu/flextime"
	"github.com/westlab/glory"
	"google.golang.org/api/drive/v3"
	"log"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

func Handler(ctx context.Context, srv *drive.Service) error {
	t := flextime.Now()
	config, err := glory.LoadConfig("../config.json")
	if err != nil {
		return fmt.Errorf("load config error: %w", err)
	}
	// 最終更新日の取得(決め打ち)
	threshold := t.Add(time.Duration(-config.UpdateIntervalMin) * time.Minute)
	for _, wg := range config.WorkingGroups {
		if wg.DirID == "" {
			log.Printf("WG %s does not have dir_name, skip", wg.Title)
			continue
		}
		subdirs, err := FetchSubDirectories(ctx, srv, wg.DirID)
		if err != nil {
			return fmt.Errorf("fetch subdirectories error, WG %s : %w", wg.Title, err)
		}
		for _, d := range subdirs {
			if !isMember(d.Name, wg) {
				continue
			}
			log.Printf("File: name %s, ID %s\n", d.Name, d.Id)
			docxFile, err := FetchLatestDocx(ctx, srv, d.Id, threshold.UTC())
			if err != nil {
				if err == ErrNotFound {
					log.Print("recent file not found")
					continue
				}
				return err
			}
			//count characters
			out, err := exec.Command("./MSword_counter.sh", path.Join("workspace", docxFile.Name)).Output()
			if err != nil {
				return fmt.Errorf("word count error: %w", err)
			}
			exec.Command("rm", "-f", path.Join("workspace", docxFile.Name)).Run()

			if strings.Contains(string(out), "not MS word file") {
				log.Printf("%v is not a MS word file", docxFile.Name)
				continue
			}
			count, err := strconv.ParseInt(strings.TrimRight(string(out), "\n"), 10, 64)
			if err != nil {
				return err
			}
			log.Printf("%s: %d characters", docxFile.Name, count)
			authorID, err := calcAuthorID(d.Name, wg.Title)
			if err != nil {
				return fmt.Errorf("calc author id error: %w", err)
			}
			LastMod, err := time.Parse(time.RFC3339, docxFile.ModifiedTime)
			if err != nil {
				return fmt.Errorf("last mod time parse error: %w", err)
			}
			err = CreateThesisHistory(authorID, count, LastMod.UTC().Format(glory.TimeFormat), t.UTC().Format(glory.TimeFormat))
			if err != nil {
				return fmt.Errorf("insert thesis history error: %w", err)
			}
		}
	}
	return nil
}

func isMember(name string, wg *glory.WorkingGroup) bool {
	for _, m := range wg.Members {
		if name == m {
			return true
		}
	}
	return false
}

func calcAuthorID(name, wg string) (int64, error) {
	id, err := FetchAuthorID(name, wg)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return -1, err
		}
		id, err = CreateAuthorID(name, wg)
		if err != nil {
			return -1, err
		}
	}
	return id, nil
}
