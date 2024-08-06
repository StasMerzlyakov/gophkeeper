package fs_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"io"
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
	t.Run("test_steram_writer", func(t *testing.T) {
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
	})

	t.Run("test_steram_reader", func(t *testing.T) {

		tempDir := os.TempDir()

		bucket := uuid.New().String()

		storagePath := filepath.Join(tempDir, "temp-storage")
		defer func() {
			err := os.RemoveAll(storagePath)
			require.NoError(t, err)
		}()

		bucketPath := filepath.Join(storagePath, bucket)
		err := os.MkdirAll(bucketPath, os.ModePerm)
		require.NoError(t, err)

		f, err := os.CreateTemp(bucketPath, "sample")
		require.NoError(t, err)

		defer func() {
			err := os.Remove(f.Name())
			require.NoError(t, err)
		}()

		chunkSize := domain.FileChunkSize

		bufSize := chunkSize*2 + chunkSize/2
		buf := make([]byte, bufSize)

		_, err = rand.Read(buf)
		require.NoError(t, err)

		_, err = f.Write(buf)
		require.NoError(t, err)

		err = f.Close()
		require.NoError(t, err)

		fs := fs.NewFileStorage(&config.ServerConf{
			FStoragePath: storagePath,
		})

		name := filepath.Base(f.Name())

		ctx := context.Background()
		reader, err := fs.CreateStreamFileReader(ctx, bucket, name)

		require.NoError(t, err)
		require.NotNil(t, reader)

		require.Equal(t, int64(bufSize), reader.FileSize())

		var resultBuf bytes.Buffer

		chunk1, err := reader.Next()
		require.NoError(t, err)
		require.NotNil(t, chunk1)
		resultBuf.Write(chunk1)

		chunk2, err := reader.Next()
		require.NoError(t, err)
		require.NotNil(t, chunk2)
		resultBuf.Write(chunk2)

		chunk3, err := reader.Next()
		require.NoError(t, err)
		require.NotNil(t, chunk3)
		resultBuf.Write(chunk3)

		require.Equal(t, 0, bytes.Compare(buf, resultBuf.Bytes()))
		chunk4, err := reader.Next()
		require.ErrorIs(t, err, io.EOF)
		require.Equal(t, 0, len(chunk4))

		reader.Close()

	})
}
