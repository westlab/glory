package web

import (
	"github.com/Songmu/flextime"
	"github.com/gin-gonic/gin"
	"github.com/westlab/glory"
	"log"
	"net/http"
	"strconv"
	"time"
)

var config *glory.Conf

func init() {
	var err error
	if config, err = glory.LoadConfig("../config.json"); err != nil {
		log.Fatalf("load config error: %v", err)
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
	for _, wg := range config.WorkingGroups {
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
	for _, wg := range config.WorkingGroups {
		ret = append(ret, wg.Title)
	}
	return &ret
}

func calcDeadlines() *map[string]string {
	ret := map[string]string{}
	for _, wg := range config.WorkingGroups {
		var left string
		deadline, err := time.Parse(glory.DateFormat, wg.Deadline)
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
