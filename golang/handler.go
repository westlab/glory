package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Songmu/flextime"
	"github.com/gin-gonic/gin"

	"github.com/westlab/glory/config"
	"github.com/westlab/glory/db"
)

var gloryConfig *config.Conf

func init() {
	var err error
	if gloryConfig, err = config.LoadConfig("./config.json"); err != nil {
		log.Fatal(err)
	}
}

func IndexHandler(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index", gin.H{
		"title":          "glory",
		"working_groups": *getWorkingGroup(),
		"deadlines":      *calcDeadlines(),
	})
}

func ProgressHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	for _, wg := range gloryConfig.WorkingGroups {
		if id == wg.Title {
			th, err := FetchAllHistory(id)
			if err != nil {
				log.Printf("fetch history error: %v", err)
			}
			chart := NewProgressChart(wg, th)
			ctx.HTML(http.StatusOK, "progress", gin.H{
				"title":          wg.Describe,
				"heading":        wg.Describe + "の進捗状況",
				"working_groups": *getWorkingGroup(),
				"chart":          chart,
			})
			return
		}
	}
	NotFoundHandler(ctx)
}

func NotFoundHandler(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "notFound", gin.H{
		"working_groups": *getWorkingGroup(),
	})
}

func getWorkingGroup() *[]string {
	var ret []string
	for _, wg := range gloryConfig.WorkingGroups {
		ret = append(ret, wg.Title)
	}
	return &ret
}

func calcDeadlines() *map[string]string {
	ret := map[string]string{}
	for _, wg := range gloryConfig.WorkingGroups {
		var left string
		deadline, err := time.Parse(db.DateFormat, wg.Deadline)
		if err != nil {
			log.Printf("deadline parse error: %v", err)
			left = "error"
		} else {
			duration := deadline.Sub(flextime.Now())
			hours0 := int(duration.Hours())
			left = strconv.Itoa(hours0 / 24)
		}
		ret[wg.Describe] = left
	}
	return &ret
}

func FetchAllHistory(wg string) ([]*ThesisHistoryJoinAuthor, error) {
	var ret []*ThesisHistoryJoinAuthor

	rows, _ := db.DB.Query("SELECT name, char_count, fetch_time FROM thesis_history JOIN author ON author.author_id=thesis_history.author_id WHERE working_group=? and fetch_time <= ? ORDER BY fetch_time", wg, flextime.Now().UTC().Format(db.TimeFormat))
	if db.IsNoRows(rows.Err()) {
		return nil, nil
	}
	defer rows.Close()

	for rows.Next() {
		var th ThesisHistoryJoinAuthor
		var ft string
		err := rows.Scan(&th.Name, &th.CharCount, &ft)
		if err != nil {
			return nil, fmt.Errorf("fetch record error: %w", err)
		}
		th.FetchTime, err = time.Parse(db.TimeFormat, ft)
		if err != nil {
			return nil, fmt.Errorf("convert to time error: %w", err)
		}
		th.FetchTime = th.FetchTime.Local()
		ret = append(ret, &th)
	}

	return ret, nil
}
