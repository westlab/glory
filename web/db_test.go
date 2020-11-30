package web

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/westlab/glory"
	"testing"
	"time"
)

func TestFetchAllHistory(t *testing.T) {
	glory.SetupTest([]string{
		glory.ExecSQL("./testdata/teardown.sql"),
		glory.ExecSQL("./testdata/input.sql"),
	})
	defer glory.TearDown([]string{
		glory.ExecSQL("./testdata/teardown.sql"),
	})

	tests := []struct {
		name         string
		workingGroup string
		want         []ThesisHistoryJoinAuthor
		wantErr      error
	}{
		{
			name:         "履歴を取得できる",
			workingGroup: "B4",
			want: []ThesisHistoryJoinAuthor{
				{
					Name:      "suzuki",
					CharCount: 120,
					FetchTime: time.Date(2020, 11, 15, 12, 0, 0, 0, time.Local),
				},
				{
					Name:      "takahashi",
					CharCount: 140,
					FetchTime: time.Date(2020, 11, 15, 12, 0, 0, 0, time.Local),
				},
				{
					Name:      "suzuki",
					CharCount: 140,
					FetchTime: time.Date(2020, 11, 16, 12, 0, 0, 0, time.Local),
				},
				{
					Name:      "suzuki",
					CharCount: 220,
					FetchTime: time.Date(2020, 11, 17, 12, 0, 0, 0, time.Local),
				},
				{
					Name:      "takahashi",
					CharCount: 160,
					FetchTime: time.Date(2020, 11, 17, 12, 0, 0, 0, time.Local),
				},
				{
					Name:      "takahashi",
					CharCount: 360,
					FetchTime: time.Date(2020, 11, 19, 12, 0, 0, 0, time.Local),
				},
			},
			wantErr: nil,
		},
		{
			name:         "履歴が存在しない",
			workingGroup: "D3",
			want:         []ThesisHistoryJoinAuthor{},
			wantErr:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := FetchAllHistory(tt.workingGroup)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error is wrong: want %v, got %v", tt.wantErr, err)
				return
			}
			if resp == nil {
				return
			}

			if len(tt.want) != len(resp) {
				t.Errorf("length of resp is wrong:want %d, got %d", len(tt.want), len(resp))
			}

			for i := range resp {
				if diff := cmp.Diff(*resp[i], tt.want[i]); diff != "" {
					t.Errorf("author %v differs: (-got +want)\n%v", i, diff)
				}
			}
		})
	}

}
