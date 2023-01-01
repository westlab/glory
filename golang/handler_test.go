package main

import (
	"fmt"
	"log"
	"testing"

	"github.com/westlab/glory/db"
)

func TestFetchRanking(t *testing.T) {

	done, err := db.InitializeDB(dsn)
	if err != nil {
		log.Fatalf("failed to initialize db: %v", err)
	}
	defer done()

	ranking, err := FetchRanking("B4")
	if err != nil {
		t.Errorf("fetchRanking failed: %v", err)
	}

	for _, row := range ranking {
		fmt.Println(row.Rank, row.Name, row.CharCount)
	}

}
