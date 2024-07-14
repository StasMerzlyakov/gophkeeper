package handler

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func NewGRPCHandler(conf *config.ServerConf) *grpcHandler {
	return &grpcHandler{
		conf: conf,
	}
}

type grpcHandler struct {
	conf *config.ServerConf
	s    *grpc.Server
}

func (grpcHandler *grpcHandler) loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair(grpcHandler.conf.TLSCert, grpcHandler.conf.TLSKey)
	if err != nil {
		return nil, fmt.Errorf("can't loadTLSCrdentials %w", err)
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}

func (grpcHandler *grpcHandler) Start(srcCtx context.Context) error {
	listen, err := net.Listen("tcp", grpcHandler.conf.Port)
	if err != nil {
		return fmt.Errorf("can't listen %w", err)
	}
	if grpcHandler.conf.TLSKey != "" {
		tlsCredentials, err := grpcHandler.loadTLSCredentials()
		if err != nil {
			return err
		}
		grpcHandler.s = grpc.NewServer(
			grpc.Creds(tlsCredentials),
		)
	} else {
		grpcHandler.s = grpc.NewServer()
	}

	proto.RegisterPingerServer(grpcHandler.s, &pinger{})
	return grpcHandler.s.Serve(listen)
}

func (grpcHandler *grpcHandler) Stop() {
	grpcHandler.s.GracefulStop()
}
