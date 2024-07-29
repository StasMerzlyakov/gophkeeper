package handler

import (
	"context"
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/StasMerzlyakov/gophkeeper/internal/proto"
)

func NewDataAccessor(accessor DataAccessor) *dataAccessor {
	return &dataAccessor{
		accessor: accessor,
	}
}

type dataAccessor struct {
	proto.UnimplementedDataAccessorServer
	accessor DataAccessor
}

func (da *dataAccessor) Hello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloResponse, error) {
	action := domain.GetAction(0)

	dd, err := da.accessor.GetHelloData(ctx)
	if err != nil {
		return nil, fmt.Errorf("%v err - %w", action, err)
	}

	resp := &proto.HelloResponse{
		HelloEncrypted:     dd.HelloEncrypted,
		EncryptedMasterKey: dd.EncryptedMasterKey,
	}

	return resp, nil
}
