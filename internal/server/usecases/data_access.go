package usecases

import (
	"context"
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

func NewDataAccessor(conf *config.ServerConf) *dataAccessor {
	return &dataAccessor{
		conf: conf,
	}
}

func (dAcc *dataAccessor) StateFullStorage(stflStorage StateFullStorage) *dataAccessor {
	dAcc.stflStorage = stflStorage
	return dAcc
}

type dataAccessor struct {
	conf        *config.ServerConf
	stflStorage StateFullStorage
}

func (da *dataAccessor) GetHelloData(ctx context.Context) (*domain.HelloData, error) {
	log := domain.GetCtxLogger(ctx)
	action := domain.GetAction(1)

	log.Debugw(action, "msg", "getHelloData start")

	res, err := da.stflStorage.GetHelloData(ctx)
	if err != nil {
		err := fmt.Errorf("getHelloData err %w", err)
		log.Infow(action, "err", err.Error())
		return res, fmt.Errorf("getHelloData err %w", err)
	}

	log.Debugw(action, "msg", "getHelloData success")
	return res, nil
}

func (da *dataAccessor) GetBankCardList(ctx context.Context) ([]domain.EncryptedBankCard, error) {
	log := domain.GetCtxLogger(ctx)
	action := domain.GetAction(1)

	log.Debugw(action, "msg", fmt.Sprintf("%s start", action))
	data, err := da.stflStorage.GetBankCardList(ctx)
	if err != nil {
		err := fmt.Errorf("%s err %w", action, err)
		log.Infow(action, "err", err.Error())
		return nil, err
	}

	log.Debugw(action, "msg", fmt.Sprintf("%s success", action))
	return data, nil
}

func (da *dataAccessor) CreateBankCard(ctx context.Context, bnkCard *domain.EncryptedBankCard) error {
	log := domain.GetCtxLogger(ctx)
	action := domain.GetAction(1)

	log.Debugw(action, "msg", fmt.Sprintf("%s start", action))
	if err := da.stflStorage.CreateBankCard(ctx, bnkCard); err != nil {
		err := fmt.Errorf("%s err %w", action, err)
		log.Infow(action, "err", err.Error())
		return err
	}

	log.Debugw(action, "msg", fmt.Sprintf("%s success", action))
	return nil
}

func (da *dataAccessor) UpdateBankCard(ctx context.Context, bnkCard *domain.EncryptedBankCard) error {
	log := domain.GetCtxLogger(ctx)
	action := domain.GetAction(1)

	log.Debugw(action, "msg", fmt.Sprintf("%s start", action))
	if err := da.stflStorage.UpdateBankCard(ctx, bnkCard); err != nil {
		err := fmt.Errorf("%s err %w", action, err)
		log.Infow(action, "err", err.Error())
		return err
	}

	log.Debugw(action, "msg", fmt.Sprintf("%s success", action))
	return nil
}

func (da *dataAccessor) DeleteBankCard(ctx context.Context, number string) error {
	log := domain.GetCtxLogger(ctx)
	action := domain.GetAction(1)

	log.Debugw(action, "msg", fmt.Sprintf("%s start", action))
	if err := da.stflStorage.DeleteBankCard(ctx, number); err != nil {
		err := fmt.Errorf("%s err %w", action, err)
		log.Infow(action, "err", err.Error())
		return err
	}

	log.Debugw(action, "msg", fmt.Sprintf("%s success", action))
	return nil
}

func (da *dataAccessor) GetUserPasswordDataList(ctx context.Context) ([]domain.EncryptedUserPasswordData, error) {
	log := domain.GetCtxLogger(ctx)
	action := domain.GetAction(1)

	log.Debugw(action, "msg", fmt.Sprintf("%s start", action))
	data, err := da.stflStorage.GetUserPasswordDataList(ctx)
	if err != nil {
		err := fmt.Errorf("%s err %w", action, err)
		log.Infow(action, "err", err.Error())
		return nil, err
	}

	log.Debugw(action, "msg", fmt.Sprintf("%s success", action))
	return data, nil
}

func (da *dataAccessor) CreateUserPasswordData(ctx context.Context, bnkCard *domain.EncryptedUserPasswordData) error {
	log := domain.GetCtxLogger(ctx)
	action := domain.GetAction(1)

	log.Debugw(action, "msg", fmt.Sprintf("%s start", action))
	if err := da.stflStorage.CreateUserPasswordData(ctx, bnkCard); err != nil {
		err := fmt.Errorf("%s err %w", action, err)
		log.Infow(action, "err", err.Error())
		return err
	}

	log.Debugw(action, "msg", fmt.Sprintf("%s success", action))
	return nil
}

func (da *dataAccessor) UpdateUserPasswordData(ctx context.Context, bnkCard *domain.EncryptedUserPasswordData) error {
	log := domain.GetCtxLogger(ctx)
	action := domain.GetAction(1)

	log.Debugw(action, "msg", fmt.Sprintf("%s start", action))
	if err := da.stflStorage.UpdateUserPasswordData(ctx, bnkCard); err != nil {
		err := fmt.Errorf("%s err %w", action, err)
		log.Infow(action, "err", err.Error())
		return err
	}

	log.Debugw(action, "msg", fmt.Sprintf("%s success", action))
	return nil
}

func (da *dataAccessor) DeleteUserPasswordData(ctx context.Context, hint string) error {
	log := domain.GetCtxLogger(ctx)
	action := domain.GetAction(1)

	log.Debugw(action, "msg", fmt.Sprintf("%s start", action))
	if err := da.stflStorage.DeleteUserPasswordData(ctx, hint); err != nil {
		err := fmt.Errorf("%s err %w", action, err)
		log.Infow(action, "err", err.Error())
		return err
	}

	log.Debugw(action, "msg", fmt.Sprintf("%s success", action))
	return nil
}
