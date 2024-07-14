package handler

import (
	"context"

	"github.com/StasMerzlyakov/gophkeeper/internal/proto"
)

type pinger struct {
	proto.UnimplementedPingerServer
}

func (pn *pinger) Ping(context.Context, *proto.PingRequest) (*proto.PingResponse, error) {
	var response proto.PingResponse
	return &response, nil
}
