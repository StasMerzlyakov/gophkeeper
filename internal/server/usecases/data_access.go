package usecases

import (
	"context"
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
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

func (da *dataAccessor) GetHelloData(ctx context.Context) (string, error) {
	res, err := da.stflStorage.GetHelloData(ctx)
	if err != nil {
		return res, fmt.Errorf("getHelloData err %w", err)
	}
	return res, err
}
