package handler_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/StasMerzlyakov/gophkeeper/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var srv *grpc.Server
var srvPort int

var dtAccessor = &dataAccessor{}
var rgHandler = &regHandler{}
var athService = &authService{}

var wg sync.WaitGroup

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		panic(err)
	}
	code := m.Run()
	shutdown()
	os.Exit(code)
}

const TestDataDirectory = "../../../../../testdata/"

func shutdown() {
	srv.GracefulStop()
	wg.Wait()
}

func setup() error {
	keyFile := filepath.Join(TestDataDirectory, "test-server-key.pem")
	certFile := filepath.Join(TestDataDirectory, "test-server-cert.pem")

	var err error
	srvPort, err = getFreePort()
	if err != nil {
		return err
	}

	serverCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	tlsCredentials := credentials.NewTLS(config)

	srv = grpc.NewServer(
		grpc.Creds(tlsCredentials),
	)

	proto.RegisterRegistrationServiceServer(srv, rgHandler)
	proto.RegisterDataAccessorServer(srv, dtAccessor)
	proto.RegisterAuthServiceServer(srv, athService)

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", srvPort))
	if err != nil {
		return err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.Serve(listen); err != nil {
			panic(err)
		}
	}()
	return nil
}

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

type dataAccessor struct {
	proto.UnimplementedDataAccessorServer
	helloFn func(ctx context.Context, req *proto.HelloRequest) (*proto.HelloResponse, error)
}

func (da *dataAccessor) Hello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloResponse, error) {
	return da.helloFn(ctx, req)
}

type regHandler struct {
	proto.UnimplementedRegistrationServiceServer
	checkEMailFn   func(ctx context.Context, req *proto.CheckEMailRequest) (*proto.CheckEMailResponse, error)
	registrateFn   func(ctx context.Context, req *proto.RegistrationRequest) (*proto.RegistrationResponse, error)
	passOTPFn      func(ctx context.Context, req *proto.PassOTPRequest) (*proto.PassOTPResponse, error)
	setMasterKeyFn func(ctx context.Context, req *proto.MasterKeyRequest) (*proto.MasterKeyResponse, error)
}

func (rh *regHandler) CheckEMail(ctx context.Context, req *proto.CheckEMailRequest) (*proto.CheckEMailResponse, error) {
	return rh.checkEMailFn(ctx, req)
}

func (rh *regHandler) Registrate(ctx context.Context, req *proto.RegistrationRequest) (*proto.RegistrationResponse, error) {
	return rh.registrateFn(ctx, req)
}

func (rh *regHandler) PassOTP(ctx context.Context, req *proto.PassOTPRequest) (*proto.PassOTPResponse, error) {
	return rh.passOTPFn(ctx, req)
}

func (rh *regHandler) SetMasterKey(ctx context.Context, req *proto.MasterKeyRequest) (*proto.MasterKeyResponse, error) {
	return rh.setMasterKeyFn(ctx, req)
}

type authService struct {
	proto.UnimplementedAuthServiceServer
	loginFn   func(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error)
	passOTPFn func(ctx context.Context, req *proto.PassOTPRequest) (*proto.AuthResponse, error)
}

func (aS *authService) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	return aS.loginFn(ctx, req)
}

func (aS *authService) PassOTP(ctx context.Context, req *proto.PassOTPRequest) (*proto.AuthResponse, error) {
	return aS.passOTPFn(ctx, req)
}
