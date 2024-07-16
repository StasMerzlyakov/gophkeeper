package email_test

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
	mail "github.com/xhit/go-simple-mail"
)

const htmlBody = `<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
		<title>Hello Gophers!</title>
	</head>
	<body>
		<p>This is the <b>Go gopher</b>.</p>
		<p><img src="cid:Gopher.png" alt="Go gopher" /></p>
		<p>Image created by Renee French</p>
	</body>
</html>`

const htmlBodySimple = `<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
		<title>Hello Gophers!</title>
	</head>
	<body>
		<p>This is the <b>Go gopher</b>.</p>
		<p>Image created by Renee French</p>
	</body>
</html>`

func TestSendMail(t *testing.T) {
	server := mail.NewSMTPClient()

	server.Host = "127.0.0.1"
	server.Port = 25

	smtpClient, err := server.Connect()
	require.NoError(t, err)

	email := mail.NewMSG()
	email.SetFrom("From Example <test@test.com>").
		AddTo("st.merzlyakov@yandex.ru").
		SetSubject("New Go Email")

	email.SetBody(mail.TextHTML, htmlBodySimple)

	if email.Error != nil {
		log.Fatal(email.Error)
	}

	err = email.Send(smtpClient)
	require.NoError(t, err)
}
