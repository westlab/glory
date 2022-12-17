package glory

import "time"

const (
	TimeFormat = "2006-01-02 15:04:05"
	DateFormat = "2006-01-02"
)

type Author struct {
	Id           int
	Name         string
	WorkingGroup string
}

type ThesisHistory struct {
	Id        int64
	AuthorId  int64
	CharCount int64
	LastMod   time.Time
	FetchTime time.Time
}
