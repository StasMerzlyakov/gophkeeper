package domain

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
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
	var err error
	che.readHeaderOnce.Do(func() {
		if len(chunk) < aes.BlockSize+Pbkdf2SaltLen {
			err = errors.New("chunk too short") // it is possible to use buffer to store sufficient number of bytes
			return
		}

		salt := chunk[:Pbkdf2SaltLen]

		// get the IV from the ciphertext
		iv := chunk[Pbkdf2SaltLen : aes.BlockSize+Pbkdf2SaltLen]

		key := pbkdf2.Key([]byte(che.pass), salt, Pbkdf2Iter, EncryptAESKeyLen, sha256.New)

		var block cipher.Block
		block, err = aes.NewCipher(key)
		if err != nil {
			err = fmt.Errorf("init cipher err %w", err)
			return
		}

		che.cipher = cipher.NewCFBDecrypter(block, iv)

		che.hasher = hmac.New(sha256.New, key)

		che.pass = ""

		chunk = chunk[aes.BlockSize+Pbkdf2SaltLen:] // remove salt + iv
	})

	if err != nil {
		return nil, err
	}

	// decrypt chunk

	res := che.tailStorage.Write(chunk)

	chn := make([]byte, len(res))
	copy(chn, res)
	che.cipher.XORKeyStream(chn, chn)

	if _, err := che.hasher.Write(chn); err != nil {
		return nil, fmt.Errorf("decrypt err - %w", err)
	}

	return chn, nil
}

func (che *chunkDecrypter) Finish() error {
	mac, err := che.tailStorage.Finish()
	if err != nil {
		return fmt.Errorf("decrypt err - %w", err)
	}
	fmt.Println("decr   mac " + base64.StdEncoding.EncodeToString(mac))

	chn := make([]byte, len(mac))
	copy(chn, mac)
	che.cipher.XORKeyStream(chn, chn) // encrypt mac

	fmt.Println("decryptmac " + base64.StdEncoding.EncodeToString(chn))

	extractedMac := che.hasher.Sum(nil)

	fmt.Println("extracted " + base64.StdEncoding.EncodeToString(extractedMac))

	if !hmac.Equal(extractedMac, chn) {
		return fmt.Errorf("decrypt err - hmac not equal")
	}

	return nil
}
