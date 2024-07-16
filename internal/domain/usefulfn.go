package domain

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/mail"
	"path/filepath"
	"runtime"

	pasVld "github.com/wagslane/go-password-validator"
)

func GetAction(depth int) string {
	pc, _, _, _ := runtime.Caller(depth)
	action := runtime.FuncForPC(pc).Name()
	return filepath.Base(action)
}

// Minimal account password complexity level
// https://github.com/wagslane/go-password-validator
const minAccountPassEntropyBits = 30

func CheckEMailData(data *EMailData) (bool, error) {
	action := GetAction(2)
	log := GetApplicationLogger()

	if _, err := mail.ParseAddress(data.EMail); err != nil {
		log.Warnf(action, "err", fmt.Sprintf("email %s parse error", data.EMail))
		return false, fmt.Errorf("%w - email is bad", ErrDataFormat)
	}

	if err := pasVld.Validate(data.Password, minAccountPassEntropyBits); err != nil {
		log.Warnf(action, "err", fmt.Sprintf("password for email %s is too simple", data.EMail))
		return false, fmt.Errorf("%w - pass is bad", ErrDataFormat)
	}

	return true, nil
}

type saltFn func(b []byte) (n int, err error)

type hashData struct {
	Hash string
	Salt string
}

func HashPassword(pass string, saltFn saltFn) (*hashData, error) {
	salt := make([]byte, 16)
	if _, err := saltFn(salt); err != nil {
		action := GetAction(2)
		log := GetApplicationLogger()
		log.Warnf(action, "err", fmt.Sprintf("salt generation error %s", err.Error()))
		return nil, fmt.Errorf("%w - salt generation", ErrInternalServer)
	}

	saltB64 := base64.URLEncoding.EncodeToString(salt)

	h := sha256.New()
	h.Write(salt)
	h.Write([]byte(pass))
	hex := h.Sum(nil)
	hexB64 := base64.URLEncoding.EncodeToString(hex)

	return &hashData{
		Hash: hexB64,
		Salt: saltB64,
	}, nil
}

// Minimal server secret key complexity level
// https://github.com/wagslane/go-password-validator
const minSecretKeyPassEntropyBits = 120

func CheckServerSecretKey(pass string) error {
	if len(pass) != 2*aes.BlockSize {
		return fmt.Errorf("wrong secret key length, expected %d", 2*aes.BlockSize)
	}

	if err := pasVld.Validate(pass, minSecretKeyPassEntropyBits); err != nil {
		return fmt.Errorf("secret key too simple %w", err)
	}

	return nil
}

func CheckPassword(pass string, hashB64 string, saltB64 string) (bool, error) {
	hash, err := base64.URLEncoding.DecodeString(hashB64)
	if err != nil {
		action := GetAction(2)
		log := GetApplicationLogger()
		log.Warnf(action, "err", fmt.Sprintf("hash decoding error %s", err.Error()))
		return false, fmt.Errorf("%w - hash decoding error", ErrDataFormat)
	}

	salt, err := base64.URLEncoding.DecodeString(saltB64)
	if err != nil {
		action := GetAction(2)
		log := GetApplicationLogger()
		log.Warnf(action, "err", fmt.Sprintf("salt decoding error %s", err.Error()))
		return false, fmt.Errorf("%w - salt decoding error", ErrDataFormat)
	}

	h := sha256.New()
	h.Write(salt)
	h.Write([]byte(pass))
	hex := h.Sum(nil)

	if !bytes.Equal(hash, hex) {
		return false, fmt.Errorf("%w - authentification failed", ErrWrongLoginPassword)
	}

	return true, nil
}

func EncryptData(secretKey string, plaintext string, saltFn saltFn) (string, error) {
	aes, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		action := GetAction(2)
		log := GetApplicationLogger()
		log.Warnf(action, "err", fmt.Sprintf("encrypt error %s", err.Error()))
		return "", fmt.Errorf("%w - cipher error", ErrInternalServer)
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		action := GetAction(2)
		log := GetApplicationLogger()
		log.Warnf(action, "err", fmt.Sprintf("encrypt error %s", err.Error()))
		return "", fmt.Errorf("%w - gcm error", ErrInternalServer)
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = saltFn(nonce)
	if err != nil {
		action := GetAction(2)
		log := GetApplicationLogger()
		log.Warnf(action, "err", fmt.Sprintf("salt generation error %s", err.Error()))
		return "", fmt.Errorf("%w - salt generation", ErrInternalServer)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	return string(ciphertext), nil
}

func DecryptData(secretKey string, ciphertext string) (string, error) {
	aes, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		action := GetAction(2)
		log := GetApplicationLogger()
		log.Warnf(action, "err", fmt.Sprintf("decrypt error %s", err.Error()))
		return "", fmt.Errorf("%w - cipher error", ErrInternalServer)
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		action := GetAction(2)
		log := GetApplicationLogger()
		log.Warnf(action, "err", fmt.Sprintf("decrypt error %s", err.Error()))
		return "", fmt.Errorf("%w - gcm error", ErrInternalServer)
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		action := GetAction(2)
		log := GetApplicationLogger()
		log.Warnf(action, "err", fmt.Sprintf("decrypt error %s", err.Error()))
		return "", fmt.Errorf("%w - gcm error", ErrInternalServer)
	}

	return string(plaintext), nil
}
