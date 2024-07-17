package email

import (
	"context"
	"encoding/base64"
	"fmt"
	"sync"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	mail "github.com/xhit/go-simple-mail"
)

func NewSender(conf *config.ServerConf) *sender {

	server := mail.NewSMTPClient()
	server.Host = conf.SMTPHost
	server.Port = conf.SMTPPort

	return &sender{
		conf:   conf,
		server: server,
	}
}

type sender struct {
	client     *mail.SMTPClient
	server     *mail.SMTPServer
	conf       *config.ServerConf
	once       sync.Once
	connectErr error
}

func (snd *sender) Connect(ctx context.Context) error {
	snd.once.Do(func() {
		snd.client, snd.connectErr = snd.server.Connect()
	})
	return snd.connectErr
}

func (snd *sender) Close() {
	snd.client.Close()
}

const htmlBody = `<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
	</head>
	<body>
		<p>Save and delete</p>
		<p><img src="cid:QR.png" alt="Go gopher"/></p>
		<p><b>Do not show anybody!!</b></p>
	</body>
</html>`

func (snd *sender) Send(ctx context.Context, email string, png []byte) error {
	emailMsg := mail.NewMSG()
	emailMsg.SetFrom(fmt.Sprintf("GophKeeper <%s>", snd.conf.SMTPServerEMail)).
		AddTo(email).
		SetSubject("GophKeeper OTP QR")
	emailMsg.SetBody(mail.TextHTML, htmlBody)
	pngB64 := base64.StdEncoding.EncodeToString(png)
	emailMsg.AddAttachmentBase64(pngB64, "QR.png")
	if emailMsg.Error != nil {
		return fmt.Errorf("%w: email creation error: %s", domain.ErrInternalServer, emailMsg.Error.Error())
	}

	if err := emailMsg.Send(snd.client); err != nil {
		return fmt.Errorf("%w: send error: %s", domain.ErrInternalServer, err.Error())
	}

	return nil
}