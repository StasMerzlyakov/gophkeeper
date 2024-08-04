package fs_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/StasMerzlyakov/gophkeeper/internal/server/adapters/bigfiles/fs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestFileStorage(t *testing.T) {

	tempDir := os.TempDir()

	storagePath := filepath.Join(tempDir, "temp-storage")
	defer func() {
		err := os.RemoveAll(storagePath)
		require.NoError(t, err)
	}()

	fs := fs.NewFileStorage(&config.ServerConf{
		FStoragePath: storagePath,
	})

	ctx := context.Background()
	bucket := uuid.New().String()

	lst, err := fs.GetFileInfoList(ctx, bucket)

	require.NoError(t, err)
	require.Equal(t, 0, len(lst))

	writer, err := fs.CreateStreamFileWriter(ctx, bucket)
	require.NoError(t, err)
	require.NotNil(t, writer)

	chunk1 := "hello "
	chunk2 := "world"
	fileName := "name"

	err = writer.WriteChunk(ctx, fileName, []byte(chunk1))
	require.NoError(t, err)

	err = writer.WriteChunk(ctx, fileName, []byte(chunk2))
	require.NoError(t, err)

	err = writer.WriteChunk(ctx, fileName+"!", []byte(chunk2))
	require.ErrorIs(t, err, domain.ErrServerInternal)

	lst, err = fs.GetFileInfoList(ctx, bucket)
	require.NoError(t, err)
	require.Equal(t, 0, len(lst))

	err = writer.Commit(ctx)
	require.NoError(t, err)

	lst, err = fs.GetFileInfoList(ctx, bucket)
	require.NoError(t, err)
	require.Equal(t, 1, len(lst))

	err = fs.DeleteFileInfo(ctx, bucket, fileName)
	require.NoError(t, err)

	err = fs.DeleteFileInfo(ctx, bucket+"?", fileName)
	require.ErrorIs(t, err, domain.ErrClientDataIncorrect)

}
