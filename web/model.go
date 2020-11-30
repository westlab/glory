package web

import (
	"github.com/Songmu/flextime"
	"github.com/westlab/glory"
	"time"
)

const (
	dateFormat  = "01/02"
	defaultMaxi = 100
	steps       = 10
)

type ThesisHistoryJoinAuthor struct {
	Name      string
	CharCount int
	FetchTime time.Time
}

type ProgressChart struct {
	Dates          []string
	Maxi           int
	StepSize       int
	AuthorDataList []*AuthorData
}

type AuthorData struct {
	Name  string
	Data  []int
	Color RGB
}

func NewProgressChart(wg *glory.WorkingGroup, th []*ThesisHistoryJoinAuthor) *ProgressChart {
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
		ret.AuthorDataList[authorIdx[t.Name]].Data[idx] = t.CharCount
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

type RGB struct {
	R int
	G int
	B int
}

var colors = []RGB{
	{241, 154, 56},
	{226, 67, 64},
	{181, 181, 173},
	{68, 153, 187},
	{84, 184, 137},
	{189, 165, 119},
	{139, 118, 208},
	{77, 169, 155},
	{147, 97, 58},
	{208, 78, 60},
	{46, 106, 177},
	{179, 193, 70},
	{182, 39, 93},
}
