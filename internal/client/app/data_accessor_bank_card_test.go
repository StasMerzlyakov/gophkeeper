package app_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/StasMerzlyakov/gophkeeper/internal/client/app"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetBankCardList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {

		masterKey := "masterKey"
		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().GetMasterPassword().Return(masterKey).Times(1)

		mockServer := NewMockAppServer(ctrl)

		bankCardList := []domain.EncryptedBankCard{
			{
				Number:  "100",
				Content: "content",
			},
		}
		mockServer.EXPECT().GetBankCardList(gomock.Any()).Return(bankCardList, nil).Times(1)

		mockHelper := NewMockDomainHelper(ctrl)
		mockHelper.EXPECT().DecryptShortData(gomock.Any(), gomock.Any()).DoAndReturn(func(key string, cnt string) (string, error) {
			assert.Equal(t, masterKey, key)
			assert.Equal(t, bankCardList[0].Content, cnt)
			return `{"number":"100", "type":"MIR"}`, nil
		}).Times(1)

		mockStorage.EXPECT().SetBankCards(gomock.Any()).Do(func(crds []domain.BankCard) {
			require.Equal(t, 1, len(crds))
			crd := crds[0]
			assert.Equal(t, "100", crd.Number)
			assert.Equal(t, "MIR", crd.Type)

			assert.Equal(t, 0, crd.ExpiryMonth)
			assert.Equal(t, 0, crd.ExpiryYear)
			assert.Equal(t, "", crd.CVV)
		}).Times(1)

		da := app.NewDataAccessor().AppSever(mockServer).AppStorage(mockStorage).DomainHelper(mockHelper)
		err := da.GetBankCardList(context.Background())
		require.NoError(t, err)
	})

	t.Run("get_bank_card_list_err", func(t *testing.T) {

		masterKey := "masterKey"
		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().GetMasterPassword().Return(masterKey).Times(1)

		mockServer := NewMockAppServer(ctrl)

		testErr := errors.New("testErr")
		mockServer.EXPECT().GetBankCardList(gomock.Any()).Return(nil, testErr).Times(1)

		da := app.NewDataAccessor().AppSever(mockServer).AppStorage(mockStorage)
		err := da.GetBankCardList(context.Background())
		require.ErrorIs(t, err, testErr)
	})

	t.Run("decrypt_err", func(t *testing.T) {

		masterKey := "masterKey"
		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().GetMasterPassword().Return(masterKey).Times(1)

		mockServer := NewMockAppServer(ctrl)

		bankCardList := []domain.EncryptedBankCard{
			{
				Number:  "100",
				Content: "content",
			},
		}
		mockServer.EXPECT().GetBankCardList(gomock.Any()).Return(bankCardList, nil).Times(1)

		mockHelper := NewMockDomainHelper(ctrl)
		testErr := errors.New("testErr")
		mockHelper.EXPECT().DecryptShortData(gomock.Any(), gomock.Any()).DoAndReturn(func(key string, cnt string) (string, error) {
			assert.Equal(t, masterKey, key)
			assert.Equal(t, bankCardList[0].Content, cnt)
			return `{"number":"100", "type":"MIR"}`, testErr
		}).Times(1)

		da := app.NewDataAccessor().AppSever(mockServer).AppStorage(mockStorage).DomainHelper(mockHelper)
		err := da.GetBankCardList(context.Background())
		require.ErrorIs(t, err, domain.ErrClientDataIncorrect)
	})

	t.Run("json_err", func(t *testing.T) {

		masterKey := "masterKey"
		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().GetMasterPassword().Return(masterKey).Times(1)

		mockServer := NewMockAppServer(ctrl)

		bankCardList := []domain.EncryptedBankCard{
			{
				Number:  "100",
				Content: "content",
			},
		}
		mockServer.EXPECT().GetBankCardList(gomock.Any()).Return(bankCardList, nil).Times(1)

		mockHelper := NewMockDomainHelper(ctrl)
		mockHelper.EXPECT().DecryptShortData(gomock.Any(), gomock.Any()).DoAndReturn(func(key string, cnt string) (string, error) {
			assert.Equal(t, masterKey, key)
			assert.Equal(t, bankCardList[0].Content, cnt)
			return `-`, nil
		}).Times(1)

		da := app.NewDataAccessor().AppSever(mockServer).AppStorage(mockStorage).DomainHelper(mockHelper)
		err := da.GetBankCardList(context.Background())
		require.ErrorIs(t, err, domain.ErrClientDataIncorrect)
	})
}

func TestAddBankCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {

		mockHelper := NewMockDomainHelper(ctrl)

		bankCard := &domain.BankCard{
			Number: "100",
			Type:   "MIR",
		}

		mockHelper.EXPECT().CheckBankCardData(gomock.Any()).DoAndReturn(func(crd *domain.BankCard) error {
			assert.Same(t, bankCard, crd)
			return nil
		}).Times(1)

		masterKey := "masterKey"
		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().GetMasterPassword().Return(masterKey).Times(1)

		encryptedCnt := "encrypted"
		mockHelper.EXPECT().EncryptShortData(gomock.Any(), gomock.Any()).DoAndReturn(func(key string, content string) (string, error) {

			assert.Equal(t, masterKey, key)
			cnt, err := json.Marshal(bankCard)
			require.NoError(t, err)
			assert.JSONEq(t, string(cnt), content)
			return encryptedCnt, nil
		}).Times(1)

		mockServer := NewMockAppServer(ctrl)
		mockServer.EXPECT().CreateBankCard(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, bnkCard *domain.EncryptedBankCard) error {
			require.NotNil(t, bnkCard)
			assert.Equal(t, bankCard.Number, bnkCard.Number)
			assert.Equal(t, encryptedCnt, bnkCard.Content)
			return nil
		}).Times(1)

		da := app.NewDataAccessor().AppSever(mockServer).AppStorage(mockStorage).DomainHelper(mockHelper)
		err := da.AddBankCard(context.Background(), bankCard)
		require.NoError(t, err)
	})

	t.Run("check_card_err", func(t *testing.T) {

		mockHelper := NewMockDomainHelper(ctrl)

		bankCard := &domain.BankCard{
			Number: "100",
			Type:   "MIR",
		}

		testErr := errors.New("testErr")
		mockHelper.EXPECT().CheckBankCardData(gomock.Any()).DoAndReturn(func(crd *domain.BankCard) error {
			assert.Same(t, bankCard, crd)
			return testErr
		}).Times(1)

		da := app.NewDataAccessor().DomainHelper(mockHelper)
		err := da.AddBankCard(context.Background(), bankCard)
		require.ErrorIs(t, err, domain.ErrClientDataIncorrect)
	})

	t.Run("encrypt_err", func(t *testing.T) {

		mockHelper := NewMockDomainHelper(ctrl)

		bankCard := &domain.BankCard{
			Number: "100",
			Type:   "MIR",
		}

		mockHelper.EXPECT().CheckBankCardData(gomock.Any()).DoAndReturn(func(crd *domain.BankCard) error {
			assert.Same(t, bankCard, crd)
			return nil
		}).Times(1)

		masterKey := "masterKey"
		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().GetMasterPassword().Return(masterKey).Times(1)

		testErr := errors.New("testErr")
		mockHelper.EXPECT().EncryptShortData(gomock.Any(), gomock.Any()).DoAndReturn(func(key string, content string) (string, error) {

			assert.Equal(t, masterKey, key)
			cnt, err := json.Marshal(bankCard)
			require.NoError(t, err)
			assert.JSONEq(t, string(cnt), content)
			return "", testErr
		}).Times(1)

		da := app.NewDataAccessor().AppStorage(mockStorage).DomainHelper(mockHelper)
		err := da.AddBankCard(context.Background(), bankCard)
		require.ErrorIs(t, err, domain.ErrClientInternal)
	})

	t.Run("server_err", func(t *testing.T) {

		mockHelper := NewMockDomainHelper(ctrl)

		bankCard := &domain.BankCard{
			Number: "100",
			Type:   "MIR",
		}

		mockHelper.EXPECT().CheckBankCardData(gomock.Any()).DoAndReturn(func(crd *domain.BankCard) error {
			assert.Same(t, bankCard, crd)
			return nil
		}).Times(1)

		masterKey := "masterKey"
		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().GetMasterPassword().Return(masterKey).Times(1)

		encryptedCnt := "encrypted"
		mockHelper.EXPECT().EncryptShortData(gomock.Any(), gomock.Any()).DoAndReturn(func(key string, content string) (string, error) {

			assert.Equal(t, masterKey, key)
			cnt, err := json.Marshal(bankCard)
			require.NoError(t, err)
			assert.JSONEq(t, string(cnt), content)
			return encryptedCnt, nil
		}).Times(1)

		mockServer := NewMockAppServer(ctrl)
		testErr := errors.New("testErr")
		mockServer.EXPECT().CreateBankCard(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, bnkCard *domain.EncryptedBankCard) error {
			require.NotNil(t, bnkCard)
			assert.Equal(t, bankCard.Number, bnkCard.Number)
			assert.Equal(t, encryptedCnt, bnkCard.Content)
			return testErr
		}).Times(1)

		da := app.NewDataAccessor().AppSever(mockServer).AppStorage(mockStorage).DomainHelper(mockHelper)
		err := da.AddBankCard(context.Background(), bankCard)
		require.ErrorIs(t, err, testErr)
	})
}

func TestUpdateBankCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {

		mockHelper := NewMockDomainHelper(ctrl)

		bankCard := &domain.BankCard{
			Number: "100",
			Type:   "MIR",
		}

		mockHelper.EXPECT().CheckBankCardData(gomock.Any()).DoAndReturn(func(crd *domain.BankCard) error {
			assert.Same(t, bankCard, crd)
			return nil
		}).Times(1)

		masterKey := "masterKey"
		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().GetMasterPassword().Return(masterKey).Times(1)

		encryptedCnt := "encrypted"
		mockHelper.EXPECT().EncryptShortData(gomock.Any(), gomock.Any()).DoAndReturn(func(key string, content string) (string, error) {

			assert.Equal(t, masterKey, key)
			cnt, err := json.Marshal(bankCard)
			require.NoError(t, err)
			assert.JSONEq(t, string(cnt), content)
			return encryptedCnt, nil
		}).Times(1)

		mockServer := NewMockAppServer(ctrl)
		mockServer.EXPECT().UpdateBankCard(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, bnkCard *domain.EncryptedBankCard) error {
			require.NotNil(t, bnkCard)
			assert.Equal(t, bankCard.Number, bnkCard.Number)
			assert.Equal(t, encryptedCnt, bnkCard.Content)
			return nil
		}).Times(1)

		da := app.NewDataAccessor().AppSever(mockServer).AppStorage(mockStorage).DomainHelper(mockHelper)
		err := da.UpdateBankCard(context.Background(), bankCard)
		require.NoError(t, err)
	})

	t.Run("check_card_err", func(t *testing.T) {

		mockHelper := NewMockDomainHelper(ctrl)

		bankCard := &domain.BankCard{
			Number: "100",
			Type:   "MIR",
		}

		testErr := errors.New("testErr")
		mockHelper.EXPECT().CheckBankCardData(gomock.Any()).DoAndReturn(func(crd *domain.BankCard) error {
			assert.Same(t, bankCard, crd)
			return testErr
		}).Times(1)

		da := app.NewDataAccessor().DomainHelper(mockHelper)
		err := da.AddBankCard(context.Background(), bankCard)
		require.ErrorIs(t, err, domain.ErrClientDataIncorrect)
	})

	t.Run("emcrypt_err", func(t *testing.T) {

		mockHelper := NewMockDomainHelper(ctrl)

		bankCard := &domain.BankCard{
			Number: "100",
			Type:   "MIR",
		}

		mockHelper.EXPECT().CheckBankCardData(gomock.Any()).DoAndReturn(func(crd *domain.BankCard) error {
			assert.Same(t, bankCard, crd)
			return nil
		}).Times(1)

		masterKey := "masterKey"
		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().GetMasterPassword().Return(masterKey).Times(1)

		testErr := errors.New("testErr")
		mockHelper.EXPECT().EncryptShortData(gomock.Any(), gomock.Any()).DoAndReturn(func(key string, content string) (string, error) {

			assert.Equal(t, masterKey, key)
			cnt, err := json.Marshal(bankCard)
			require.NoError(t, err)
			assert.JSONEq(t, string(cnt), content)
			return "", testErr
		}).Times(1)

		da := app.NewDataAccessor().AppStorage(mockStorage).DomainHelper(mockHelper)
		err := da.AddBankCard(context.Background(), bankCard)
		require.ErrorIs(t, err, domain.ErrClientInternal)
	})

	t.Run("server_err", func(t *testing.T) {

		mockHelper := NewMockDomainHelper(ctrl)

		bankCard := &domain.BankCard{
			Number: "100",
			Type:   "MIR",
		}

		mockHelper.EXPECT().CheckBankCardData(gomock.Any()).DoAndReturn(func(crd *domain.BankCard) error {
			assert.Same(t, bankCard, crd)
			return nil
		}).Times(1)

		masterKey := "masterKey"
		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().GetMasterPassword().Return(masterKey).Times(1)

		encryptedCnt := "encrypted"
		mockHelper.EXPECT().EncryptShortData(gomock.Any(), gomock.Any()).DoAndReturn(func(key string, content string) (string, error) {

			assert.Equal(t, masterKey, key)
			cnt, err := json.Marshal(bankCard)
			require.NoError(t, err)
			assert.JSONEq(t, string(cnt), content)
			return encryptedCnt, nil
		}).Times(1)

		mockServer := NewMockAppServer(ctrl)
		testErr := errors.New("testErr")
		mockServer.EXPECT().UpdateBankCard(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, bnkCard *domain.EncryptedBankCard) error {
			require.NotNil(t, bnkCard)
			assert.Equal(t, bankCard.Number, bnkCard.Number)
			assert.Equal(t, encryptedCnt, bnkCard.Content)
			return testErr
		}).Times(1)

		da := app.NewDataAccessor().AppSever(mockServer).AppStorage(mockStorage).DomainHelper(mockHelper)
		err := da.UpdateBankCard(context.Background(), bankCard)
		require.ErrorIs(t, err, testErr)
	})
}

func TestDeleteBankCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {

		number := "100"

		mockServer := NewMockAppServer(ctrl)
		mockServer.EXPECT().DeleteBankCard(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, nmb string) error {
			assert.Equal(t, number, nmb)
			return nil
		}).Times(1)

		da := app.NewDataAccessor().AppSever(mockServer)
		err := da.DeleteBankCard(context.Background(), number)
		require.NoError(t, err)
	})

	t.Run("server_err", func(t *testing.T) {

		number := "100"

		mockServer := NewMockAppServer(ctrl)
		testErr := errors.New("testErr")
		mockServer.EXPECT().DeleteBankCard(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, nmb string) error {
			assert.Equal(t, number, nmb)
			return testErr
		}).Times(1)

		da := app.NewDataAccessor().AppSever(mockServer)
		err := da.DeleteBankCard(context.Background(), number)
		require.ErrorIs(t, err, testErr)

	})
}
