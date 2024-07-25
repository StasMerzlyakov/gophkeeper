package domain_test

import (
	"bytes"
	"crypto/rand"
	"errors"
	"testing"
	"time"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAction(t *testing.T) {
	val := domain.GetAction(1)
	require.Equal(t, "domain_test.TestGetAction", val)
}

func TestValidateAuthPassword(t *testing.T) {
	testData := []struct {
		name string
		pass string
		res  bool
	}{
		{

			"ok",
			"IK0exasdF!",
			true,
		},
		{

			"bad",
			"123",
			false,
		},
	}

	for _, test := range testData {
		t.Run(test.name, func(t *testing.T) {
			ok := domain.CheckAuthPasswordComplexityLevel(test.pass)
			assert.Equal(t, test.res, ok)
		})
	}
}

func TestValidateEncryptionPassword(t *testing.T) {
	testData := []struct {
		name string
		pass string
		res  bool
	}{
		{
			"ok",
			"289!asdKeqvas!~",
			true,
		},
		{

			"bad",
			"IK0exasdF!",
			false,
		},
		{

			"bad2",
			"123",
			false,
		},
	}

	for _, test := range testData {
		t.Run(test.name, func(t *testing.T) {
			ok := domain.ValidateEncryptionPassword(test.pass)
			assert.Equal(t, test.res, ok)
		})
	}
}

func TestCheckEMailData(t *testing.T) {
	testData := []struct {
		name string
		data domain.EMailData
		err  error
	}{
		{

			"ok",
			domain.EMailData{
				EMail:    "test@gmail.com",
				Password: "IK0exasdF!",
			},
			nil,
		},
		{

			"email bad",
			domain.EMailData{
				EMail:    "test",
				Password: "IK0exasdF!",
			},
			domain.ErrClientDataIncorrect,
		},
		{

			"password bad",
			domain.EMailData{
				EMail:    "test@email",
				Password: "123",
			},
			domain.ErrClientDataIncorrect,
		},
	}

	for _, test := range testData {
		t.Run(test.name, func(t *testing.T) {
			res, err := domain.CheckEMailData(&test.data)
			if test.err == nil {
				assert.NoError(t, err)
				assert.True(t, res)
			} else {
				assert.ErrorIs(t, err, test.err)
				assert.False(t, res)
			}
		})
	}
}

var testOKSaltFn = func(p []byte) (n int, err error) {
	for i := range p {
		p[i] = byte(i)
	}
	return len(p), nil
}

var testErrSaltFn = func(p []byte) (n int, err error) {
	return -1, errors.New("test error")
}

func TestPasswordOperation(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		saltFn := testOKSaltFn

		pass := "12345678"
		hash, err := domain.HashPassword(pass, saltFn)
		require.NoError(t, err)

		require.NotNil(t, hash)

		assert.Equal(t, "In+BZhwpWKZH/S1QtMWcAOONZcrO9jVDaMDoJqgOfWM=", hash.Hash)
		assert.Equal(t, "AAECAwQFBgcICQoLDA0ODw==", hash.Salt)

		ok, err := domain.ValidateAccountPass(pass, hash.Hash, hash.Salt)
		require.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("errSaltFn", func(t *testing.T) {
		saltFn := testErrSaltFn
		pass := "12345678"
		hash, err := domain.HashPassword(pass, saltFn)
		require.Nil(t, hash)
		require.ErrorIs(t, err, domain.ErrServerInternal)
	})

	t.Run("errHashDecode", func(t *testing.T) {
		ok, err := domain.ValidateAccountPass("", "In+BZhwpWKZH/S1QtMWcAOONZcrO9jVDaMDoJqgOfWM", "")
		require.False(t, ok)
		assert.ErrorIs(t, err, domain.ErrClientDataIncorrect)
	})

	t.Run("errHashDecode", func(t *testing.T) {
		ok, err := domain.ValidateAccountPass("", "In+BZhwpWKZH/S1QtMWcAOONZcrO9jVDaMDoJqgOfWM=", "AAECAwQFBgcICQoLDA0ODw=")
		require.False(t, ok)
		assert.ErrorIs(t, err, domain.ErrClientDataIncorrect)
	})

	t.Run("password incorrect", func(t *testing.T) {
		ok, err := domain.ValidateAccountPass("123456789", "In+BZhwpWKZH/S1QtMWcAOONZcrO9jVDaMDoJqgOfWM=", "AAECAwQFBgcICQoLDA0ODw==")
		require.ErrorIs(t, err, domain.ErrAuthDataIncorrect)
		assert.False(t, ok)
	})
}

func TestCheckServerSecretKey(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		secretKey := "N1PCdw3M2B1TfJhoaY2mL736p2vCUc47"
		err := domain.CheckServerSecretKeyComplexityLevel(secretKey)
		require.NoError(t, err)
	})

	t.Run("wong length", func(t *testing.T) {
		secretKey := "N1PCdw3M2B1TfJhoaY2mL736p2vC"
		err := domain.CheckServerSecretKeyComplexityLevel(secretKey)
		require.Error(t, err)
	})

	t.Run("simple pass", func(t *testing.T) {
		secretKey := "12341234123412341234123412341234"
		err := domain.CheckServerSecretKeyComplexityLevel(secretKey)
		require.Error(t, err)
	})
}

