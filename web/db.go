package web

import (
	"database/sql"
	"fmt"
	"github.com/Songmu/flextime"
	_ "github.com/go-sql-driver/mysql"
	"github.com/westlab/glory"
	"time"
)

func FetchAllHistory(wg string) ([]*ThesisHistoryJoinAuthor, error) {
	var ret []*ThesisHistoryJoinAuthor
	db, err := sql.Open("mysql", glory.DataSourceName)
	if err != nil {
		return nil, fmt.Errorf("db connection error: %w", err)
	}

	rows, err := db.Query("SELECT name, char_count, fetch_time FROM thesis_history JOIN author ON author.author_id=thesis_history.author_id WHERE working_group=? and fetch_time <= ? ORDER BY fetch_time", wg, flextime.Now().UTC().Format(glory.TimeFormat))

	if err != nil {
		return nil, fmt.Errorf("db read error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var th ThesisHistoryJoinAuthor
		var ft string
		err := rows.Scan(&th.Name, &th.CharCount, &ft)
		if err != nil {
			return nil, fmt.Errorf("fetch record error: %w", err)
		}
		th.FetchTime, err = time.Parse(glory.TimeFormat, ft)
		if err != nil {
			return nil, fmt.Errorf("convert to time error: %w", err)
		}
		th.FetchTime = th.FetchTime.Local()
		ret = append(ret, &th)
	}

	return ret, nil
}
