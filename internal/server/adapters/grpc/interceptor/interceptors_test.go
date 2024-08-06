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
	"github.com/golang/protobuf/ptypes/empty"
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

func (tr *testRequestIDPinger) Ping(ctx context.Context, empty *empty.Empty) (*empty.Empty, error) {
	val := ctx.Value(domain.LoggerKey)
	if val == nil {
		return nil, errors.New("LoggerKey is not exists in context")
	}

	if _, ok := val.(domain.Logger); !ok {
		return nil, errors.New("LoggerKey has wrong type")
	}

	return nil, nil
}

func TestEncrichWithRequestIDUnaryInterceptor(t *testing.T) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	require.NoError(t, err)

	l, err := net.ListenTCP("tcp", addr)
	require.NoError(t, err)

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.EncrichWithRequestIDUnaryInterceptor(),
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

	_, err = pinger.Ping(ctx, nil)
	require.NoError(t, err)

	stopFn()
	wg.Wait()
}

type testErrPinger struct {
	proto.UnimplementedPingerServer
	err error
}

func (tr *testErrPinger) Ping(ctx context.Context, empty *empty.Empty) (*empty.Empty, error) {
	return nil, tr.err
}

type testIDPinger struct {
	proto.UnimplementedPingerServer
}

func (tr *testIDPinger) Ping(ctx context.Context, empty *empty.Empty) (*empty.Empty, error) {
	val := ctx.Value(domain.LoggerKey)
	if val != nil {
		return nil, errors.New("LoggerKey is not set")
	}

	return nil, nil
}

func TestErrorCodeInteceptor(t *testing.T) {

	t.Run("no_error", func(t *testing.T) {

		addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
		require.NoError(t, err)

		l, err := net.ListenTCP("tcp", addr)
		require.NoError(t, err)

		s := grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				interceptor.ErrorCodeUnaryInteceptor(),
			),
		)
		ctx, stopFn := context.WithCancel(context.Background())
		defer stopFn()

		proto.RegisterPingerServer(s, &testIDPinger{})

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

		_, err = pinger.Ping(ctx, nil)
		require.NoError(t, err)

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
				interceptor.ErrorCodeUnaryInteceptor(),
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

		_, err = pinger.Ping(ctx, nil)

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
				interceptor.JWTUnaryInterceptor([]byte(tokenSecret), []string{}),
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

		_, err = pinger.Ping(ctx, nil)

		require.NoError(t, err)
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
				interceptor.ErrorCodeUnaryInteceptor(),
				interceptor.JWTUnaryInterceptor([]byte(tokenSecret), []string{"proto.Pinger"}),
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

		_, err = pinger.Ping(ctx, nil)
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
				interceptor.ErrorCodeUnaryInteceptor(),
				interceptor.JWTUnaryInterceptor([]byte(tokenSecret), []string{"proto.Pinger"}),
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

		authCtx := metadata.AppendToOutgoingContext(ctx, domain.AuthorizationMetadataTokenName, "wrong_token")

		_, err = pinger.Ping(authCtx, nil)
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
				interceptor.ErrorCodeUnaryInteceptor(),
				interceptor.JWTUnaryInterceptor([]byte(tokenSecret), []string{"proto.Pinger"}),
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

		jwtTok, err := domain.CreateJWTToken([]byte(tokenSecret), 10*time.Second, 1)
		require.NoError(t, err)
		authCtx := metadata.AppendToOutgoingContext(ctx, domain.AuthorizationMetadataTokenName, string(jwtTok))

		resp, err := pinger.Ping(authCtx, nil)

		require.NoError(t, err)
		require.NotNil(t, resp)

		stopFn()
		wg.Wait()
	})
}

type testLoggerKeyAccessor struct {
	proto.UnimplementedFileAccessorServer
}

func (fa *testLoggerKeyAccessor) UploadFile(fs proto.FileAccessor_UploadFileServer) error {

	ctx := fs.Context()

	val := ctx.Value(domain.LoggerKey)
	if val != nil {
		return errors.New("LoggerKey is not set")
	}

	return nil
}

