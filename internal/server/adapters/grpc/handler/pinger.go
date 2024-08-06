package handler

import (
	"context"

	"github.com/StasMerzlyakov/gophkeeper/internal/proto"
	"github.com/golang/protobuf/ptypes/empty"
)

type pinger struct {
	proto.UnimplementedPingerServer
}

func (pn *pinger) Ping(context.Context, *empty.Empty) (*empty.Empty, error) {
	return nil, nil
}
