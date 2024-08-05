package domain_test

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTailStorage(t *testing.T) {

	t.Run("long_chunk", func(t *testing.T) {

		testCase := []struct {
			name string
			size int
		}{
			{
				name: "len_0",
				size: 0,
			},
			{
				name: "len_32",
				size: 32,
			},
		}

		for _, tst := range testCase {
			t.Run(tst.name, func(t *testing.T) {
				storage := domain.NewTailStorage(tst.size)
				var actualBuf bytes.Buffer
				var expectedBuf bytes.Buffer

				chunk := make([]byte, 1024)
				for i := 0; i < 10; i++ {
					_, err := rand.Read(chunk)
					require.NoError(t, err)
					_, err = expectedBuf.Write(chunk)
					require.NoError(t, err)
					writen := storage.Write(chunk)
					_, err = actualBuf.Write(writen)
					require.NoError(t, err)
				}

				tail, err := storage.Finish()
				assert.NoError(t, err)
				assert.Equal(t, tst.size, len(tail))
				_, err = actualBuf.Write(tail)
				require.NoError(t, err)
				expected := expectedBuf.Bytes()
				actual := actualBuf.Bytes()
				assert.Equal(t, 0, bytes.Compare(expected, actual))
			})
		}
	})

	t.Run("shor_chunk", func(t *testing.T) {

		testCase := []struct {
			name string
			size int
		}{
			{
				name: "len_16",
				size: 16,
			},
			{
				name: "len_32",
				size: 32,
			},
		}

		for _, tst := range testCase {
			t.Run(tst.name, func(t *testing.T) {
				tailSize := 2 * tst.size
				storage := domain.NewTailStorage(tailSize)
				var actualBuf bytes.Buffer
				var expectedBuf bytes.Buffer

				chunk := make([]byte, tst.size)
				for i := 0; i < 10; i++ {
					_, err := rand.Read(chunk)
					require.NoError(t, err)
					_, err = expectedBuf.Write(chunk)
					require.NoError(t, err)
					writen := storage.Write(chunk)
					_, err = actualBuf.Write(writen)
					require.NoError(t, err)
				}

				tail, err := storage.Finish()
				assert.NoError(t, err)
				assert.Equal(t, tailSize, len(tail))
				_, err = actualBuf.Write(tail)
				require.NoError(t, err)
				expected := expectedBuf.Bytes()
				actual := actualBuf.Bytes()
				assert.Equal(t, 0, bytes.Compare(expected, actual))
			})
		}
	})

	t.Run("with_tail", func(t *testing.T) {

		testCase := []struct {
			name string
			size int
		}{
			{
				name: "len_16",
				size: 16,
			},
			{
				name: "len_32",
				size: 32,
			},
		}

		for _, tst := range testCase {
			t.Run(tst.name, func(t *testing.T) {
				tailSize := 2 * tst.size
				storage := domain.NewTailStorage(tailSize)
				var actualBuf bytes.Buffer
				var expectedBuf bytes.Buffer

				chunk := make([]byte, tailSize)
				for i := 0; i < 10; i++ {
					if _, err := rand.Read(chunk); err != nil {
						panic(err)
					}
					if _, err := expectedBuf.Write(chunk); err != nil {
						panic(err)
					}
					writen := storage.Write(chunk)
					if _, err := actualBuf.Write(writen); err != nil {
						panic(err)
					}
				}

				chunk = make([]byte, 10)
				if _, err := rand.Read(chunk); err != nil {
					panic(err)
				}

				if _, err := expectedBuf.Write(chunk); err != nil {
					panic(err)
				}
				writen := storage.Write(chunk)
				if _, err := actualBuf.Write(writen); err != nil {
					panic(err)
				}

				tail, err := storage.Finish()
				assert.NoError(t, err)
				assert.Equal(t, tailSize, len(tail))
				if _, err := actualBuf.Write(tail); err != nil {
					panic(err)
				}
				expected := expectedBuf.Bytes()
				actual := actualBuf.Bytes()
				assert.Equal(t, 0, bytes.Compare(expected, actual))
			})
		}
	})

}

func TestTailStorageHello(t *testing.T) {

	tail := make([]byte, 32)
	_, err := rand.Read(tail)
	require.NoError(t, err)

	hello := "hello"

	tailStorage := domain.NewTailStorage(32)
	res := tailStorage.Write([]byte(hello))
	assert.Nil(t, res)

	res = tailStorage.Write(tail)
	assert.Equal(t, hello, string(res))
}
