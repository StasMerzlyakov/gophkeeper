package usecases_test

import (
	"testing"
	"time"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/StasMerzlyakov/gophkeeper/internal/server/usecases"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testOKSaltFn = func(p []byte) (n int, err error) {
	for i := range p {
		p[i] = byte(i)
	}
	return len(p), nil
}

func TestRegistrationHelper(t *testing.T) {
	helper := usecases.NewRegistrationHelper(testOKSaltFn)
	t.Run("hash_password", func(t *testing.T) {

		pass := "12345678"
		hash, err := helper.HashPassword(pass)
		require.NoError(t, err)

		require.NotNil(t, hash)

		assert.Equal(t, "In-BZhwpWKZH_S1QtMWcAOONZcrO9jVDaMDoJqgOfWM=", hash.Hash)
		assert.Equal(t, "AAECAwQFBgcICQoLDA0ODw==", hash.Salt)

		ok, err := helper.CheckPassword(pass, hash.Hash, hash.Salt)
		require.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("check_email_data", func(t *testing.T) {
		data := &domain.EMailData{
			EMail:    "test@gmail.com",
			Password: "IK0exasdF!",
		}

		res, err := helper.CheckEMailData(data)

		assert.NoError(t, err)
		assert.True(t, res)
	})

	t.Run("generate_qr", func(t *testing.T) {
		issuer := "issuer"
		accountName := "accountName"

		key, png, err := helper.GenerateQR(issuer, accountName)
		require.NoError(t, err)
		require.NotEmpty(t, key)
		require.NotNil(t, png)
	})

	t.Run("encrypt_decrypt", func(t *testing.T) {
		secretKey := "N1PCdw3M2B1TfJhoaY2mL736p2vCUc47"
		plainText := "testTestTest123"
		cipherText, err := helper.EncryptData(secretKey, plainText)
		require.NoError(t, err)
		require.True(t, len(cipherText) > 0)

		text, err := helper.DecryptData(secretKey, cipherText)
		require.NoError(t, err)
		require.Equal(t, plainText, text)
	})

	t.Run("session_id", func(t *testing.T) {
		sessID := helper.NewSessionID()
		require.True(t, len(sessID) > 0)
	})

	t.Run("valudate_pass", func(t *testing.T) {
		issuer := "issuer"
		accountName := "accountName"

		keyURL, _, err := helper.GenerateQR(issuer, accountName)

		require.NoError(t, err)

		key, err := otp.NewKeyFromURL(keyURL)
		require.NoError(t, err)

		validOpts := totp.ValidateOpts{
			Period:    domain.TOTPPeriod,
			Digits:    domain.TOTPDigits,
			Algorithm: domain.TOTPAlgorithm,
		}
		passcode, err := totp.GenerateCodeCustom(key.Secret(), time.Now(), validOpts)
		require.NoError(t, err)

		ok, err := helper.ValidatePassCode(keyURL, passcode)
		require.NoError(t, err)

		assert.True(t, ok)
	})

	t.Run("hello_test", func(t *testing.T) {
		hello, err := helper.GenerateHello()
		require.NoError(t, err)
		ok, err := helper.CheckHello(hello)
		require.NoError(t, err)
		require.True(t, ok)
	})
}
