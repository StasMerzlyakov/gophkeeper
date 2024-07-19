package usecases

import (
	"context"
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
)

func NewDataAccessor(conf *config.ServerConf, stflStorage StateFullStorage) *dataAccessor {
	return &dataAccessor{
		conf:        conf,
		stflStorage: stflStorage,
	}
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
