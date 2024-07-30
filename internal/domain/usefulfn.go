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
	"image/png"
	"io"
	"net/mail"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain/creditcard"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	pasVld "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/pbkdf2"
)

const (
	AuthorizationMetadataTokenName = "authorization"
)

const (
	TOTPPeriod    = 30
	TOTPDigits    = otp.DigitsSix
	TOTPAlgorithm = otp.AlgorithmSHA512
	SaltSize      = 16
)

const (
	saltLen    = 32
	keyLen     = 32
	iterations = 100002
)

func GetAction(depth int) string {
	pc, _, _, _ := runtime.Caller(depth)
	action := runtime.FuncForPC(pc).Name()
	return filepath.Base(action)
}

// Minimal account password complexity level
// https://github.com/wagslane/go-password-validator
const minAccountPassEntropyBits = 30

// ParseEMail check email string format
func ParseEMail(email string) bool {
	if _, err := mail.ParseAddress(email); err != nil {
		return false
	}
	return true
}

// CheckEMailData check email string format and password entropy
func CheckEMailData(data *EMailData) (bool, error) {
	if !ParseEMail(data.EMail) {
		return false, fmt.Errorf("%w - email %s parse error", ErrClientDataIncorrect, data.EMail)
	}

	if err := pasVld.Validate(data.Password, minAccountPassEntropyBits); err != nil {
		return false, fmt.Errorf("%w - password is too simple", ErrClientDataIncorrect)
	}

	return true, nil
}

// CheckAuthPasswordComplexityLevel used to check auth pass
func CheckAuthPasswordComplexityLevel(pass string) bool {
	if err := pasVld.Validate(pass, minAccountPassEntropyBits); err != nil {
		return false
	}
	return true
}

// Minimal client data encryption password complexity level
// https://github.com/wagslane/go-password-validator
const minMasterPasswordEntropyBits = 80

func CheckMasterPasswordComplexityLevel(pass string) bool {
	if err := pasVld.Validate(pass, minMasterPasswordEntropyBits); err != nil {
		return false
	}
	return true
}

func Random32ByteString() string {
	//needs a randomly generated 32 character string. Exactly 32 characters. The string is 22 characters, but it's encoded to 32.
	b := make([]byte, 22)
	_, err := rand.Read(b)

	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}

// EncryptShortData encrypt data on password
// https://www.codemio.com/2023/05/advanced-golang-tutorials-aes-256.html
func EncryptShortData(masterPass string, data string) (string, error) {
	return encryptData(masterPass, data)
}

// DecryptShortData data on password
func DecryptShortData(masterKey string, ciphertext string) (string, error) {
	return decryptData(masterKey, ciphertext)
}

type SaltFn func(b []byte) (n int, err error)

// HashPassword has user password
func HashPassword(pass string, saltFn SaltFn) (*HashData, error) {
	salt := make([]byte, SaltSize)
	if _, err := saltFn(salt); err != nil {
		return nil, fmt.Errorf("%w - salt generation error %s", ErrServerInternal, err.Error())
	}

	saltB64 := base64.StdEncoding.EncodeToString(salt)

	h := sha256.New()
	h.Write(salt)
	h.Write([]byte(pass))
	hex := h.Sum(nil)
	hexB64 := base64.StdEncoding.EncodeToString(hex)

	return &HashData{
		Hash: hexB64,
		Salt: saltB64,
	}, nil
}

// Minimal server secret key complexity level
// https://github.com/wagslane/go-password-validator
const minSecretKeyPassEntropyBits = 120

// CheckServerSecretKeyComplexityLevel check server secret key complexity level
func CheckServerSecretKeyComplexityLevel(pass string) error {

	if err := pasVld.Validate(pass, minSecretKeyPassEntropyBits); err != nil {
		return fmt.Errorf("secret key too simple %w", err)
	}

	return nil
}

// ValidateAccountPass check server secret key complexity level
func ValidateAccountPass(pass string, hashB64 string, saltB64 string) (bool, error) {
	hash, err := base64.StdEncoding.DecodeString(hashB64)
	if err != nil {
		return false, fmt.Errorf("%w - hash decoding error %s", ErrClientDataIncorrect, err.Error())
	}

	salt, err := base64.StdEncoding.DecodeString(saltB64)
	if err != nil {
		return false, fmt.Errorf("%w - salt decoding error %s", ErrClientDataIncorrect, err.Error())
	}

	h := sha256.New()
	h.Write(salt)
	h.Write([]byte(pass))
	hex := h.Sum(nil)

	if !bytes.Equal(hash, hex) {
		return false, fmt.Errorf("%w - authentification failed", ErrAuthDataIncorrect)
	}

	return true, nil
}

// EncryptHello is used in registration process
func EncryptHello(masterPass string, hello string) (string, error) {
	return encryptData(masterPass, hello)
}

// DecryptHello is used to check masterPass is true
func DecryptHello(masterPass string, encryptedHello string) error {
	if _, err := decryptData(masterPass, encryptedHello); err != nil {
		return fmt.Errorf("%w - hello is not restored", err)
	}
	return nil
}

// EncryptOTPKey is used to secure qr code on server side
func EncryptOTPKey(secretKey string, otpKey string) (string, error) {
	return encryptData(secretKey, otpKey)
}

// DecryptOTPKey is used during registration process
func DecryptOTPKey(secretKey string, encryptedOTPKey string) (string, error) {
	return decryptData(secretKey, encryptedOTPKey)
}

