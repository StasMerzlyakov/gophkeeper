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