func TestEncrichWithRequestIDStreamInterceptor(t *testing.T) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	require.NoError(t, err)

	l, err := net.ListenTCP("tcp", addr)
	require.NoError(t, err)

	s := grpc.NewServer(
		grpc.ChainStreamInterceptor(
			interceptor.EncrichWithRequestIDStreamInterceptor(),
		),
	)
	ctx, stopFn := context.WithCancel(context.Background())
	defer stopFn()

	proto.RegisterFileAccessorServer(s, &testLoggerKeyAccessor{})

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

	fileClient := proto.NewFileAccessorClient(client)
	_, err = fileClient.UploadFile(ctx)
	require.NoError(t, err)

	stopFn()
	wg.Wait()
}

/*
Test commented for future learning.

Once the server returns an error from the method handler, the stream is closed and sends/receives on the stream in other goroutines will fail.


https://github.com/grpc/grpc-go/issues/2548
https://github.com/grpc/grpc-go/issues/2435


type testErrorCodeAccessor struct {
	proto.UnimplementedFileAccessorServer
}

func (fa *testErrorCodeAccessor) UploadFile(fs proto.FileAccessor_UploadFileServer) error {
	return fmt.Errorf("%w auth test", domain.ErrNotAuthorized)
}

func TestErrorCodeStreamInterceptor(t *testing.T) {

	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	require.NoError(t, err)

	l, err := net.ListenTCP("tcp", addr)
	require.NoError(t, err)

	s := grpc.NewServer(
		grpc.ChainStreamInterceptor(
			interceptor.ErrorCodeStreamInterceptor(),
		),
	)
	ctx, stopFn := context.WithCancel(context.Background())
	defer stopFn()

	proto.RegisterFileAccessorServer(s, &testErrorCodeAccessor{})

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

	fileClient := proto.NewFileAccessorClient(client)
	cln, err := fileClient.UploadFile(ctx)

	require.Error(t, err)

	// require.NotNil(t, cln)

	// err = cln.Send(&proto.UploadFileRequest{
	// 	 Name:        "file",
	//   SizeInBytes: 10,
	//   Data:        make([]byte, 10),
	// })

	e, ok := status.FromError(err)
	require.True(t, ok)

	require.Equal(t, gp.PermissionDenied, e.Code())

	stopFn()
	wg.Wait()
}
*/

type testJWTFileAccessor struct {
	userID domain.UserID

	proto.UnimplementedFileAccessorServer

	errChan chan error
	doneCh  chan struct{}
}

func (fa *testJWTFileAccessor) UploadFile(fs proto.FileAccessor_UploadFileServer) error {
	defer func() {
		fa.doneCh <- struct{}{}
	}()
	ctx := fs.Context()
	ctxId, err := domain.GetUserID(ctx)
	if err != nil {
		fa.errChan <- err
	} else {
		if ctxId != fa.userID {
			fa.errChan <- fmt.Errorf("not equals")
		}
	}
	return nil
}

func TestJWTStreamInterceptor(t *testing.T) {

	t.Run("jwt_ok", func(t *testing.T) {
		addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
		require.NoError(t, err)

		l, err := net.ListenTCP("tcp", addr)
		require.NoError(t, err)

		tokenSecret := []byte("tokenSecret")
		userID := domain.UserID(10)
		jwtTok, err := domain.CreateJWTToken(tokenSecret, 5*time.Second, userID)
		require.NoError(t, err)

		s := grpc.NewServer(
			grpc.ChainStreamInterceptor(
				interceptor.JWTStreamInterceptor(tokenSecret),
			),
		)
		ctx, stopFn := context.WithTimeout(context.Background(), 3*time.Second)
		defer stopFn()

		doneCh := make(chan struct{})
		errChan := make(chan error)

		proto.RegisterFileAccessorServer(s, &testJWTFileAccessor{
			userID:  userID,
			errChan: errChan,
			doneCh:  doneCh,
		})

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

		fileClient := proto.NewFileAccessorClient(client)

		authCtx := metadata.AppendToOutgoingContext(ctx, domain.AuthorizationMetadataTokenName, string(jwtTok))

		cln, _ := fileClient.UploadFile(authCtx) // no error check see TestErrorCodeStreamInterceptor

		require.NotNil(t, cln)

		_ = cln.Send(&proto.UploadFileRequest{
			Name:        "file",
			SizeInBytes: 10,
			Data:        make([]byte, 10),
		})

		select {
		case err := <-errChan:
			require.NoError(t, err)
		case <-doneCh:
		case <-ctx.Done():
			require.NoError(t, ctx.Err())
		}
		stopFn()
		wg.Wait()
	})
}
