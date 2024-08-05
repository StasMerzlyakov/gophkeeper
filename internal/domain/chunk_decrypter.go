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
	chn := make([]byte, len(chunk))
	copy(chn, chunk)

	che.readHeaderOnce.Do(func() {
		if len(chn) < aes.BlockSize+Pbkdf2SaltLen {
			err = errors.New("chunk too short") // it is possible to use buffer to store sufficient number of bytes
			return
		}

		salt := chn[:Pbkdf2SaltLen]

		fmt.Println("decr   salt " + base64.StdEncoding.EncodeToString(salt))

		// get the IV from the ciphertext
		iv := chn[Pbkdf2SaltLen : aes.BlockSize+Pbkdf2SaltLen]

		fmt.Println("decr   iv " + base64.StdEncoding.EncodeToString(iv))

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

		chn = chn[aes.BlockSize+Pbkdf2SaltLen:] // remove salt + iv
	})

	if err != nil {
		return nil, err
	}

	// decrypt chunk
	res := che.tailStorage.Write(chn)

	che.cipher.XORKeyStream(res, res)

	if _, err := che.hasher.Write(res); err != nil {
		return nil, fmt.Errorf("decrypt err - %w", err)
	}

	return res, nil
}

func (che *chunkDecrypter) Finish() error {
	mac, err := che.tailStorage.Finish()
	if err != nil {
		return fmt.Errorf("decrypt err - %w", err)
	}
	fmt.Println("decr   mac " + base64.StdEncoding.EncodeToString(mac))

	che.cipher.XORKeyStream(mac, mac) // encrypt mac

	fmt.Println("decryptmac " + base64.StdEncoding.EncodeToString(mac))

	extractedMac := che.hasher.Sum(nil)

	fmt.Println("extracted " + base64.StdEncoding.EncodeToString(extractedMac))

	if !hmac.Equal(extractedMac, mac) {
		return fmt.Errorf("decrypt err - hmac not equal")
	}

	return nil
}
