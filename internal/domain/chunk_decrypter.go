package domain

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"hash"
	"sync"

	"golang.org/x/crypto/pbkdf2"
)

func NewChunkDecrypter(password string) *chunkDecrypter {

	return &chunkDecrypter{
		pass:        password,
		tailStorage: NewTailStorage(sha256.BlockSize),
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
		if len(chunk) < aes.BlockSize+saltLen {
			panic("chunk too short") // it is possible to use buffer to store sufficient number of bytes
		}

		salt := chunk[:saltLen]
		fmt.Println("decr salt: " + base64.StdEncoding.EncodeToString(salt))

		// get the IV from the ciphertext
		iv := chunk[saltLen : aes.BlockSize+saltLen]
		fmt.Println("decr iv: " + base64.StdEncoding.EncodeToString(iv))

		key := pbkdf2.Key([]byte(che.pass), salt, iterations, keyLen, sha256.New)

		block, err := aes.NewCipher(key)
		if err != nil {
			panic(err)
		}

		che.cipher = cipher.NewCFBDecrypter(block, iv)
		che.hasher = hmac.New(sha256.New, key)
		che.pass = ""

		chunk = chunk[:aes.BlockSize+saltLen] // remove salt + iv
	})

	// decrypt chunk
	che.cipher.XORKeyStream(chunk, chunk)

	res := che.tailStorage.Write(chunk)
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
