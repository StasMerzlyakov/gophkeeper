package domain_test

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChunkEncryptDecryptTest(t *testing.T) {

	t.Run("test_1", func(t *testing.T) {
		masterPass := "masterPass"

		encryptor := domain.NewChunkEncrypter(masterPass)
		decryptor := domain.NewChunkDecrypter(masterPass)

		var actualBuf bytes.Buffer
		var expectedBuf bytes.Buffer
		var encryptedBuf bytes.Buffer

		str := "hello"

		_, err := expectedBuf.Write([]byte(str))
		require.NoError(t, err)

		encrypted, err := encryptor.WriteChunk([]byte(str))
		require.NoError(t, err)

		tailEncr, err := encryptor.Finish()
		require.NoError(t, err)

		decrChunk1, err := decryptor.WriteChunk(encrypted)
		require.NoError(t, err)

		_, err = encryptedBuf.Write(encrypted)
		require.NoError(t, err)

		_, err = encryptedBuf.Write(tailEncr)
		require.NoError(t, err)

		_, err = actualBuf.Write(decrChunk1)
		require.NoError(t, err)

		decrChunk2, err := decryptor.WriteChunk(tailEncr)
		require.NoError(t, err)

		_, err = actualBuf.Write(decrChunk2)
		require.NoError(t, err)

		require.NoError(t, decryptor.Finish())

		require.Equal(t, 0, bytes.Compare(expectedBuf.Bytes(), actualBuf.Bytes()))

		require.NotEqual(t, 0, bytes.Compare(expectedBuf.Bytes(), encryptedBuf.Bytes()))
		require.NotEqual(t, len(expectedBuf.Bytes()), len(encryptedBuf.Bytes()))

	})

	t.Run("test_2", func(t *testing.T) {
		masterPass := "masterPass"

		encryptor := domain.NewChunkEncrypter(masterPass)
		decryptor := domain.NewChunkDecrypter(masterPass)

		var actualBuf bytes.Buffer
		var expectedBuf bytes.Buffer

		str := "hello"

		_, err := expectedBuf.Write([]byte(str))
		require.NoError(t, err)

		mess := []byte(str)

		var tBuf bytes.Buffer
		encrypted, err := encryptor.WriteChunk(mess[:1])
		require.NoError(t, err)
		_, err = tBuf.Write(encrypted)
		require.NoError(t, err)

		encrypted, err = encryptor.WriteChunk(mess[1:2])
		require.NoError(t, err)
		_, err = tBuf.Write(encrypted)
		require.NoError(t, err)

		encrypted, err = encryptor.WriteChunk(mess[2:3])
		require.NoError(t, err)
		_, err = tBuf.Write(encrypted)
		require.NoError(t, err)

		encrypted, err = encryptor.WriteChunk(mess[3:4])
		require.NoError(t, err)
		_, err = tBuf.Write(encrypted)
		require.NoError(t, err)

		encrypted, err = encryptor.WriteChunk(mess[4:5])
		require.NoError(t, err)
		_, err = tBuf.Write(encrypted)
		require.NoError(t, err)

		tailEncr, err := encryptor.Finish()
		require.NoError(t, err)

		decrChunk1, err := decryptor.WriteChunk(tBuf.Bytes())
		require.NoError(t, err)

		_, err = actualBuf.Write(decrChunk1)
		require.NoError(t, err)

		decrChunk2, err := decryptor.WriteChunk(tailEncr)
		require.NoError(t, err)

		_, err = actualBuf.Write(decrChunk2)
		require.NoError(t, err)

		fmt.Println(str)
		fmt.Println(actualBuf.String())

		require.NoError(t, decryptor.Finish())

		require.Equal(t, 0, bytes.Compare(expectedBuf.Bytes(), actualBuf.Bytes()))

	})

	t.Run("test_2", func(t *testing.T) {
		masterPass := "masterPass"

		encryptor := domain.NewChunkEncrypter(masterPass)
		decryptor := domain.NewChunkDecrypter(masterPass)

		var actualBuf bytes.Buffer
		var expectedBuf bytes.Buffer

		chunk := make([]byte, 1024)
		for i := 0; i < 10; i++ {
			_, err := rand.Read(chunk)
			require.NoError(t, err)
			_, err = expectedBuf.Write(chunk)
			require.NoError(t, err)

			encrypted, err := encryptor.WriteChunk(chunk)
			require.NoError(t, err)

			decrypted, err := decryptor.WriteChunk(encrypted)
			require.NoError(t, err)

			_, err = actualBuf.Write(decrypted)
			require.NoError(t, err)
		}

		tailEncr, err := encryptor.Finish()
		require.NoError(t, err)

		tailDecr, err := decryptor.WriteChunk(tailEncr)
		require.NoError(t, err)

		_, err = actualBuf.Write(tailDecr)
		require.NoError(t, err)

		require.NoError(t, decryptor.Finish())

		expected := expectedBuf.Bytes()
		actual := actualBuf.Bytes()
		assert.Equal(t, 0, bytes.Compare(expected, actual))
	})

}
