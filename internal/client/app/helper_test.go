package app_test

import (
	"bytes"
	"testing"

	"github.com/StasMerzlyakov/gophkeeper/internal/client/app"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
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

func TestCheckMasterPasswordComplexityLevel(t *testing.T) {

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
			ok := helper.CheckMasterPasswordComplexityLevel(test.pass)
			assert.Equal(t, test.res, ok)
		})
	}
}

func TestEncryptionHello(t *testing.T) {
	helper := app.NewHelper(testOKSaltFn)
	passphrase := helper.Random32ByteString()

	masterPass := "hello world"

	encrypted, err := helper.EncryptHello(masterPass, passphrase)
	require.NoError(t, err)

	err = helper.DecryptHello(masterPass, encrypted)
	require.NoError(t, err)
}

func TestEncryptionText(t *testing.T) {
	passphrase := domain.Random32ByteString()
	helper := app.NewHelper(testOKSaltFn)

	randomText := "hello world"

	encrypted, err := helper.EncryptShortData(passphrase, randomText)
	require.NoError(t, err)

	data, err := helper.DecryptShortData(passphrase, encrypted)
	require.NoError(t, err)

	require.Equal(t, randomText, data)

	// check iv
	encrypted2, err := helper.EncryptShortData(passphrase, randomText)
	require.NoError(t, err)

	require.False(t, bytes.Equal([]byte(encrypted), []byte(encrypted2)))
}