func TestEncryptDecrypt(t *testing.T) {

	t.Run("ok", func(t *testing.T) {
		secretKey := "N1PCdw3M2B1TfJhoaY2mL736p2vCUc47"
		plainText := "testTestTest123"
		cipherText, err := domain.EncryptData(secretKey, plainText, testOKSaltFn)
		require.NoError(t, err)
		require.True(t, len(cipherText) > 0)

		text, err := domain.DecryptData(secretKey, cipherText)
		require.NoError(t, err)
		require.Equal(t, plainText, text)
	})

	t.Run("wrong pass encr", func(t *testing.T) {
		secretKey := "N1PCdw3M2B1TfJhoaY2mL736p2vCUc4"
		plainText := "testTestTest123"
		cipherText, err := domain.EncryptData(secretKey, plainText, testOKSaltFn)
		require.ErrorIs(t, err, domain.ErrServerInternal)
		assert.True(t, len(cipherText) == 0)
	})

	t.Run("wrong pass decr", func(t *testing.T) {
		secretKey := "N1PCdw3M2B1TfJhoaY2mL736p2vCUc47"
		plainText := "testTestTest123"
		cipherText, err := domain.EncryptData(secretKey, plainText, testOKSaltFn)
		require.NoError(t, err)
		require.True(t, len(cipherText) > 0)

		text, err := domain.DecryptData(secretKey[2:], cipherText)
		require.ErrorIs(t, err, domain.ErrServerInternal)
		assert.True(t, len(text) == 0)
	})

	t.Run("err", func(t *testing.T) {
		secretKey := "N1PCdw3M2B1TfJhoaY2mL736p2vCUc47"
		plainText := "testTestTest123"
		cipherText, err := domain.EncryptData(secretKey, plainText, testErrSaltFn)
		assert.True(t, len(cipherText) == 0)
		require.ErrorIs(t, err, domain.ErrServerInternal)
	})

}

func TestGenerateQR(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		issuer := "issuer"
		accountName := "accountName"

		key, png, err := domain.GenerateQR(issuer, accountName)
		require.NoError(t, err)
		require.NotEmpty(t, key)
		require.NotNil(t, png)
	})

	t.Run("err", func(t *testing.T) {
		issuer := ""
		accountName := "accountName"

		key, png, err := domain.GenerateQR(issuer, accountName)
		require.Error(t, err)
		require.Empty(t, key)
		require.Nil(t, png)
	})

	t.Run("err2", func(t *testing.T) {
		issuer := "issuer"
		accountName := ""

		key, png, err := domain.GenerateQR(issuer, accountName)
		require.Error(t, err)
		require.Empty(t, key)
		require.Nil(t, png)
	})
}

