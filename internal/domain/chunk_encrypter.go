package domain

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"hash"
	"io"
	"sync"

	"golang.org/x/crypto/pbkdf2"
)

func NewChunkEncrypterByReader(password string, reader io.Reader) (*chunkEncrypter, error) {
	// allocate memory to hold the header of the ciphertext
	header := make([]byte, Pbkdf2SaltLen+aes.BlockSize)

	// generate salt
	salt := header[:Pbkdf2SaltLen]
	if _, err := io.ReadFull(reader, salt); err != nil {
		return nil, fmt.Errorf("salt init err %w", err)
	}

	fmt.Println("encr   salt " + base64.StdEncoding.EncodeToString(salt))

	// generate initialization vector
	iv := header[Pbkdf2SaltLen : aes.BlockSize+Pbkdf2SaltLen]
	if _, err := io.ReadFull(reader, iv); err != nil {
		return nil, fmt.Errorf("iv init err %w", err)
	}
	fmt.Println("encr   iv " + base64.StdEncoding.EncodeToString(iv))

	decr := NewChunkDecrypter(password)

	// generate a 32 bit key with the provided password
	key := pbkdf2.Key([]byte(password), salt, Pbkdf2Iter, EncryptAESKeyLen, sha256.New)

	hasher := hmac.New(sha256.New, key)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher err %w", err)
	}

	cipher := cipher.NewCFBEncrypter(block, iv)

	return &chunkEncrypter{
		header:    header,
		hasher:    hasher,
		cipher:    cipher,
		decriptor: decr,
	}, nil
}

func NewChunkEncrypter(password string) (*chunkEncrypter, error) {
	return NewChunkEncrypterByReader(password, rand.Reader)
}

var _ ChunkEncrypter = (*chunkEncrypter)(nil)

type chunkEncrypter struct {
	header          []byte
	hasher          hash.Hash
	cipher          cipher.Stream
	writeHeaderOnce sync.Once
	decriptor       *chunkDecrypter
}

func (che *chunkEncrypter) WriteChunk(chunk []byte) ([]byte, error) {
	var res bytes.Buffer
	var err error
	che.writeHeaderOnce.Do(func() {
		if _, err = res.Write(che.header); err != nil {
			err = fmt.Errorf("write buf err %w", err)
		}
	})
	if err != nil {
		return nil, err
	}

	chn := make([]byte, len(chunk))
	copy(chn, chunk)

	// update hash
	if _, err := che.hasher.Write(chn); err != nil {
		return nil, fmt.Errorf("%w can't write chunk", err)
	}

	che.cipher.XORKeyStream(chn, chn)

	if _, err := res.Write(chn); err != nil {
		return nil, fmt.Errorf("%w can't write chunk", err)
	}

	result := res.Bytes()

	if _, err := che.decriptor.WriteChunk(result); err != nil {
		return nil, fmt.Errorf("%w can't write chunk", err)
	}

	return result, nil
}

func (che *chunkEncrypter) Finish() ([]byte, error) {
	mac := che.hasher.Sum(nil)
	// encrypt hmac

	fmt.Println("encr   mac " + base64.StdEncoding.EncodeToString(mac))

	che.cipher.XORKeyStream(mac, mac) // encrypt mac

	fmt.Println("encriptmac " + base64.StdEncoding.EncodeToString(mac))

	if _, err := che.decriptor.WriteChunk(mac); err != nil {
		return nil, fmt.Errorf("%w can't write chunk", err)
	}

	if err := che.decriptor.Finish(); err != nil {
		return nil, fmt.Errorf("%w finish err", err)
	}

	return mac, nil
}
