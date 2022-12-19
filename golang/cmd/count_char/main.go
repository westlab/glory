package main

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/Songmu/flextime"
	"github.com/westlab/glory/config"

	"github.com/westlab/glory/db"
)

func main() {

	log.Print("Start countChar process")

	done, err := db.InitializeDB(os.Getenv("DSN"))
	if err != nil {
		log.Fatalf("Failed to initialize db: %v", err)
	}
	defer done()

	gloryConfig, err := config.LoadConfig("/app/config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 最終更新日の取得(決め打ち)
	t := flextime.Now()
	for _, wg := range gloryConfig.WorkingGroups {

		if wg.DirName == "" {
			log.Printf("WG %s does not have dir_name, skip", wg.Title)
			continue
		}

		for _, m := range wg.Members {

			dirPath := wg.DirName + "/" + m
			log.Printf("start fetch latest file from %s", dirPath)

			status, err := exec.Command("/app/cmd/fetch_latest_file.bash", dirPath).Output()
			if err != nil {
				log.Fatalf("fetch latest file error: %v", err)
			}

			if strings.Contains(string(status), "no docx file") {
				log.Printf("No docx file in %s, skip", dirPath)
				continue
			}

			data := strings.Split(string(status), "\n")
			count, err := strconv.ParseInt(data[2], 10, 64)
			if err != nil {
				log.Fatalf("count char error: %v", err)
			}

			authorID, err := calcAuthorID(m, wg.Title)
			if err != nil {
				log.Fatalf("calc author id error: %v", err)
			}

			lastMod, err := time.Parse(time.RFC3339Nano, data[1])
			if err != nil {
				log.Fatalf("last mod time parse error: %v", err)
			}
			log.Printf("member: %s, updateTime: %v, ", m, lastMod)

			err = CreateThesisHistory(authorID, count, lastMod.UTC().Format(db.TimeFormat), t.UTC().Format(db.TimeFormat))
			if err != nil {
				log.Fatalf("insert thesis history error: %v", err)
			}

		}
	}

}
