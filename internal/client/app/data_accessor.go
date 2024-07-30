package app

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

func NewDataAccessor() *dataAccessor {
	return &dataAccessor{}
}

func (dcc *dataAccessor) AppSever(appServer AppServer) *dataAccessor {
	dcc.appServer = appServer
	return dcc
}

func (dcc *dataAccessor) DomainHelper(helper DomainHelper) *dataAccessor {
	dcc.helper = helper
	return dcc
}

func (dcc *dataAccessor) AppStorage(appStorage AppStorage) *dataAccessor {
	dcc.appStorage = appStorage
	return dcc
}

type dataAccessor struct {
	appServer  AppServer
	helper     DomainHelper
	appStorage AppStorage
}

func (dcc *dataAccessor) GetBankCardList(ctx context.Context) error {
	log := GetMainLogger()
	action := domain.GetAction(1)
	log.Debugf("%v start", action)

	masterPass := dcc.appStorage.GetMasterPassword()

	encrList, err := dcc.appServer.GetBankCardList(ctx)
	if err != nil {
		err := fmt.Errorf("%w - %v error", err, action)
		log.Warn(err.Error())
		return err
	}

	var decryptedList []domain.BankCard
	for _, encr := range encrList {
		content := encr.Content

		encrypted, err := dcc.helper.DecryptShortData(masterPass, content)
		if err != nil {
			err := fmt.Errorf("%w - %v error - can't encrypt card with number %v", err, action, encr.Number)
			log.Warn(err.Error())
			return err
		}

		var bankCard domain.BankCard
		if err := json.Unmarshal([]byte(encrypted), &bankCard); err != nil {
			err := fmt.Errorf("%w - %v error - can't decode card data with number %v", err, action, encr.Number)
			log.Warn(err.Error())
			return err
		}
		decryptedList = append(decryptedList, bankCard)
	}

	dcc.appStorage.SetBankCards(decryptedList)

	log.Debugf("%v success", action)
	return nil
}

func (dcc *dataAccessor) AddBankCard(ctx context.Context, bankCard *domain.BankCard) error {
	log := GetMainLogger()
	action := domain.GetAction(1)
	log.Debugf("%v start", action)

	res, err := json.Marshal(bankCard)
	if err != nil {
		err := fmt.Errorf("%w - %v error - can't marshal data %v", domain.ErrClientInternal, action, err.Error())
		log.Warn(err.Error())
		return err
	}
	masterPass := dcc.appStorage.GetMasterPassword()

	content, err := dcc.helper.EncryptShortData(masterPass, string(res))

	if err != nil {
		err := fmt.Errorf("%w - %v error - can't encrypt data %v", domain.ErrClientInternal, action, err.Error())
		log.Warn(err.Error())
		return err
	}

	err = dcc.appServer.CreateBankCard(ctx, &domain.EncryptedBankCard{
		Number:  bankCard.Number,
		Content: content,
	})

	if err != nil {
		err := fmt.Errorf("%w - %v error", err, action)
		log.Warn(err.Error())
		return err
	}

	log.Debugf("%v success", action)
	return nil
}

func (dcc *dataAccessor) UpdateBankCard(ctx context.Context, bankCard *domain.BankCard) error {
	log := GetMainLogger()
	action := domain.GetAction(1)
	log.Debugf("%v start", action)
	res, err := json.Marshal(bankCard)
	if err != nil {
		err := fmt.Errorf("%w - %v error - can't marshal data %v", domain.ErrClientInternal, action, err.Error())
		log.Warn(err.Error())
		return err
	}
	masterPass := dcc.appStorage.GetMasterPassword()

	content, err := dcc.helper.EncryptShortData(masterPass, string(res))

	if err != nil {
		err := fmt.Errorf("%w - %v error - can't encrypt data %v", domain.ErrClientInternal, action, err.Error())
		log.Warn(err.Error())
		return err
	}

	err = dcc.appServer.UpdateBankCard(ctx, &domain.EncryptedBankCard{
		Number:  bankCard.Number,
		Content: content,
	})

	if err != nil {
		err := fmt.Errorf("%w - %v error", err, action)
		log.Warn(err.Error())
		return err
	}

	log.Debugf("%v success", action)
	return nil
}

func (dcc *dataAccessor) DeleteBankCard(ctx context.Context, number string) error {
	log := GetMainLogger()
	action := domain.GetAction(1)
	log.Debugf("%v start", action)

	if err := dcc.appServer.DeleteBankCard(ctx, number); err != nil {
		err := fmt.Errorf("%w - %v error", err, action)
		log.Warn(err.Error())
		return err
	}

	log.Debugf("%v success", action)
	return nil
}