func TestValidate(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		issuer := "issuer"
		accountName := "accountName"

		keyURL, _, err := domain.GenerateQR(issuer, accountName)

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

		ok, err := domain.ValidatePassCode(key.URL(), passcode)
		require.NoError(t, err)

		assert.True(t, ok)
	})

	t.Run("err", func(t *testing.T) {
		_, err := domain.ValidatePassCode(":\\", "12345")
		require.Error(t, err)
	})

	t.Run("err2", func(t *testing.T) {
		issuer := "issuer"
		accountName := ""

		key, png, err := domain.GenerateQR(issuer, accountName)
		require.Error(t, err)
		require.Empty(t, key)
		require.Nil(t, png)
	})
}

func TestHelloWorld(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		generated, err := domain.GenerateHello(rand.Read)
		require.NoError(t, err)

		ok, err := domain.CheckHello(generated)
		require.NoError(t, err)
		require.True(t, ok)
	})

	t.Run("saltErr", func(t *testing.T) {
		generated, err := domain.GenerateHello(testErrSaltFn)
		require.ErrorIs(t, err, domain.ErrServerInternal)
		assert.Empty(t, generated)
	})

	t.Run("false", func(t *testing.T) {
		ok, err := domain.CheckHello("jFG4UOPJEJoS+tm3Z5cGZ2hlbGxvIGZyb20gR29waEtlZXBlciEhIQ==")
		require.NoError(t, err)
		require.False(t, ok)
	})

	t.Run("b64_err", func(t *testing.T) {
		ok, err := domain.CheckHello("9COumBoRUEFVXbFYg5LM1GhlbGxvIGZyb20gR29waEtlRXBlciEh=")
		require.Error(t, err, domain.ErrClientDataIncorrect)
		require.False(t, ok)
	})

	t.Run("wrong_length", func(t *testing.T) {
		ok, err := domain.CheckHello("9COu")
		require.Error(t, err, domain.ErrClientDataIncorrect)
		require.False(t, ok)
	})

	//	ok, err = domain.CheckHelloWorld()

}

func TestJWT(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		tokenSecret := "tokenSecret"
		tokenDuration := 2 * time.Second
		userID := domain.UserID(1)

		jwtTok, err := domain.CreateJWTToken([]byte(tokenSecret), tokenDuration, userID)
		require.NoError(t, err)
		require.NotEmpty(t, jwtTok)

		uID, err := domain.ParseJWTToken([]byte(tokenSecret), jwtTok)
		require.NoError(t, err)
		require.Equal(t, userID, uID)
	})

	t.Run("timeout", func(t *testing.T) {
		tokenSecret := "tokenSecret"
		tokenDuration := 1 * time.Second
		userID := domain.UserID(1)

		jwtTok, err := domain.CreateJWTToken([]byte(tokenSecret), tokenDuration, userID)
		require.NoError(t, err)
		require.NotEmpty(t, jwtTok)

		time.Sleep(2 * time.Second)

		_, err = domain.ParseJWTToken([]byte(tokenSecret), jwtTok)
		require.ErrorIs(t, err, domain.ErrAuthDataIncorrect)
	})

	t.Run("wrong_secret", func(t *testing.T) {
		tokenSecret := "tokenSecret"
		tokenDuration := 1 * time.Second
		userID := domain.UserID(1)

		jwtTok, err := domain.CreateJWTToken([]byte(tokenSecret), tokenDuration, userID)
		require.NoError(t, err)
		require.NotEmpty(t, jwtTok)

		time.Sleep(2 * time.Second)

		tokenSecret = tokenSecret + "nonce"
		_, err = domain.ParseJWTToken([]byte(tokenSecret), jwtTok)
		require.ErrorIs(t, err, domain.ErrAuthDataIncorrect)
	})
}

func TestEncryotAES256(t *testing.T) {
	passphrase := domain.Random32ByteString()

	randomText := "hello world"

	encrypted, err := domain.EncryptAES256([]byte(randomText), passphrase)
	require.NoError(t, err)

	data, err := domain.DecryptAES256(encrypted, passphrase)
	require.NoError(t, err)

	require.True(t, bytes.Equal([]byte(randomText), data))

	// check iv
	encrypted2, err := domain.EncryptAES256([]byte(randomText), passphrase)
	require.NoError(t, err)

	require.False(t, bytes.Equal([]byte(encrypted), []byte(encrypted2)))

}
