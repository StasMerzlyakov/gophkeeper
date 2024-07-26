package app_test

import (
	"bytes"
	"testing"

	"github.com/StasMerzlyakov/gophkeeper/internal/client/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testOKSaltFn = func(p []byte) (n int, err error) {
	for i := range p {
		p[i] = byte(i)
	}
	return len(p), nil
}

func TestCheckAuthPasswordComplexityLevel(t *testing.T) {
	helper := app.NewHelper(testOKSaltFn)

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
			ok := helper.CheckAuthPasswordComplexityLevel(test.pass)
			assert.Equal(t, test.res, ok)
		})
	}
}

func TestParseEMail(t *testing.T) {
	helper := app.NewHelper(testOKSaltFn)

	testData := []struct {
		name  string
		email string
		ok    bool
	}{
		{

			"ok",
			"test@gmail.com",
			true,
		},
		{

			"email bad",
			"test",
			false,
		},
	}

	for _, test := range testData {
		t.Run(test.name, func(t *testing.T) {
			res := helper.ParseEMail(test.email)
			assert.Equal(t, test.ok, res)
		})
	}
}

func TestCheckMasterKeyPasswordComplexityLevel(t *testing.T) {

	helper := app.NewHelper(testOKSaltFn)
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
			ok := helper.CheckMasterKeyPasswordComplexityLevel(test.pass)
			assert.Equal(t, test.res, ok)
		})
	}
}

func TestEncryptionText(t *testing.T) {
	helper := app.NewHelper(testOKSaltFn)
	passphrase := helper.Random32ByteString()

	randomText := "hello world"

	encrypted, err := helper.EncryptShortData([]byte(randomText), passphrase)
	require.NoError(t, err)

	data, err := helper.DecryptShortData(encrypted, passphrase)
	require.NoError(t, err)

	require.True(t, bytes.Equal([]byte(randomText), data))

	// check iv
	encrypted2, err := helper.EncryptShortData([]byte(randomText), passphrase)
	require.NoError(t, err)

	require.False(t, bytes.Equal([]byte(encrypted), []byte(encrypted2)))
}

func TestHelloWorld(t *testing.T) {
	helper := app.NewHelper(testOKSaltFn)

	generated, err := helper.GenerateHello()
	require.NoError(t, err)

	ok, err := helper.CheckHello(generated)
	require.NoError(t, err)
	require.True(t, ok)

}

func TestEncryptMasterKey(t *testing.T) {
	helper := app.NewHelper(testOKSaltFn)
	secretKey := "N1PCdw3M2B1TfJhoaY2mL736p2vCUc47"
	plainText := "testTestTest123"
	cipherText, err := helper.EncryptMasterKey(secretKey, plainText)
	require.NoError(t, err)
	require.True(t, len(cipherText) > 0)

	text, err := helper.DecryptMasterKey(secretKey, cipherText)
	require.NoError(t, err)
	require.Equal(t, plainText, text)
}
