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

func NewChunkEncrypterByReader(password string, reader io.Reader) *chunkEncrypter {
	// allocate memory to hold the header of the ciphertext
	header := make([]byte, Pbkdf2SaltLen+aes.BlockSize)

	// generate salt
	salt := header[:Pbkdf2SaltLen]
	if _, err := io.ReadFull(reader, salt); err != nil {
		panic(err)
	}
	// generate initialization vector
	iv := header[Pbkdf2SaltLen : aes.BlockSize+Pbkdf2SaltLen]
	if _, err := io.ReadFull(reader, iv); err != nil {
		panic(err)
	}

	// generate a 32 bit key with the provided password
	key := pbkdf2.Key([]byte(password), salt, Pbkdf2Iter, EncryptAESKeyLen, sha256.New)

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

func NewChunkEncrypter(password string) *chunkEncrypter {
	return NewChunkEncrypterByReader(password, rand.Reader)
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

	che.cipher.XORKeyStream(chunk, chunk)

	if _, err := res.Write(chunk); err != nil {
		panic(err)
	}

	result := res.Bytes()

	return result, nil
}

func (che *chunkEncrypter) Finish() ([]byte, error) {
	mac := che.hasher.Sum(nil)
	// encrypt hmac

	che.cipher.XORKeyStream(mac, mac) // encrypt mac

	return mac, nil
}
