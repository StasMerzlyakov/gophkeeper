package domain_test

import (
	"context"
	"crypto/aes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/stretchr/testify/require"
)

const TestDataDirectory = "../../testdata"

func TestBigFileEncryptDecryptChunkReader(t *testing.T) {

	masterPass := "masterPass"

	encryptor, err := domain.NewChunkEncrypter(masterPass)
	require.NoError(t, err)
	f, err := os.CreateTemp("", "bif-file-encrypted-")
	require.NoError(t, err)

	defer func() {
		require.NoError(t, os.RemoveAll(f.Name()))
	}()

	// encrypt big file
	encrChunkSize := 4 * domain.KiB

	chunkReader, err := domain.NewStreamFileReader(filepath.Join(TestDataDirectory, "bigfile.mp4"), encrChunkSize)

	require.NoError(t, err)

	for {

		chunk, err := chunkReader.Next()
		if err == nil {
			encryChun, err := encryptor.WriteChunk(chunk)
			require.NoError(t, err)

			_, err = f.Write(encryChun)
			require.NoError(t, err)

			_, err = f.Write([]byte{}) // test writes
			require.NoError(t, err)

			_, err = f.Write(nil) // test writes
			require.NoError(t, err)
		} else {
			if err == io.EOF {
				chunkReader.Close()
				require.Equal(t, 0, len(chunk))

				encrTail, err := encryptor.Finish()
				require.NoError(t, err)

				_, err = f.Write(encrTail)
				require.NoError(t, err)

				err = f.Close()
				require.NoError(t, err)
				break
			}
			require.NoError(t, err) // unexpected err
		}
		require.NotEqual(t, 0, len(chunk)) // test impl err
	}

	testCase := []struct {
		name          string
		decrChunkSize int
	}{
		{
			"block_size_min_length",
			aes.BlockSize + domain.Pbkdf2SaltLen,
		},
		{
			"block_size_4_kb",
			4 * domain.KiB,
		},
		{
			"block_size_4_aes_key",
			4 * aes.BlockSize,
		},
	}

	for _, test := range testCase {
		t.Run(test.name, func(t *testing.T) {

			decryptor := domain.NewChunkDecrypter(masterPass)
			chunkReader, err := domain.NewStreamFileReader(f.Name(), test.decrChunkSize)
			require.NoError(t, err)

			for {
				chunk, err := chunkReader.Next()
				if err == nil {
					require.NotEqual(t, 0, len(chunk))
					_, err := decryptor.WriteChunk(chunk)
					require.NoError(t, err)
				} else {
					if err == io.EOF {
						require.Equal(t, 0, len(chunk))
						chunkReader.Close()
						err = decryptor.Finish()
						require.NoError(t, err)
						break
					}
					require.NoError(t, err)
				}

			}

		})
	}

}

func TestBigFileEncryptDecryptChunkWriter(t *testing.T) {

	masterPass := "masterPass"

	encryptor, err := domain.NewChunkEncrypter(masterPass)
	require.NoError(t, err)
	tempFile, err := os.CreateTemp("", "bif-file-encrypted-")
	require.NoError(t, err)
	err = tempFile.Close()
	require.NoError(t, err)

	require.NoError(t, os.Remove(tempFile.Name())) // only name need

	defer func() {
		require.NoError(t, os.RemoveAll(tempFile.Name()))
	}()

	// encrypt big file
	encrChunkSize := 4 * domain.KiB

	testFile, err := os.Open(filepath.Join(TestDataDirectory, "bigfile.mp4"))
	require.NoError(t, err)
	defer func() {
		_ = testFile.Close()
	}()

	chunk := make([]byte, encrChunkSize)

	tempFileWriter, err := domain.CreateStreamFileWriter(filepath.Dir(tempFile.Name()))
	require.NoError(t, err)

	ctx := context.Background()
	tempFileName := filepath.Base(tempFile.Name())

	for {
		n, err := testFile.Read(chunk)
		if err == nil {
			encryChun, err := encryptor.WriteChunk(chunk[:n])
			require.NoError(t, err)

			err = tempFileWriter.WriteChunk(ctx, tempFileName, encryChun)
			require.NoError(t, err)

			err = tempFileWriter.WriteChunk(ctx, tempFileName, []byte{}) // test writes
			require.NoError(t, err)

			err = tempFileWriter.WriteChunk(ctx, tempFileName, nil) // test writes
			require.NoError(t, err)
		} else {
			if err == io.EOF {
				require.NoError(t, testFile.Close())
				require.Equal(t, 0, n)

				encrTail, err := encryptor.Finish()
				require.NoError(t, err)

				err = tempFileWriter.WriteChunk(ctx, tempFileName, encrTail)
				require.NoError(t, err)

				err = tempFileWriter.Commit(ctx)
				require.NoError(t, err)
				break
			}
			require.NoError(t, err) // unexpected err
		}
		require.NotEqual(t, 0, len(chunk)) // test impl err
	}

	testCase := []struct {
		name          string
		decrChunkSize int
	}{

		{
			"block_size_min_length",
			aes.BlockSize + domain.Pbkdf2SaltLen,
		},
		{
			"block_size_4_kb",
			4 * domain.KiB,
		},
		{
			"block_size_4_aes_key",
			4 * aes.BlockSize,
		},
	}

	for _, test := range testCase {
		t.Run(test.name, func(t *testing.T) {

			decryptor := domain.NewChunkDecrypter(masterPass)
			chunkReader, err := domain.NewStreamFileReader(tempFile.Name(), test.decrChunkSize)
			require.NoError(t, err)

			for {
				chunk, err := chunkReader.Next()
				if err == nil {
					require.NotEqual(t, 0, len(chunk))
					_, err := decryptor.WriteChunk(chunk)
					require.NoError(t, err)
				} else {
					if err == io.EOF {
						require.Equal(t, 0, len(chunk))
						chunkReader.Close()
						err = decryptor.Finish()
						require.NoError(t, err)
						break
					}
					require.NoError(t, err)
				}
			}
		})
	}
}
