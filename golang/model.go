package main

import (
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

type RankingRow struct {
	Rank      int
	Name      string
	CharCount int
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
