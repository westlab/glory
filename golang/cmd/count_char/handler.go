package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/westlab/glory/db"
)

func FetchAuthorID(name, wg string) (int64, error) {

	var id int64

	if err := db.DB.QueryRow("SELECT author_id FROM author WHERE name = ? and working_group = ? LIMIT 1", name, wg).Scan(&id); err != nil {
		return -1, err
	}

	return id, nil
}

func CreateAuthorID(name, wg string) (int64, error) {

	ins, err := db.DB.Prepare("INSERT INTO author VALUE(0, ?, ?)")
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

	ins, err := db.DB.Prepare("INSERT INTO thesis_history VALUE(0, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("sql preparation error: %w", err)
	}

	if _, err := ins.Exec(authorID, count, lastMod, fetchTime); err != nil {
		return fmt.Errorf("db insert error: %w", err)
	}
	return nil
}

func calcAuthorID(name, wg string) (int64, error) {
	id, err := FetchAuthorID(name, wg)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return -1, err
		}
		id, err = CreateAuthorID(name, wg)
		if err != nil {
			return -1, err
		}
	}
	return id, nil
}
