package interceptor_test

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/StasMerzlyakov/gophkeeper/internal/proto"
	"github.com/StasMerzlyakov/gophkeeper/internal/server/adapters/grpc/interceptor"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	gp "google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type testRequestIDPinger struct {
	proto.UnimplementedPingerServer
}

func (tr *testRequestIDPinger) Ping(ctx context.Context, req *proto.PingRequest) (*proto.PingResponse, error) {
	val := ctx.Value(domain.LoggerKey)
	if val == nil {
		return nil, errors.New("LoggerKey is not exists in context")
	}

	if _, ok := val.(domain.Logger); !ok {
		return nil, errors.New("LoggerKey has wrong type")
	}

	return &proto.PingResponse{}, nil
}

func TestEncrichWithRequestIDInterceptor(t *testing.T) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	require.NoError(t, err)

	l, err := net.ListenTCP("tcp", addr)
	require.NoError(t, err)

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.EncrichWithRequestIDInterceptor(),
		),
	)
	ctx, stopFn := context.WithCancel(context.Background())
	defer stopFn()

	proto.RegisterPingerServer(s, &testRequestIDPinger{})

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = s.Serve(l)
		require.NoError(t, err)

	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		s.GracefulStop()
	}()

	// client
	port := l.Addr().(*net.TCPAddr).Port
	client, err := grpc.NewClient(fmt.Sprintf("localhost:%d", port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	pinger := proto.NewPingerClient(client)

	pingReq := proto.PingRequest{}

	resp, err := pinger.Ping(ctx, &pingReq)
	require.NoError(t, err)
	require.NotNil(t, resp)

	stopFn()
	wg.Wait()
}

type testErrPinger struct {
	proto.UnimplementedPingerServer
	err error
}

func (tr *testErrPinger) Ping(ctx context.Context, req *proto.PingRequest) (*proto.PingResponse, error) {
	return &proto.PingResponse{}, tr.err
}

func TestErrorCodeInteceptor(t *testing.T) {

	t.Run("no_error", func(t *testing.T) {

		addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
		require.NoError(t, err)

		l, err := net.ListenTCP("tcp", addr)
		require.NoError(t, err)

		s := grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				interceptor.ErrorCodeInteceptor(),
			),
		)
		ctx, stopFn := context.WithCancel(context.Background())
		defer stopFn()

		proto.RegisterPingerServer(s, &testErrPinger{})

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = s.Serve(l)
			require.NoError(t, err)

		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ctx.Done()
			s.GracefulStop()
		}()

		// client
		port := l.Addr().(*net.TCPAddr).Port
		client, err := grpc.NewClient(fmt.Sprintf("localhost:%d", port),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err)

		pinger := proto.NewPingerClient(client)

		pingReq := proto.PingRequest{}

		resp, err := pinger.Ping(ctx, &pingReq)
		require.NoError(t, err)
		require.NotNil(t, resp)

		stopFn()
		wg.Wait()
	})

	t.Run("error", func(t *testing.T) {

		addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
		require.NoError(t, err)

		l, err := net.ListenTCP("tcp", addr)
		require.NoError(t, err)

		s := grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				interceptor.ErrorCodeInteceptor(),
			),
		)
		ctx, stopFn := context.WithCancel(context.Background())
		defer stopFn()

		proto.RegisterPingerServer(s, &testErrPinger{err: fmt.Errorf("%w err", domain.ErrAuthDataIncorrect)})

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = s.Serve(l)
			require.NoError(t, err)

		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ctx.Done()
			s.GracefulStop()
		}()

		// client
		port := l.Addr().(*net.TCPAddr).Port
		client, err := grpc.NewClient(fmt.Sprintf("localhost:%d", port),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err)

		pinger := proto.NewPingerClient(client)

		pingReq := proto.PingRequest{}

		_, err = pinger.Ping(ctx, &pingReq)

		e, ok := status.FromError(err)
		require.True(t, ok)

		require.Equal(t, gp.InvalidArgument, e.Code())
		stopFn()
		wg.Wait()
	})
}

