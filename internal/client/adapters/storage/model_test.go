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

		app.SetBankCards([]domain.BankCard{
			{
				Number:      "6250941006528599",
				ExpiryMonth: 06,
				ExpiryYear:  2026,
				CVV:         "123",
			},
		})
		require.Equal(t, 1, len(app.GetBankCardNumberList()))

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

		dt, err := app.GetUserPasswordData(data.Hint)
		require.NoError(t, err)

		require.Equal(t, dt.Hint, data.Hint)
		require.Equal(t, dt.Login, data.Login)
		require.Equal(t, dt.Passwrod, data.Passwrod)

		err = app.DeleteUserPasswordData("yayyy")
		require.ErrorIs(t, err, domain.ErrClientInternal)

		err = app.DeleteUserPasswordData("ya.ru")
		require.NoError(t, err)

		require.Equal(t, 0, len(app.GetUserPasswordDataList()))

		app.SetUserPasswordDatas([]domain.UserPasswordData{
			{
				Hint:     "ya.ru",
				Login:    "login",
				Passwrod: "pass",
			},
		})
		require.Equal(t, 1, len(app.GetUserPasswordDataList()))
	})

	t.Run("files_operations", func(t *testing.T) {
		app := storage.NewStorage()

		require.Equal(t, 0, len(app.GetFileInfoList()))

		data := &domain.FileInfo{
			Name: "ya.ru",
		}

		err := app.AddFileInfo(data)
		require.NoError(t, err)

		require.Equal(t, 1, len(app.GetFileInfoList()))

		err = app.AddFileInfo(data)
		require.ErrorIs(t, err, domain.ErrClientInternal)

		dt, err := app.GetFileInfo(data.Name)
		require.NoError(t, err)

		require.Equal(t, dt.Name, data.Name)
		require.Equal(t, dt.Path, data.Path)

		data = &domain.FileInfo{
			Name: "ya.ru",
			Path: "path",
		}
		err = app.UpdateFileInfo(data)
		require.NoError(t, err)

		dt, err = app.GetFileInfo(data.Name)
		require.NoError(t, err)

		ok := app.IsFileInfoExists(data.Name)
		require.True(t, ok)

		require.Equal(t, dt.Name, data.Name)
		require.Equal(t, dt.Path, data.Path)

		err = app.DeleteFileInfo("yayyy")
		require.ErrorIs(t, err, domain.ErrClientInternal)

		err = app.DeleteFileInfo("ya.ru")
		require.NoError(t, err)

		require.Equal(t, 0, len(app.GetFileInfoList()))

		app.SetUserPasswordDatas([]domain.UserPasswordData{
			{
				Hint:     "ya.ru",
				Login:    "login",
				Passwrod: "pass",
			},
		})
		require.Equal(t, 1, len(app.GetUserPasswordDataList()))
	})
}