func encryptData(password string, plaintext string) (string, error) {

	// allocate memory to hold the header of the ciphertext
	header := make([]byte, saltLen+aes.BlockSize)

	// generate salt
	salt := header[:saltLen]
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		panic(err)
	}

	// generate initialization vector
	iv := header[saltLen : aes.BlockSize+saltLen]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	// generate a 32 bit key with the provided password
	key := pbkdf2.Key([]byte(password), salt, iterations, keyLen, sha256.New)

	// generate a hmac for the message with the key
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(plaintext))
	hmac := mac.Sum(nil)

	// append this hmac to the plaintext
	plaintext = string(hmac) + plaintext

	//create the cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	// allocate space for the ciphertext and write the header to it
	ciphertext := make([]byte, len(header)+len(plaintext))
	copy(ciphertext, header)

	// encrypt
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize+saltLen:], []byte(plaintext))
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decryptData(password string, encrypted string) (string, error) {

	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)

	if err != nil {
		return "", fmt.Errorf("%w decrypt err %s", ErrServerInternal, err.Error())
	}
	// get the salt from the ciphertext
	salt := ciphertext[:saltLen]
	// get the IV from the ciphertext
	iv := ciphertext[saltLen : aes.BlockSize+saltLen]

	// generate the key with the KDF
	key := pbkdf2.Key([]byte(password), salt, iterations, keyLen, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("%w wrong key", ErrServerInternal)
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("%w wrong key length size", ErrServerInternal)
	}

	decrypted := ciphertext[saltLen+aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(decrypted, decrypted)

	// extract hmac from plaintext
	extractedMac := decrypted[:32]
	plaintext := decrypted[32:]

	// validate the hmac
	mac := hmac.New(sha256.New, key)
	mac.Write(plaintext)
	expectedMac := mac.Sum(nil)
	if !hmac.Equal(extractedMac, expectedMac) {
		return "", fmt.Errorf("%w hmac not equal", ErrClientDataIncorrect)
	}

	return string(plaintext), nil
}

// GenerateQR
func GenerateQR(issuer string, userEmail string) (string, []byte, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: userEmail,
		Period:      TOTPPeriod,
		Digits:      TOTPDigits,
		Algorithm:   TOTPAlgorithm,
	})

	if err != nil {
		return "", nil, fmt.Errorf("%w GenerateQR TOTP gererate err - %s", ErrServerInternal, err.Error())
	}

	var buf bytes.Buffer
	img, err := key.Image(450, 450)
	if err != nil {
		return "", nil, fmt.Errorf("%w GenerateQR TOTP image generation err - %s", ErrServerInternal, err.Error())
	}
	if err = png.Encode(&buf, img); err != nil {
		return "", nil, fmt.Errorf("%w GenerateQR png.Encode err - %s", ErrServerInternal, err.Error())
	}

	keyURL := key.URL()
	return keyURL, buf.Bytes(), nil
}

// ValidateOTPCode check otpPassCode
func ValidateOTPCode(keyURL string, otpPassCode string) (bool, error) {
	key, err := otp.NewKeyFromURL(keyURL)
	if err != nil {
		return false, fmt.Errorf("%w ValidatePassCode restore key err - %s", err, err.Error())
	}

	validOpts := totp.ValidateOpts{
		Period:    TOTPPeriod,
		Digits:    TOTPDigits,
		Algorithm: TOTPAlgorithm,
	}

	valid, err := totp.ValidateCustom(otpPassCode, key.Secret(), time.Now(), validOpts)
	return valid, err
}

// CreateJWTToken
func CreateJWTToken(tokenSecret []byte, tokenExp time.Duration, userID UserID) (JWTToken, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString(tokenSecret)
	if err != nil {
		return "", fmt.Errorf("%w: can't sign token %v", ErrServerInternal, err.Error())
	}

	return JWTToken(tokenString), nil
}

// ParseJWTToken
func ParseJWTToken(tokenSecret []byte, token JWTToken) (UserID, error) {
	claims := &Claims{}

	_, err := jwt.ParseWithClaims(string(token), claims, func(t *jwt.Token) (interface{}, error) {
		return tokenSecret, nil
	})

	if err != nil {
		return -1, fmt.Errorf("%w: authorization failed - %s", ErrAuthDataIncorrect, err.Error())
	}

	return claims.UserID, nil
}

func CheckBankCardData(data *BankCard) error {

	card := (*creditcard.Card)(data)

	if len(card.Number) != 16 { // The github.com/retgits/creditcard library has a bug if card.Number length too short.
		return fmt.Errorf("%w wrong card number length", ErrClientDataIncorrect)
	}
	validation := card.Validate()
	if len(validation.Errors) > 0 {
		errStr := strings.Join(validation.Errors, ", ")
		return fmt.Errorf("%w bank data validation err %s", ErrClientDataIncorrect, errStr)
	}
	return nil
}

func CheckUserPasswordData(data *UserPasswordData) error {
	if len(data.Hint) < 5 {
		return fmt.Errorf("%w hint too short", ErrClientDataIncorrect)
	}

	if len(data.Login) == 0 {
		return fmt.Errorf("%w login is not set", ErrClientDataIncorrect)
	}

	if len(data.Passwrod) == 0 {
		return fmt.Errorf("%w password is not set", ErrClientDataIncorrect)
	}

	return nil
}
