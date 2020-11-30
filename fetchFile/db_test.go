package fetchFile

import (
	"database/sql"
	"errors"
	"github.com/Songmu/flextime"
	"github.com/google/go-cmp/cmp"
	"github.com/westlab/glory"
	"testing"
	"time"
)

func TestGetAllAuthor(t *testing.T) {
	glory.SetupTest([]string{
		glory.ExecSQL("./testdata/teardown.sql"),
		glory.ExecSQL("./testdata/input.sql"),
	})
	defer glory.TearDown([]string{
		glory.ExecSQL("./testdata/teardown.sql"),
	})

	tests := []struct {
		name string
		want []glory.Author
	}{
		{
			name: "正常系　全てのauthorを取得できる",
			want: []glory.Author{
				{
					Id:           1,
					Name:         "suzuki",
					WorkingGroup: "B4",
				},
				{
					Id:           2,
					Name:         "takahashi",
					WorkingGroup: "B4",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := GetAllAuthor()
			if err != nil {
				t.Errorf("get authors error: %v", err)
			}

			if len(resp) != len(tt.want) {
				t.Errorf("length is wrong want %v, got %v", len(tt.want), len(resp))
				return
			}

			for i := range resp {
				if diff := cmp.Diff(*resp[i], tt.want[i]); diff != "" {
					t.Errorf("author %v differs: (-got +want)\n%v", i, diff)
				}
			}

		})
	}
}

func TestGetAllThesisHistory(t *testing.T) {
	glory.SetupTest([]string{
		glory.ExecSQL("./testdata/teardown.sql"),
		glory.ExecSQL("./testdata/input.sql"),
	})
	defer glory.TearDown([]string{
		glory.ExecSQL("./testdata/teardown.sql"),
	})

	tests := []struct {
		name string
		want []glory.ThesisHistory
	}{
		{
			name: "正常系　全てのauthorを取得できる",
			want: []glory.ThesisHistory{
				{
					Id:        1,
					AuthorId:  1,
					CharCount: 120,
					LastMod:   time.Date(2020, 11, 15, 11, 34, 0, 0, time.Local),
					FetchTime: time.Date(2020, 11, 15, 12, 0, 0, 0, time.Local),
				},
				{
					Id:        2,
					AuthorId:  2,
					CharCount: 140,
					LastMod:   time.Date(2020, 11, 15, 11, 44, 0, 0, time.Local),
					FetchTime: time.Date(2020, 11, 15, 12, 0, 0, 0, time.Local),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := GetAllThesisHistory()
			if err != nil {
				t.Errorf("get thesis history error: %v", err)
			}

			if len(resp) != len(tt.want) {
				t.Errorf("length is wrong want %v, got %v", len(tt.want), len(resp))
			}

			for i := range resp {
				if diff := cmp.Diff(*resp[i], tt.want[i]); diff != "" {
					t.Errorf("author %v differs: (-got +want)\n%v", i, diff)
				}
			}

		})
	}

}

func TestFetchAuthorID(t *testing.T) {
	glory.SetupTest([]string{
		glory.ExecSQL("./testdata/teardown.sql"),
		glory.ExecSQL("./testdata/input.sql"),
	})
	defer glory.TearDown([]string{
		glory.ExecSQL("./testdata/teardown.sql"),
	})

	tests := []struct {
		name         string
		author       string
		workingGroup string
		want         int64
		wantErr      *error
	}{
		{
			name:         "[正常系]　該当するauthorが存在する",
			author:       "suzuki",
			workingGroup: "B4",
			want:         1,
			wantErr:      nil,
		},
		{
			name:         "[正常系]　該当するauthorが存在しない",
			author:       "smith",
			workingGroup: "B4",
			want:         -1,
			wantErr:      &sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := FetchAuthorID(tt.author, tt.workingGroup)
			if err != nil {
				if tt.wantErr == nil {
					t.Error(err)
				}
				if errors.Is(*tt.wantErr, err) {
					t.Errorf("error of fetch author id: got %v, want %v", err, *tt.wantErr)
				}
			}

			if id != tt.want {
				t.Errorf("return id is wrong: got %d, want %d", id, tt.want)
			}
		})
	}
}

func TestCreateAuthorID(t *testing.T) {
	glory.SetupTest([]string{
		glory.ExecSQL("./testdata/teardown.sql"),
	})
	defer glory.TearDown([]string{
		glory.ExecSQL("./testdata/teardown.sql"),
	})

	tests := []struct {
		name         string
		author       string
		workingGroup string
		want         []glory.Author
	}{
		{
			name:         "[正常系] authorを追加できる",
			author:       "tanabe",
			workingGroup: "M2",
			want: []glory.Author{
				{
					Id:           1,
					Name:         "tanabe",
					WorkingGroup: "M2",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateAuthorID(tt.author, tt.workingGroup)
			if err != nil {
				t.Error(err)
			}
			resp, err := GetAllAuthor()
			for i := range resp {
				if diff := cmp.Diff(*resp[i], tt.want[i]); diff != "" {
					t.Errorf("author %v differs: (-got +want)\n%v", i, diff)
				}
			}
		})
	}
}

func TestCreateThesisHistory(t *testing.T) {
	glory.SetupTest([]string{
		glory.ExecSQL("./testdata/teardown.sql"),
		glory.ExecSQL("./testdata/input2.sql"),
	})
	defer glory.TearDown([]string{
		glory.ExecSQL("./testdata/teardown.sql"),
	})

	flextime.Fix(time.Date(2020, 4, 3, 12, 0, 0, 0, time.Local))

	tests := []struct {
		name     string
		authorID int64
		count    int64
		lastMod  time.Time
		want     []glory.ThesisHistory
	}{
		{
			name:     "thesis historyを追加できる",
			authorID: 1,
			count:    1200,
			lastMod:  time.Date(2020, 4, 3, 10, 0, 0, 0, time.Local),
			want: []glory.ThesisHistory{
				{
					Id:        1,
					AuthorId:  1,
					CharCount: 1200,
					LastMod:   time.Date(2020, 4, 3, 10, 0, 0, 0, time.Local),
					FetchTime: flextime.Now(),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CreateThesisHistory(tt.authorID, tt.count, tt.lastMod.UTC().Format(glory.TimeFormat), flextime.Now().UTC().Format(glory.TimeFormat))
			if err != nil {
				t.Errorf("create thesis history error: %v", err)
			}
			resp, err := GetAllThesisHistory()
			for i := range resp {
				if diff := cmp.Diff(*resp[i], tt.want[i]); diff != "" {
					t.Errorf("author %v differs: (-got +want)\n%v", i, diff)
				}
			}

		})
	}
}
