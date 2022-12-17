package countchars

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/Songmu/flextime"
	"github.com/westlab/glory"
)

func Handler(ctx context.Context) error {
	t := flextime.Now()
	config, err := glory.LoadConfig("../config.json")
	if err != nil {
		return fmt.Errorf("load config error: %w", err)
	}
	// 最終更新日の取得(決め打ち)
	for _, wg := range config.WorkingGroups {
		if wg.DirID == "" {
			log.Printf("WG %s does not have dir_name, skip", wg.Title)
			continue
		}
		for _, m := range wg.Members {
			dirPath := wg.DirID + "/" + m
			log.Printf("start fetch latest file from %s\n", dirPath)

			status, err := exec.Command("./fetch_latest_file.sh", dirPath).Output()
			if err != nil {
				return fmt.Errorf("fetch latest file error: %v", err)
			}
			if strings.Contains(string(status), "no docx file") {
				continue
			}

			data := strings.Split(string(status), "\n")
			count, err := strconv.ParseInt(data[2], 10, 64)
			if err != nil {
				return err
			}

			authorID, err := calcAuthorID(m, wg.Title)
			if err != nil {
				return fmt.Errorf("calc author id error: %w", err)
			}

			lastMod, err := time.Parse(time.RFC3339Nano, data[1])
			if err != nil {
				return fmt.Errorf("last mod time parse error: %w", err)
			}
			log.Printf("member: %s, updateTime: %v, ", m, lastMod)

			err = CreateThesisHistory(authorID, count, lastMod.UTC().Format(glory.TimeFormat), t.UTC().Format(glory.TimeFormat))
			if err != nil {
				return fmt.Errorf("insert thesis history error: %w", err)
			}

		}
	}
	return nil
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
