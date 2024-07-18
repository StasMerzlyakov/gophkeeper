package domain

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"image/png"
	"net/mail"
	"path/filepath"
	"runtime"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	pasVld "github.com/wagslane/go-password-validator"
)

const (
	TOTPPeriod    = 30
	TOTPDigits    = otp.DigitsSix
	TOTPAlgorithm = otp.AlgorithmSHA512
	HelloWorld    = "hello from GophKeeper!!!"
	SaltSize      = 16
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
	if _, err := mail.ParseAddress(data.EMail); err != nil {
		return false, fmt.Errorf("%w - email %s parse error", ErrClientDataIncorrect, data.EMail)
	}

	if err := pasVld.Validate(data.Password, minAccountPassEntropyBits); err != nil {
		return false, fmt.Errorf("%w - password is too simple", ErrClientDataIncorrect)
	}

	return true, nil
}

type SaltFn func(b []byte) (n int, err error)

func HashPassword(pass string, saltFn SaltFn) (*HashData, error) {
	salt := make([]byte, SaltSize)
	if _, err := saltFn(salt); err != nil {
		return nil, fmt.Errorf("%w - salt generation error %s", ErrInternalServer, err.Error())
	}

	saltB64 := base64.URLEncoding.EncodeToString(salt)

	h := sha256.New()
	h.Write(salt)
	h.Write([]byte(pass))
	hex := h.Sum(nil)
	hexB64 := base64.URLEncoding.EncodeToString(hex)

	return &HashData{
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
		return false, fmt.Errorf("%w - hash decoding error %s", ErrClientDataIncorrect, err.Error())
	}

	salt, err := base64.URLEncoding.DecodeString(saltB64)
	if err != nil {
		return false, fmt.Errorf("%w - salt decoding error %s", ErrClientDataIncorrect, err.Error())
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

func EncryptData(secretKey string, plaintext string, saltFn SaltFn) (string, error) {
	aes, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("%w - encrypt NewCipher error %s", ErrInternalServer, err.Error())
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return "", fmt.Errorf("%w - encrypt NewGCM error %s", ErrInternalServer, err.Error())
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = saltFn(nonce)
	if err != nil {
		return "", fmt.Errorf("%w - encrypt salt generation %s", ErrInternalServer, err.Error())
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	return string(ciphertext), nil
}

func DecryptData(secretKey string, ciphertext string) (string, error) {
	aes, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("%w - decrypt newCipher error %s", ErrInternalServer, err.Error())
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return "", fmt.Errorf("%w - decrypt NewGCM error %s", ErrInternalServer, err.Error())
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return "", fmt.Errorf("%w - decrypt gcm open err %s", ErrInternalServer, err.Error())
	}

	return string(plaintext), nil
}

func GenerateQR(issuer string, accountName string) (string, []byte, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: accountName,
		Period:      TOTPPeriod,
		Digits:      TOTPDigits,
		Algorithm:   TOTPAlgorithm,
	})

	if err != nil {
		return "", nil, fmt.Errorf("%w GenerateQR TOTP gererate err - %s", ErrInternalServer, err.Error())
	}

	var buf bytes.Buffer
	img, err := key.Image(450, 450)
	if err != nil {
		return "", nil, fmt.Errorf("%w GenerateQR TOTP image generation err - %s", ErrInternalServer, err.Error())
	}
	if err = png.Encode(&buf, img); err != nil {
		return "", nil, fmt.Errorf("%w GenerateQR png.Encode err - %s", ErrInternalServer, err.Error())
	}

	keyURL := key.URL()
	return keyURL, buf.Bytes(), nil
}

func ValidatePassCode(keyURL string, passcode string) (bool, error) {
	key, err := otp.NewKeyFromURL(keyURL)
	if err != nil {
		return false, fmt.Errorf("%w ValidatePassCode restore key err - %s", err, err.Error())
	}

	validOpts := totp.ValidateOpts{
		Period:    TOTPPeriod,
		Digits:    TOTPDigits,
		Algorithm: TOTPAlgorithm,
	}

	valid, err := totp.ValidateCustom(passcode, key.Secret(), time.Now(), validOpts)
	return valid, err
}

func GenerateHello(saltFn SaltFn) (string, error) {

	salt := make([]byte, SaltSize)
	_, err := saltFn(salt)
	if err != nil {
		return "", fmt.Errorf("%w GenerateHelloWorld generate salt err - %s", ErrInternalServer, err.Error())
	}

	var bytesToGen bytes.Buffer
	bytesToGen.Write(salt)
	bytesToGen.Write([]byte(HelloWorld))
	return base64.RawURLEncoding.EncodeToString(bytesToGen.Bytes()), nil

}

func CheckHello(toCheck string) (bool, error) {
	bytes, err := base64.RawStdEncoding.DecodeString(toCheck)
	if err != nil {
		return false, fmt.Errorf("%w CheckHelloWorld decode err %s", ErrInternalServer, err.Error())
	}

	if len(bytes) <= SaltSize {
		return false, fmt.Errorf("%w CheckHelloWorld decode err - unexpected input size", ErrInternalServer)
	}

	bytes = bytes[SaltSize:]
	return string(bytes) == HelloWorld, nil
}
