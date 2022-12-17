package countchars

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/westlab/glory"
)

func GetAllAuthor() ([]*glory.Author, error) {
	var ret []*glory.Author

	db, err := sql.Open("mysql", glory.DataSourceName)
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("SELECT * FROM author")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var author glory.Author
		err := rows.Scan(&author.Id, &author.Name, &author.WorkingGroup)
		if err != nil {
			return nil, err
		}
		ret = append(ret, &author)
	}
	return ret, nil
}

func GetAllThesisHistory() ([]*glory.ThesisHistory, error) {
	var ret []*glory.ThesisHistory

	db, err := sql.Open("mysql", glory.DataSourceName)
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("SELECT * FROM thesis_history")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var th glory.ThesisHistory
		var lm, ft string
		err := rows.Scan(&th.Id, &th.AuthorId, &th.CharCount, &lm, &ft)
		if err != nil {
			return nil, err
		}
		th.LastMod, err = time.Parse(glory.TimeFormat, lm)
		if err != nil {
			return nil, err
		}
		th.FetchTime, err = time.Parse(glory.TimeFormat, ft)
		if err != nil {
			return nil, err
		}
		th.LastMod = th.LastMod.Local()
		th.FetchTime = th.FetchTime.Local()
		ret = append(ret, &th)
	}
	return ret, nil
}

func FetchAuthorID(name, wg string) (int64, error) {
	var id int64
	db, err := sql.Open("mysql", glory.DataSourceName)
	if err != nil {
		return -1, fmt.Errorf("db connection error: %w", err)
	}
	if err = db.QueryRow("SELECT author_id FROM author WHERE name = ? and working_group = ? LIMIT 1", name, wg).Scan(&id); err != nil {
		return -1, fmt.Errorf("db read error: %w", err)
	}
	return id, nil
}

func CreateAuthorID(name, wg string) (int64, error) {
	db, err := sql.Open("mysql", glory.DataSourceName)
	if err != nil {
		return -1, fmt.Errorf("db connection error: %w", err)
	}
	ins, err := db.Prepare("INSERT INTO author VALUE(0, ?, ?)")
	if err != nil {
		return -1, fmt.Errorf("sql preparation error: %w", err)
	}
	result, err := ins.Exec(name, wg)
	if err != nil {
		return -1, fmt.Errorf("db insert error: %w", err)
	}
	id, _ := result.LastInsertId()
	return id, nil
}

func CreateThesisHistory(authorID, count int64, lastMod, fetchTime string) error {
	db, err := sql.Open("mysql", glory.DataSourceName)
	if err != nil {
		return fmt.Errorf("db connection error: %w", err)
	}
	ins, err := db.Prepare("INSERT INTO thesis_history VALUE(0, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("sql preparation error: %w", err)
	}
	if _, err := ins.Exec(authorID, count, lastMod, fetchTime); err != nil {
		return fmt.Errorf("db insert error: %w", err)
	}
	return nil
}
