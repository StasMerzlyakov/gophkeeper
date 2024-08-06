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
			err := fmt.Errorf("%w - %v error - %v - can't encrypt card with number %v", domain.ErrClientDataIncorrect, action, err.Error(), encr.Number)
			log.Warn(err.Error())
			return err
		}

		var bankCard domain.BankCard
		if err := json.Unmarshal([]byte(encrypted), &bankCard); err != nil {
			err := fmt.Errorf("%w - %v error - %v - can't decode card data with number %v", domain.ErrClientDataIncorrect, action, err.Error(), encr.Number)
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

	if err := dcc.helper.CheckBankCardData(bankCard); err != nil {
		err := fmt.Errorf("%w - %v error - wrong card data %v", domain.ErrClientDataIncorrect, action, err.Error())
		log.Warn(err.Error())
		return err
	}

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

	if err := dcc.helper.CheckBankCardData(bankCard); err != nil {
		err := fmt.Errorf("%w - %v error - wrong card data %v", domain.ErrClientDataIncorrect, action, err.Error())
		log.Warn(err.Error())
		return err
	}

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

func (dcc *dataAccessor) GetUserPasswordDataList(ctx context.Context) error {
	log := GetMainLogger()
	action := domain.GetAction(1)
	log.Debugf("%v start", action)

	masterPass := dcc.appStorage.GetMasterPassword()

	encrList, err := dcc.appServer.GetUserPasswordDataList(ctx)
	if err != nil {
		err := fmt.Errorf("%w - %v error", err, action)
		log.Warn(err.Error())
		return err
	}

	var decryptedList []domain.UserPasswordData
	for _, encr := range encrList {
		content := encr.Content

		encrypted, err := dcc.helper.DecryptShortData(masterPass, content)
		if err != nil {
			err := fmt.Errorf("%w - %v error - %s - can't encrypt userPassData with hint %v", domain.ErrClientDataIncorrect, action, err.Error(), encr.Hint)
			log.Warn(err.Error())
			return err
		}

		var uPassData domain.UserPasswordData
		if err := json.Unmarshal([]byte(encrypted), &uPassData); err != nil {
			err := fmt.Errorf("%w - %v error - %s - can't decode userPassData with hint %v", domain.ErrClientDataIncorrect, action, err.Error(), encr.Hint)
			log.Warn(err.Error())
			return err
		}
		decryptedList = append(decryptedList, uPassData)
	}

	dcc.appStorage.SetUserPasswordDatas(decryptedList)

	log.Debugf("%v success", action)
	return nil
}

func (dcc *dataAccessor) AddUserPasswordData(ctx context.Context, data *domain.UserPasswordData) error {
	log := GetMainLogger()
	action := domain.GetAction(1)
	log.Debugf("%v start", action)

	if err := dcc.helper.CheckUserPasswordData(data); err != nil {
		err := fmt.Errorf("%w - %v error - wrong userPassword data %v", domain.ErrClientDataIncorrect, action, err.Error())
		log.Warn(err.Error())
		return err
	}

	res, err := json.Marshal(data)
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

	err = dcc.appServer.CreateUserPasswordData(ctx, &domain.EncryptedUserPasswordData{
		Hint:    data.Hint,
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

func (dcc *dataAccessor) UpdateUserPasswordData(ctx context.Context, data *domain.UserPasswordData) error {
	log := GetMainLogger()
	action := domain.GetAction(1)
	log.Debugf("%v start", action)

	if err := dcc.helper.CheckUserPasswordData(data); err != nil {
		err := fmt.Errorf("%w - %v error - wrong userPassword data %v", domain.ErrClientDataIncorrect, action, err.Error())
		log.Warn(err.Error())
		return err
	}

	res, err := json.Marshal(data)
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

	err = dcc.appServer.UpdateUserPasswordData(ctx, &domain.EncryptedUserPasswordData{
		Hint:    data.Hint,
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

func (dcc *dataAccessor) DeleteUserPasswordData(ctx context.Context, hint string) error {
	log := GetMainLogger()
	action := domain.GetAction(1)
	log.Debugf("%v start", action)

	if err := dcc.appServer.DeleteUserPasswordData(ctx, hint); err != nil {
		err := fmt.Errorf("%w - %v error", err, action)
		log.Warn(err.Error())
		return err
	}

	log.Debugf("%v success", action)
	return nil
}
