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

func (gh *grpcHandler) RegHandler(regHandler *regHandler) *grpcHandler {
	gh.regHandler = regHandler
	return gh
}

func (gh *grpcHandler) DataAccessor(dataAccessor *dataAccessor) *grpcHandler {
	gh.dataAccessor = dataAccessor
	return gh
}

func (gh *grpcHandler) AuthService(authService *authService) *grpcHandler {
	gh.authService = authService
	return gh
}

type grpcHandler struct {
	conf         *config.ServerConf
	s            *grpc.Server
	regHandler   *regHandler
	dataAccessor *dataAccessor
	authService  *authService
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
	proto.RegisterRegistrationServiceServer(grpcHandler.s, grpcHandler.regHandler)
	proto.RegisterDataAccessorServer(grpcHandler.s, grpcHandler.dataAccessor)
	proto.RegisterAuthServiceServer(grpcHandler.s, grpcHandler.authService)
	return grpcHandler.s.Serve(listen)
}

func (grpcHandler *grpcHandler) Stop() {
	grpcHandler.s.GracefulStop()
}
