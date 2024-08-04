package domain

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"hash"
	"sync"

	"golang.org/x/crypto/pbkdf2"
)

func NewChunkDecrypter(password string) *chunkDecrypter {

	return &chunkDecrypter{
		pass:        password,
		tailStorage: NewTailStorage(sha256.Size),
	}
}

var _ ChunkDecrypter = (*chunkDecrypter)(nil)

type chunkDecrypter struct {
	pass           string
	hasher         hash.Hash
	cipher         cipher.Stream
	readHeaderOnce sync.Once
	tailStorage    *tailStorage
}

func (che *chunkDecrypter) WriteChunk(chunk []byte) ([]byte, error) {
	che.readHeaderOnce.Do(func() {
		if len(chunk) < aes.BlockSize+pbkdf2SaltLen {
			panic("chunk too short") // it is possible to use buffer to store sufficient number of bytes
		}

		salt := chunk[:pbkdf2SaltLen]

		// get the IV from the ciphertext
		iv := chunk[pbkdf2SaltLen : aes.BlockSize+pbkdf2SaltLen]

		key := pbkdf2.Key([]byte(che.pass), salt, pbkdf2Iter, keyLen, sha256.New)

		block, err := aes.NewCipher(key)
		if err != nil {
			panic(err)
		}

		che.cipher = cipher.NewCFBDecrypter(block, iv)
		che.hasher = hmac.New(sha256.New, key)
		che.pass = ""

		chunk = chunk[aes.BlockSize+pbkdf2SaltLen:] // remove salt + iv
	})

	// decrypt chunk
	che.cipher.XORKeyStream(chunk, chunk)

	res := che.tailStorage.Write(chunk)
	if _, err := che.hasher.Write(res); err != nil {
		return nil, fmt.Errorf("decrypt err - %w", err)
	}
	return res, nil
}

func (che *chunkDecrypter) Finish() error {
	expectedMac, err := che.tailStorage.Finish()
	if err != nil {
		return fmt.Errorf("decrypt err - %w", err)
	}

	extractedMac := che.hasher.Sum(nil)

	if !hmac.Equal(extractedMac, expectedMac) {
		return fmt.Errorf("decrypt err - hmac not equal")
	}

	return nil
}
