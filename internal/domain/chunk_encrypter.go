package domain

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"hash"
	"io"
	"sync"

	"golang.org/x/crypto/pbkdf2"
)

func NewChunkEncrypter(password string) *chunkEncrypter {

	// allocate memory to hold the header of the ciphertext
	header := make([]byte, pbkdf2SaltLen+aes.BlockSize)

	// generate salt
	salt := header[:pbkdf2SaltLen]
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		panic(err)
	}

	// generate initialization vector
	iv := header[pbkdf2SaltLen : aes.BlockSize+pbkdf2SaltLen]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	// generate a 32 bit key with the provided password
	key := pbkdf2.Key([]byte(password), salt, pbkdf2Iter, keyLen, sha256.New)

	hasher := hmac.New(sha256.New, key)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	cipher := cipher.NewCFBEncrypter(block, iv)

	return &chunkEncrypter{
		header: header,
		hasher: hasher,
		cipher: cipher,
	}
}

var _ ChunkEncrypter = (*chunkEncrypter)(nil)

type chunkEncrypter struct {
	header          []byte
	hasher          hash.Hash
	cipher          cipher.Stream
	writeHeaderOnce sync.Once
}

func (che *chunkEncrypter) WriteChunk(chunk []byte) ([]byte, error) {
	var res bytes.Buffer
	che.writeHeaderOnce.Do(func() {
		if _, err := res.Write(che.header); err != nil {
			panic(err)
		}
	})

	// update hash
	if _, err := che.hasher.Write(chunk); err != nil {
		panic(err)
	}

	encrChank := make([]byte, len(chunk))

	che.cipher.XORKeyStream(encrChank, chunk)

	if _, err := res.Write(encrChank); err != nil {
		panic(err)
	}

	return res.Bytes(), nil
}

func (che *chunkEncrypter) Finish() ([]byte, error) {
	hmac := che.hasher.Sum(nil)
	// encrypt hmac
	che.cipher.XORKeyStream(hmac, hmac)
	return hmac, nil
}
