package domain_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestFileStorage(t *testing.T) {

	tempDir := os.TempDir()

	storagePath := filepath.Join(tempDir, "temp-storage")
	os.RemoveAll(storagePath)
	defer os.RemoveAll(storagePath)

	writer, err := domain.NewStreamFileWriter(storagePath)

	require.NoError(t, err)
	require.NotNil(t, writer)

	ctx := context.Background()

	chunk1 := "hello "
	chunk2 := "world"
	fileName := "name"

	err = writer.WriteChunk(ctx, fileName, []byte(chunk1))
	require.NoError(t, err)

	err = writer.WriteChunk(ctx, fileName, []byte(chunk2))
	require.NoError(t, err)

	err = writer.WriteChunk(ctx, fileName+"!", []byte(chunk2))
	require.ErrorIs(t, err, domain.ErrServerInternal)

	err = writer.Commit(ctx)
	require.NoError(t, err)

	writer2, err := domain.NewStreamFileWriter(storagePath)

	require.NoError(t, err)
	require.NotNil(t, writer2)

	err = writer2.WriteChunk(ctx, fileName, []byte(chunk1))
	require.ErrorIs(t, err, domain.ErrClientDataIncorrect)

}
