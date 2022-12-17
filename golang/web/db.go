package main

import (
	"database/sql"
	"fmt"
	"github.com/Songmu/flextime"
	"github.com/jmoiron/sqlx"
	"github.com/westlab/glory/utils"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sqlx.DB

func InitializeDB(dsn string) (func() error, error) {
	var err error
	DB, err = sqlx.Open("mysql", dsn)
	if err != nil {
		return func() error { return nil }, err
	}
	return DB.Close, DB.Ping()
}
func IsNoRows(err error) bool {
	return err == sql.ErrNoRows
}

func FetchAllHistory(wg string) ([]*ThesisHistoryJoinAuthor, error) {
	var ret []*ThesisHistoryJoinAuthor
	//db, err := sql.Open("mysql", utils.DataSourceName)
	//if err != nil {
	//	return nil, fmt.Errorf("db connection error: %w", err)
	//}

	//rows, err := db.Query("SELECT name, char_count, fetch_time FROM thesis_history JOIN author ON author.author_id=thesis_history.author_id WHERE working_group=? and fetch_time <= ? ORDER BY fetch_time", wg, flextime.Now().UTC().Format(utils.TimeFormat))
	rows, _ := DB.Query("SELECT name, char_count, fetch_time FROM thesis_history JOIN author ON author.author_id=thesis_history.author_id WHERE working_group=? and fetch_time <= ? ORDER BY fetch_time", wg, flextime.Now().UTC().Format(utils.TimeFormat))
	if IsNoRows(rows.Err()) {
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
		th.FetchTime, err = time.Parse(utils.TimeFormat, ft)
		if err != nil {
			return nil, fmt.Errorf("convert to time error: %w", err)
		}
		th.FetchTime = th.FetchTime.Local()
		ret = append(ret, &th)
	}

	return ret, nil
}