func TestJWTInterceptor(t *testing.T) {
	t.Run("no_auth_need", func(t *testing.T) {
		addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
		require.NoError(t, err)

		l, err := net.ListenTCP("tcp", addr)
		require.NoError(t, err)

		tokenSecret := "tokenSecret"

		s := grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				interceptor.JWTInterceptor([]byte(tokenSecret), []string{}),
			),
		)
		ctx, stopFn := context.WithCancel(context.Background())
		defer stopFn()

		proto.RegisterPingerServer(s, &testErrPinger{err: nil})

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = s.Serve(l)
			require.NoError(t, err)

		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ctx.Done()
			s.GracefulStop()
		}()

		// client
		port := l.Addr().(*net.TCPAddr).Port
		client, err := grpc.NewClient(fmt.Sprintf("localhost:%d", port),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err)

		pinger := proto.NewPingerClient(client)

		pingReq := proto.PingRequest{}

		resp, err := pinger.Ping(ctx, &pingReq)

		require.NoError(t, err)
		require.NotNil(t, resp)
		stopFn()
		wg.Wait()
	})

	t.Run("auth_need", func(t *testing.T) {
		addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
		require.NoError(t, err)

		l, err := net.ListenTCP("tcp", addr)
		require.NoError(t, err)

		tokenSecret := "tokenSecret"

		s := grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				interceptor.ErrorCodeInteceptor(),
				interceptor.JWTInterceptor([]byte(tokenSecret), []string{"proto.Pinger"}),
			),
		)
		ctx, stopFn := context.WithCancel(context.Background())
		defer stopFn()

		proto.RegisterPingerServer(s, &testErrPinger{err: nil})

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = s.Serve(l)
			require.NoError(t, err)

		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ctx.Done()
			s.GracefulStop()
		}()

		// client
		port := l.Addr().(*net.TCPAddr).Port
		client, err := grpc.NewClient(fmt.Sprintf("localhost:%d", port),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err)

		pinger := proto.NewPingerClient(client)

		pingReq := proto.PingRequest{}

		_, err = pinger.Ping(ctx, &pingReq)
		e, ok := status.FromError(err)
		require.True(t, ok)

		require.Equal(t, gp.PermissionDenied, e.Code())

		stopFn()
		wg.Wait()
	})

	t.Run("wrong_jwt", func(t *testing.T) {
		addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
		require.NoError(t, err)

		l, err := net.ListenTCP("tcp", addr)
		require.NoError(t, err)

		tokenSecret := "tokenSecret"

		s := grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				interceptor.ErrorCodeInteceptor(),
				interceptor.JWTInterceptor([]byte(tokenSecret), []string{"proto.Pinger"}),
			),
		)
		ctx, stopFn := context.WithCancel(context.Background())
		defer stopFn()

		proto.RegisterPingerServer(s, &testErrPinger{err: nil})

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = s.Serve(l)
			require.NoError(t, err)

		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ctx.Done()
			s.GracefulStop()
		}()

		// client
		port := l.Addr().(*net.TCPAddr).Port
		client, err := grpc.NewClient(fmt.Sprintf("localhost:%d", port),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err)

		pinger := proto.NewPingerClient(client)

		pingReq := proto.PingRequest{}

		authCtx := metadata.AppendToOutgoingContext(ctx, domain.AuthorizationMetadataTokenName, "wrong_token")

		_, err = pinger.Ping(authCtx, &pingReq)
		e, ok := status.FromError(err)
		require.True(t, ok)

		require.Equal(t, gp.PermissionDenied, e.Code())

		stopFn()
		wg.Wait()
	})

	t.Run("jwt_ok", func(t *testing.T) {
		addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
		require.NoError(t, err)

		l, err := net.ListenTCP("tcp", addr)
		require.NoError(t, err)

		tokenSecret := "tokenSecret"

		s := grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				interceptor.ErrorCodeInteceptor(),
				interceptor.JWTInterceptor([]byte(tokenSecret), []string{"proto.Pinger"}),
			),
		)
		ctx, stopFn := context.WithCancel(context.Background())
		defer stopFn()

		proto.RegisterPingerServer(s, &testErrPinger{err: nil})

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = s.Serve(l)
			require.NoError(t, err)

		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ctx.Done()
			s.GracefulStop()
		}()

		// client
		port := l.Addr().(*net.TCPAddr).Port
		client, err := grpc.NewClient(fmt.Sprintf("localhost:%d", port),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err)

		pinger := proto.NewPingerClient(client)

		pingReq := proto.PingRequest{}

		jwtTok, err := domain.CreateJWTToken([]byte(tokenSecret), 10*time.Second, 1)
		require.NoError(t, err)
		authCtx := metadata.AppendToOutgoingContext(ctx, domain.AuthorizationMetadataTokenName, string(jwtTok))

		resp, err := pinger.Ping(authCtx, &pingReq)

		require.NoError(t, err)
		require.NotNil(t, resp)

		stopFn()
		wg.Wait()
	})
}
