package fs

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

func NewFileStorage(conf *config.ServerConf) *fileStorage {
	return &fileStorage{
		path: conf.FStoragePath,
	}
}

type fileStorage struct {
	path string
}

func (fs *fileStorage) GetFileInfoList(ctx context.Context, bucket string) ([]domain.FileInfo, error) {
	dirPath := filepath.Join(fs.path, bucket)
	stat, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	if !stat.IsDir() {
		return nil, fmt.Errorf("%w bucker dir is not directory", domain.ErrServerInternal)
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	var result []domain.FileInfo
	for _, e := range entries {
		result = append(result, domain.FileInfo{
			Name: e.Name(),
		})
	}
	return result, nil
}

func (fs *fileStorage) DeleteFileInfo(ctx context.Context, bucket string, name string) error {
	filePath := filepath.Join(fs.path, bucket, name)
	return os.Remove(filePath)
}
func (fs *fileStorage) CreateStreamFileWriter(ctx context.Context, bucket string) (domain.StreamFileWriter, error) {
	dirPath := filepath.Join(fs.path, bucket)
	return domain.NewStreamFileWriter(dirPath)
}