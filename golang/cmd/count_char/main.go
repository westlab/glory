package main

import (
	"github.com/Songmu/flextime"
	"github.com/joho/godotenv"
	"github.com/westlab/glory/config"
	"github.com/westlab/glory/db"
	"log"
	"os"
)

func main() {

	log.Print("[INFO] Start countChar process")

	envFile := "/opt/glory/.env"
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("[ERROR] Failed to load .env '%s' : %v", envFile, err)
	}

	done, err := db.InitializeDB(os.Getenv("COUNT_CHAR_DSN"))
	if err != nil {
		log.Fatalf("[ERROR] Failed to initialize db: %v", err)
	}
	defer done()

	gloryConfig, err := config.LoadConfig("/opt/glory/config.json")
	if err != nil {
		log.Fatalf("[ERROR] Failed to load config: %v", err)
	}

	// 最終更新日の取得(決め打ち)
	t := flextime.Now()
	for _, wg := range gloryConfig.WorkingGroups {

		if wg.DirName == "" {
			log.Printf("[INFO] WG %s does not have dir_name, skip", wg.Title)
			continue
		}

		for _, m := range wg.Members {

			dirPath := wg.DirName + "/" + m
			log.Printf("[INFO] Start FetchLatestDocx from %s", dirPath)

			latestDocx, err := FetchLatestDocx(dirPath)
			if err == NoDocxError {
				log.Printf("[INFO] No docx file in %s, skip", dirPath)
				continue
			} else if err != nil {
				log.Fatalf("[ERROR] FetchLatestDocx failed in '%s'", dirPath)
			}

			log.Printf("[INFO] %s's latest docx file is %s modified at %s", m, latestDocx, t)

			authorID, err := calcAuthorID(m, wg.Title)
			if err != nil {
				log.Fatalf("[ERROR] calcAuthorID failed: %v", err)
			}

			count, modified, err := CountCharsInDocx(latestDocx)
			if err != nil {
				log.Fatalf("[ERROR] CountCharsInDocx failed with '%s': %v", latestDocx, err)
			}

			err = CreateThesisHistory(authorID, int64(count), modified.UTC().Format(db.TimeFormat), t.UTC().Format(db.TimeFormat))
			if err != nil {
				log.Fatalf("[ERROR] insert thesis history error: %v", err)
			}

			log.Printf("[INFO] member: %s, updateTime: %v, numChar: %d", m, modified, count)

		}
	}

}
