package domain_test

import (
	"errors"
	"testing"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAction(t *testing.T) {
	val := domain.GetAction(1)
	require.Equal(t, "domain_test.TestGetAction", val)
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
			domain.ErrDataFormat,
		},
		{

			"password bad",
			domain.EMailData{
				EMail:    "test@email",
				Password: "123",
			},
			domain.ErrDataFormat,
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

		assert.Equal(t, "In-BZhwpWKZH_S1QtMWcAOONZcrO9jVDaMDoJqgOfWM=", hash.Hash)
		assert.Equal(t, "AAECAwQFBgcICQoLDA0ODw==", hash.Salt)

		ok, err := domain.CheckPassword(pass, hash.Hash, hash.Salt)
		require.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("errSaltFn", func(t *testing.T) {
		saltFn := testErrSaltFn
		pass := "12345678"
		hash, err := domain.HashPassword(pass, saltFn)
		require.Nil(t, hash)
		require.ErrorIs(t, err, domain.ErrInternalServer)
	})

	t.Run("errHashDecode", func(t *testing.T) {
		ok, err := domain.CheckPassword("", "In-BZhwpWKZH_S1QtMWcAOONZcrO9jVDaMDoJqgOfWM", "")
		require.False(t, ok)
		assert.ErrorIs(t, err, domain.ErrDataFormat)
	})

	t.Run("errHashDecode", func(t *testing.T) {
		ok, err := domain.CheckPassword("", "In-BZhwpWKZH_S1QtMWcAOONZcrO9jVDaMDoJqgOfWM=", "AAECAwQFBgcICQoLDA0ODw=")
		require.False(t, ok)
		assert.ErrorIs(t, err, domain.ErrDataFormat)
	})

	t.Run("password incorrect", func(t *testing.T) {
		ok, err := domain.CheckPassword("123456789", "In-BZhwpWKZH_S1QtMWcAOONZcrO9jVDaMDoJqgOfWM=", "AAECAwQFBgcICQoLDA0ODw==")
		require.ErrorIs(t, err, domain.ErrWrongLoginPassword)
		assert.False(t, ok)
	})
}

func TestCheckServerSecretKey(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		secretKey := "N1PCdw3M2B1TfJhoaY2mL736p2vCUc47"
		err := domain.CheckServerSecretKey(secretKey)
		require.NoError(t, err)
	})

	t.Run("wong length", func(t *testing.T) {
		secretKey := "N1PCdw3M2B1TfJhoaY2mL736p2vC"
		err := domain.CheckServerSecretKey(secretKey)
		require.Error(t, err)
	})

	t.Run("simple pass", func(t *testing.T) {
		secretKey := "12341234123412341234123412341234"
		err := domain.CheckServerSecretKey(secretKey)
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
		require.ErrorIs(t, err, domain.ErrInternalServer)
		assert.True(t, len(cipherText) == 0)
	})

	t.Run("wrong pass decr", func(t *testing.T) {
		secretKey := "N1PCdw3M2B1TfJhoaY2mL736p2vCUc47"
		plainText := "testTestTest123"
		cipherText, err := domain.EncryptData(secretKey, plainText, testOKSaltFn)
		require.NoError(t, err)
		require.True(t, len(cipherText) > 0)

		text, err := domain.DecryptData(secretKey[2:], cipherText)
		require.ErrorIs(t, err, domain.ErrInternalServer)
		assert.True(t, len(text) == 0)
	})

	t.Run("err", func(t *testing.T) {
		secretKey := "N1PCdw3M2B1TfJhoaY2mL736p2vCUc47"
		plainText := "testTestTest123"
		cipherText, err := domain.EncryptData(secretKey, plainText, testErrSaltFn)
		assert.True(t, len(cipherText) == 0)
		require.ErrorIs(t, err, domain.ErrInternalServer)
	})

}
