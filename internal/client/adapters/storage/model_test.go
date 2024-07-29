package storage_test

import (
	"testing"

	"github.com/StasMerzlyakov/gophkeeper/internal/client/adapters/storage"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApp(t *testing.T) {

	t.Run("masterPassword", func(t *testing.T) {
		app := storage.NewStorage()
		mKey := "MasterPassword"
		app.SetMasterPassword(mKey)
		assert.Equal(t, mKey, app.GetMasterPassword())
	})

	t.Run("bank_card_operations", func(t *testing.T) {
		app := storage.NewStorage()

		require.Equal(t, 0, len(app.GetBankCardNumberList()))

		card := &domain.BankCard{
			Number:      "6250941006528599",
			ExpiryMonth: 06,
			ExpiryYear:  2026,
			CVV:         "123",
		}

		err := app.AddBankCard(card)
		require.NoError(t, err)

		require.Equal(t, 1, len(app.GetBankCardNumberList()))

		err = app.AddBankCard(card)
		require.ErrorIs(t, err, domain.ErrClientInternal)

		crd, err := app.GetBankCard(card.Number)
		require.NoError(t, err)

		require.Equal(t, crd.CVV, card.CVV)
		require.Equal(t, crd.Number, card.Number)
		require.Equal(t, crd.ExpiryMonth, card.ExpiryMonth)
		require.Equal(t, crd.ExpiryYear, card.ExpiryYear)

		err = app.DeleteBankCard("6250941006528598")
		require.ErrorIs(t, err, domain.ErrClientInternal)

		err = app.DeleteBankCard("6250941006528599")
		require.NoError(t, err)

		require.Equal(t, 0, len(app.GetBankCardNumberList()))
	})

	t.Run("user_data_operations", func(t *testing.T) {
		app := storage.NewStorage()

		require.Equal(t, 0, len(app.GetUserPasswordDataList()))

		data := &domain.UserPasswordData{
			Hint:     "ya.ru",
			Login:    "login",
			Passwrod: "pass",
		}

		err := app.AddUserPasswordData(data)
		require.NoError(t, err)

		require.Equal(t, 1, len(app.GetUserPasswordDataList()))

		err = app.AddUserPasswordData(data)
		require.ErrorIs(t, err, domain.ErrClientInternal)

		dt, err := app.GetUpdatePasswordData(data.Hint)
		require.NoError(t, err)

		require.Equal(t, dt.Hint, data.Hint)
		require.Equal(t, dt.Login, data.Login)
		require.Equal(t, dt.Passwrod, data.Passwrod)

		err = app.DeleteUpdatePasswordData("yayyy")
		require.ErrorIs(t, err, domain.ErrClientInternal)

		err = app.DeleteUpdatePasswordData("ya.ru")
		require.NoError(t, err)

		require.Equal(t, 0, len(app.GetUserPasswordDataList()))
	})
}
