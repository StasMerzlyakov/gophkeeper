package email_test

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/server/adapters/email"
	"github.com/stretchr/testify/require"
)

const TestDataDirectory = "../../../../testdata/"

func TestSendMail(t *testing.T) {

	hostAddress, portNumber := "127.0.0.1", mockServer.PortNumber()

	serverEmail := "gookeeper@localdomain.ru"
	clientEmail := "user1@yandex.ru"

	qrFile := filepath.Join(TestDataDirectory, "QR.png")

	fl, err := os.Open(qrFile)
	require.NoError(t, err)
	defer func() {
		err := fl.Close()
		require.NoError(t, err)
	}()

	qr, err := io.ReadAll(fl)
	require.NoError(t, err)

	conf := &config.ServerConf{
		SMTPHost:    hostAddress,
		SMTPPort:    portNumber,
		ServerEMail: serverEmail,
	}

	emailSender := email.NewSender(conf)
	defer func() {
		_ = emailSender.Close() // It is possible connection already closed on error
	}()

	ctx := context.Background()

	err = emailSender.Connect(ctx)
	require.NoError(t, err)

	err = emailSender.Send(ctx, clientEmail, qr)
	require.NoError(t, err)

	err = emailSender.Close()
	require.NoError(t, err)

	time.Sleep(2 * time.Second)

	msgs := mockServer.MessagesAndPurge()
	require.True(t, len(msgs) == 1)
}

func TestSendMailErr(t *testing.T) {

	hostAddress, portNumber := "127.0.0.1", mockServer.PortNumber()

	serverEmail := "gookeeper@localdomain"
	clientEmail := "st.merzlyakov@yandex.ru"

	qrFile := filepath.Join(TestDataDirectory, "QR.png")

	fl, err := os.Open(qrFile)
	require.NoError(t, err)

	defer func() {
		err := fl.Close()
		require.NoError(t, err)
	}()

	qr, err := io.ReadAll(fl)
	require.NoError(t, err)

	conf := &config.ServerConf{
		SMTPHost:    hostAddress,
		SMTPPort:    portNumber,
		ServerEMail: serverEmail,
	}

	emailSender := email.NewSender(conf)

	ctx := context.Background()

	err = emailSender.Connect(ctx)
	require.NoError(t, err)
	defer func() {
		err := emailSender.Close()
		require.NoError(t, err)
	}()

	err = emailSender.Send(ctx, clientEmail, qr)
	require.Error(t, err)
}
