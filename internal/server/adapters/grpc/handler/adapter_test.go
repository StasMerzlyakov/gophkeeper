package handler_test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/proto"
	"github.com/StasMerzlyakov/gophkeeper/internal/server/adapters/grpc/handler"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const TestDataDirectory = "../../../../../testdata/"

func TestPingNoTls(t *testing.T) {
	port, err := getFreePort()
	require.NoError(t, err)

	srvConf := &config.ServerConf{
		Port: fmt.Sprintf(":%d", port),
	}

	srv := handler.NewGRPCHandler(srvConf)

	ctx, stopFn := context.WithCancel(context.Background())
	defer stopFn()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := srv.Start(ctx)
		require.NoError(t, err)
	}()

	time.Sleep(2 * time.Second)

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		srv.Stop()
	}()

	client, err := grpc.NewClient(fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	pinger := proto.NewPingerClient(client)

	pingReq := proto.PingRequest{}
	resp, err := pinger.Ping(ctx, &pingReq)
	require.NoError(t, err)
	require.NotNil(t, resp)

	stopFn()
	wg.Wait()
}

func TestPingWithTls(t *testing.T) {
	port, err := getFreePort()
	require.NoError(t, err)

	keyFile := filepath.Join(TestDataDirectory, "test-server-key.pem")
	certFile := filepath.Join(TestDataDirectory, "test-server-cert.pem")

	srvConf := &config.ServerConf{
		Port:    fmt.Sprintf(":%d", port),
		TLSKey:  keyFile,
		TLSCert: certFile,
	}

	srv := handler.NewGRPCHandler(srvConf)

	ctx, stopFn := context.WithCancel(context.Background())
	defer stopFn()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := srv.Start(ctx)
		require.NoError(t, err)
	}()

	time.Sleep(2 * time.Second)

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		srv.Stop()
	}()

	caFile := filepath.Join(TestDataDirectory, "test-ca-cert.pem")

	cred, err := loadTLSCredentials(caFile)
	require.NoError(t, err)

	client, err := grpc.NewClient(fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(cred))
	require.NoError(t, err)

	pinger := proto.NewPingerClient(client)

	pingReq := proto.PingRequest{}
	resp, err := pinger.Ping(ctx, &pingReq)
	require.NoError(t, err)
	require.NotNil(t, resp)

	stopFn()
	wg.Wait()
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

func loadTLSCredentials(caFile string) (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate
	file, err := os.Open(caFile)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	pemServerCA, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// Create the credentials and return it
	config := &tls.Config{
		RootCAs: certPool,
	}

	return credentials.NewTLS(config), nil
}
