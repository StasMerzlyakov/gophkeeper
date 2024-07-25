package handler

import (
	"context"

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
	dd, err := da.accessor.GetHelloData(ctx)
	if err != nil {
		return nil, WrapErr(err)
	}

	resp := &proto.HelloResponse{
		HelloEncrypted:     dd.HelloEncrypted,
		EncryptedMasterKey: dd.EncryptedMasterKey,
	}

	return resp, nil
}
