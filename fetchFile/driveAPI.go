package fetchFile

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/drive/v3"
	"io"
	"log"
	"os"
	"path"
	"time"
)

var ErrNotFound = errors.New("file not found")

// FetchFileID はファイル名（フォルダ名）からIDを取得する
func FetchFileID(ctx context.Context, srv *drive.Service, name string) (string, error) {
	r, err := srv.Files.List().Fields("nextPageToken, files(parents, id, name)").
		Q(fmt.Sprintf("name = '%s'", name)).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("read file error: %w", err)
	}
	if len(r.Files) == 0 {
		return "", ErrNotFound
	}

	return r.Files[0].Id, nil
}

// FetchLatestDocx は基準時刻以降で最新のdocxファイルを取得する
func FetchLatestDocx(ctx context.Context, srv *drive.Service, dirID string, threshold time.Time) (*drive.File, error) {
	r, err := srv.Files.List().Fields("nextPageToken, files(parents, id, name, modifiedTime, mimeType)").
		Q(fmt.Sprintf("'%s' in parents and name contains 'docx' and modifiedTime > '%s'", dirID, threshold.Format(time.RFC3339))).OrderBy("modifiedTime").Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("fetch latest docx attributes error: %w", err)
	}
	if len(r.Files) == 0 {
		return nil, ErrNotFound
	}
	fileName := r.Files[0].Name
	fileID := r.Files[0].Id
	resp, err := srv.Files.Get(fileID).Download()
	if err != nil {
		return nil, fmt.Errorf("download file error: %w", err)
	}
	output, err := os.Create(path.Join("workspace", fileName))
	if err != nil {
		return nil, fmt.Errorf("create file error: %w", err)
	}
	defer output.Close()
	n, err := io.Copy(output, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("write file error: %w", err)
	}
	log.Printf("%v: %v bytes downloaded", fileName, n)

	return r.Files[0], nil
}

// FetchSubDirectories は入力のIDがparentIDであるディレクトリを取得する
func FetchSubDirectories(ctx context.Context, srv *drive.Service, parentID string) ([]*drive.File, error) {
	r, err := srv.Files.List().Fields("nextPageToken, files(parents, id, name, owners, kind, mimeType)").
		Q(fmt.Sprintf("'%s' in parents and mimeType = 'application/vnd.google-apps.folder'", parentID)).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return r.Files, nil
}
