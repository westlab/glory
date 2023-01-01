package main

import (
	"fmt"
	"github.com/westlab/glory/config"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Songmu/flextime"
	"github.com/gin-gonic/gin"

	"github.com/westlab/glory/db"
)

func IndexHandler(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index", gin.H{
		"title":          "glory",
		"working_groups": *getWorkingGroup(),
		"deadlines":      *calcDeadlines(),
	})
}

func ProgressHandler(ctx *gin.Context) {

	id := ctx.Param("id")
	for _, wg := range GloryConfig.WorkingGroups {
		if id == wg.Title {

			th, err := FetchAllHistory(id)
			if err != nil {
				log.Printf("fetch history error: %v", err)
			}

			chart := NewProgressChart(wg, th)
			ranking, err := FetchRanking(id)
			if err != nil {
				log.Printf("fetchRanking error: %v", err)
			}

			ctx.HTML(http.StatusOK, "progress", gin.H{
				"title":          wg.Describe,
				"heading":        wg.Describe + "の進捗状況",
				"working_groups": *getWorkingGroup(),
				"chart":          chart,
				"ranking":        ranking,
			})
			return
		}
	}
	NotFoundHandler(ctx)
}

func NotFoundHandler(ctx *gin.Context) {
	ctx.HTML(http.StatusNotFound, "notFound", gin.H{
		"working_groups": *getWorkingGroup(),
	})
}

func getWorkingGroup() *[]string {
	var ret []string
	for _, wg := range GloryConfig.WorkingGroups {
		ret = append(ret, wg.Title)
	}
	return &ret
}

func calcDeadlines() *map[string]string {
	ret := map[string]string{}
	for _, wg := range GloryConfig.WorkingGroups {
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

func FetchAllHistory(wgTitle string) ([]*ThesisHistoryJoinAuthor, error) {

	var ret []*ThesisHistoryJoinAuthor

	query := `SELECT name, char_count, fetch_time
		FROM thesis_history JOIN author ON author.author_id=thesis_history.author_id
		WHERE working_group=? and fetch_time <= ? ORDER BY fetch_time`

	rows, err := db.DB.Query(query, wgTitle, flextime.Now().UTC().Format(db.TimeFormat))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if db.IsNoRows(rows.Err()) {
		return nil, nil
	}

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

func NewProgressChart(wg *config.WorkingGroup, th []*ThesisHistoryJoinAuthor) *ProgressChart {
	ret := ProgressChart{
		Maxi:     defaultMaxi,
		StepSize: defaultMaxi / steps,
	}
	authorIdx := map[string]int{}
	for i, v := range wg.Members {
		author := AuthorData{
			Name:  v,
			Color: colors[i%len(colors)],
		}
		ret.AuthorDataList = append(ret.AuthorDataList, &author)
		authorIdx[v] = i
	}
	if len(th) == 0 {
		ret.Dates = []string{toDateStr(flextime.Now().AddDate(0, 0, -1)),
			toDateStr(flextime.Now())}
		return &ret
	}
	ret.Dates = makeDates(th[0].FetchTime.AddDate(0, 0, -1))
	for _, v := range ret.AuthorDataList {
		v.Data = make([]int, len(ret.Dates))
	}
	idx := 1
	for _, t := range th {
		for ret.Dates[idx] != toDateStr(t.FetchTime) {
			idx++
		}
		authorNumber, ok := authorIdx[t.Name]
		if !ok {
			continue
		}
		ret.AuthorDataList[authorNumber].Data[idx] = t.CharCount
		if t.CharCount > ret.Maxi {
			ret.Maxi = t.CharCount
		}
	}
	for _, author := range ret.AuthorDataList {
		for i := 1; i < len(ret.Dates); i++ {
			if author.Data[i] == 0 {
				author.Data[i] = author.Data[i-1]
			}
		}
	}
	ret.Maxi = fixMaxi(ret.Maxi)
	ret.StepSize = ret.Maxi / steps

	return &ret
}

func toDateStr(t time.Time) string {
	return t.Format(dateFormat)
}

func fixMaxi(m int) int {
	ret := int(float64(m) * 1.1)
	for ret >= 10 {
		ret /= 10
	}
	ret++
	for ret <= m {
		ret *= 10
	}
	return ret
}

func makeDates(start time.Time) []string {
	var ret []string
	d := start
	for d.Before(flextime.Now().AddDate(0, 0, 1)) {
		ret = append(ret, toDateStr(d))
		d = d.AddDate(0, 0, 1)
	}
	return ret
}

func FetchRanking(wgTitle string) ([]*RankingRow, error) {
	var ret []*RankingRow

	query := `SELECT b.name, a.char_count,
		    ROW_NUMBER() OVER(ORDER BY a.char_count DESC)
		FROM (
		    SELECT
		        author_id,
		        char_count,
		        ROW_NUMBER() OVER(PARTITION BY author_id ORDER BY fetch_time DESC) latest_order
		    FROM thesis_history 
		) a JOIN author b ON a.author_id = b.author_id
		WHERE latest_order = 1 AND b.working_group=?
		ORDER BY a.char_count DESC`

	rows, err := db.DB.Query(query, wgTitle)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if db.IsNoRows(rows.Err()) {
		return nil, nil
	}

	for rows.Next() {

		var th RankingRow
		err := rows.Scan(&th.Name, &th.CharCount, &th.Rank)
		if err != nil {
			return nil, fmt.Errorf("fetch record error: %w", err)
		}

		ret = append(ret, &th)
	}

	return ret, nil
}
