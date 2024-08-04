package domain_test

import (
	"bytes"
	"crypto/rand"
	"io"
	"os"
	"testing"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestStreamFileReader(t *testing.T) {

	f, err := os.CreateTemp("", "sample")
	require.NoError(t, err)

	defer func() {
		err := os.Remove(f.Name())
		require.NoError(t, err)
	}()

	chunkSize := 512

	bufSize := chunkSize*2 + chunkSize/2
	buf := make([]byte, bufSize)

	_, err = rand.Read(buf)
	require.NoError(t, err)

	_, err = f.Write(buf)
	require.NoError(t, err)

	err = f.Close()
	require.NoError(t, err)

	reader, err := domain.NewStreamFileReader(f.Name(), chunkSize)
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

}
