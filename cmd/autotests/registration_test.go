package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/StasMerzlyakov/gophkeeper/internal/client/adapters/grpc/handler"
	"github.com/StasMerzlyakov/gophkeeper/internal/client/app"
	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/StasMerzlyakov/gophkeeper/internal/fork"
	"github.com/go-resty/resty/v2"
	"github.com/liyue201/goqr"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/suite"
)

type RegistrationTest struct {
	suite.Suite
	gophKeeperProcess *fork.BackgroundProcess
	serverPort        string
	client            app.AppServer
	smtpSrvHttpUrl    string
}

var _ suite.SetupAllSuite = (*RegistrationTest)(nil)
var _ suite.TearDownAllSuite = (*RegistrationTest)(nil)

func (suite *RegistrationTest) SetupSuite() {
	suite.T().Logf("Запускаем тест на проверку регистрации")

	suite.Require().NotEmpty(flagGophKeeperServerBinaryPath, "-gophkeeper-binary-path flag required")
	suite.Require().NotEmpty(flagGophKeeperTlsCaCert, "-gophkeeper-tls-ca-cert flag required")
	suite.Require().NotEmpty(flagGophKeeperTlsKey, "-gophkeeper-tls-key flag required")
	suite.Require().NotEmpty(flagGophKeeperTlsCert, "-gophkeeper-tls-cert flag required")
	suite.Require().NotEmpty(flagGophKeeperServerSecret, "-gophkeeper-server-secret non-empty flag required")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	suite.serverUp(ctx)
	suite.clientUp(ctx)
}

func (suite *RegistrationTest) clientDown(context.Context) {
	suite.client.Stop()
}

func (suite *RegistrationTest) clientUp(ctx context.Context) {
	serverAddr := fmt.Sprintf("localhost:%s", suite.serverPort)

	clntConf := &config.ClientConf{
		ServerAddress: serverAddr,
		CACert:        flagGophKeeperTlsCaCert,
	}

	var err error
	suite.client, err = handler.NewHandler(clntConf)
	suite.Require().NoErrorf(err, "Невозможно создать соединение сервером: %w", err)

	err = suite.client.Ping(ctx)
	suite.Require().NoErrorf(err, "Невозможно вызвать метод Ping на сервере: %w", err)
}

func (suite *RegistrationTest) serverUp(ctx context.Context) {
	envs := suite.getGophKeeperEnv(ctx)
	args := []string{} // оставлю на будущее
	suite.gophKeeperProcess = fork.NewBackgroundProcess(context.Background(),
		flagGophKeeperServerBinaryPath,
		fork.WithEnv(envs...),
		fork.WithArgs(args...),
	)
	err := suite.gophKeeperProcess.Start(ctx)
	suite.Require().NoErrorf(err, "Невозможно запустить процесс командой %q: %s. Переменные окружения: %+v, флаги командной строки: %+v", suite.gophKeeperProcess, err, envs, args)

	err = suite.gophKeeperProcess.WaitPort(ctx, "tcp", suite.serverPort)
	suite.Require().NoErrorf(err, "Не удалось дождаться пока порт %s станет доступен для запроса: %s\n%s", suite.serverPort, err, string(suite.gophKeeperProcess.Stdout(ctx)))

}

func (suite *RegistrationTest) getGophKeeperEnv(ctx context.Context) []string {

	var envs []string

	smtpPort, err := smtpServer.MappedPort(ctx, portSmtpServSMTP)
	suite.Require().NoError(err, "Не удалось получить порт SMTP сервера")

	smtpPortNumber := smtpPort.Port()
	suite.T().Logf("Для взаимодействия через SMTP получен порт контейнера %s", smtpPortNumber)

	smtpHostAddress := "localhost"

	smtpHttpPort, err := smtpServer.MappedPort(ctx, portSmtpServHTTP)
	suite.Require().NoError(err, "Не удалось получить порт SMTP-сервера для доступа по http")

	smtpHttpPortNumber := smtpHttpPort.Port()
	suite.T().Logf("Для получения сообщенй от SMTP по http получен порт контейнера %s", smtpHttpPortNumber)

	envs = append(envs, fmt.Sprintf("SMTP_HOST=%v", smtpHostAddress))
	envs = append(envs, fmt.Sprintf("SMTP_PORT=%v", smtpPortNumber))

	suite.smtpSrvHttpUrl = fmt.Sprintf("http://%s:%s", smtpHostAddress, smtpHttpPortNumber)

	serverEmail := "gookeeper@gookeeper.local"
	envs = append(envs, fmt.Sprintf("SERVER_EMAIL=%v", serverEmail))

	domainName := "gookeeper.local"
	envs = append(envs, fmt.Sprintf("DOMAIN_NAME=%v", domainName))

	databaseDN, err := postgresContainer.ConnectionString(context.Background())
	suite.Require().NoError(err)
	envs = append(envs, fmt.Sprintf("DATABASE_DN=%v", databaseDN))

	suite.serverPort = GetFreePort()
	serverAddr := fmt.Sprintf(":%s", suite.serverPort)
	envs = append(envs, fmt.Sprintf("PORT=%v", serverAddr))

	envs = append(envs, fmt.Sprintf("TLS_KEY=%v", flagGophKeeperTlsKey))
	envs = append(envs, fmt.Sprintf("TLS_CERT=%v", flagGophKeeperTlsCert))

	envs = append(envs, fmt.Sprintf("SERVER_SECRET=%v", flagGophKeeperServerSecret))
	return envs
}

func (suite *RegistrationTest) TearDownSuite() {

	if suite.gophKeeperProcess == nil {
		return
	}

	suite.clientDown(context.Background())

	exitCode, err := suite.gophKeeperProcess.Stop(syscall.SIGINT, syscall.SIGKILL)
	if err != nil {
		if errors.Is(err, os.ErrProcessDone) {
			return
		}
		suite.T().Logf("Не удалось остановить процесс с помощью сигнала ОС: %s", err)
		return
	}

	if exitCode > 0 {
		suite.T().Logf("Процесс завершился с не нулевым статусом %d", exitCode)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	out := suite.gophKeeperProcess.Stderr(ctx)
	if len(out) > 0 {
		suite.T().Logf("Получен STDERR лог процесса:\n\n%s", string(out))
	}
	out = suite.gophKeeperProcess.Stdout(ctx)
	if len(out) > 0 {
		suite.T().Logf("Получен STDOUT лог процесса:\n\n%s", string(out))
	}
}

func (suite *RegistrationTest) TestRegistration() {
	suite.Run("registration_login_success", func() {
		ctx := context.Background()

		userEmail := "tester@yandex.ru"
		userPassword := "testPasswordIKDe"
		//userMasterPassword := "testMasterKey!~?{asd}"

		// Проверка email
		status, err := suite.client.CheckEMail(ctx, userEmail)
		suite.Require().NoError(err, "Невозможно обратиться к сервису для проверки email: %w", err)
		suite.Require().Equalf(domain.EMailAvailable, status, "Статус email не соответствует ожидаемому")

		// Регистрация
		err = suite.client.Registrate(ctx, &domain.EMailData{
			EMail:    userEmail,
			Password: userPassword,
		})

		suite.Require().NoErrorf(err, "Ошибка при регистрации email: %w", err)

		// Провека email
		time.Sleep(3 * time.Second)

		keyURL := suite.exatractKeyFromEMail(ctx)

		otpKey, err := otp.NewKeyFromURL(keyURL)

		suite.Require().NoError(err, "Ошибка при восстановлении ключа")

		// Параметры валидации берем исключительно из восстановленного ключа
		validOpts := totp.ValidateOpts{
			Period:    uint(otpKey.Period()),
			Digits:    otpKey.Digits(),
			Algorithm: otpKey.Algorithm(),
		}

		regPass, err := totp.GenerateCodeCustom(otpKey.Secret(), time.Now(), validOpts)

		suite.Require().NoErrorf(err, "Не удалось сгенерировать otpPass")

		err = suite.client.PassRegOTP(ctx, regPass)
		suite.Require().NoError(err, "Ошибка при проверка ключа opt")

		masterPassword := "masterPassword"

		hellStr := domain.Random32ByteString()
		hellEcnrypted, err := domain.EncryptHello(masterPassword, hellStr)
		suite.Require().NoError(err, "Ошибка при кодировании на master-ключе")

		err = suite.client.InitMasterKey(ctx, &domain.MasterKeyData{
			MasterPasswordHint: "Hint",
			HelloEncrypted:     hellEcnrypted,
		})
		suite.Require().NoError(err, "Ошибка при инициализации master-ключа")

		suite.T().Log("Регистрация успешной пройдена")
		time.Sleep(3 * time.Second)

		// Вход в систему
		err = suite.client.Login(ctx, &domain.EMailData{
			EMail:    userEmail,
			Password: userPassword,
		})
		suite.Require().NoError(err, "Ошибка при логине в систему")

		logingOTP, err := totp.GenerateCodeCustom(otpKey.Secret(), time.Now(), validOpts)
		suite.Require().NoErrorf(err, "Не удалось сгенерировать otpPass для логина")

		err = suite.client.PassLoginOTP(ctx, logingOTP)
		suite.Require().NoErrorf(err, "Не удалось ввести OTP пароль для логина")

		suite.T().Log("Логин в систему произведен")

		hData, err := suite.client.GetHelloData(ctx)
		suite.Require().NoErrorf(err, "Не удалось получить данные для проверки master-ключа")

		err = domain.DecryptHello(masterPassword, hData.HelloEncrypted)
		suite.Require().NoErrorf(err, "Ошибка при проверке hello")

		suite.T().Log("Ключ master-ключ успешно проверен")
	})

}

func (suite *RegistrationTest) exatractKeyFromEMail(ctx context.Context) string {
	email := suite.getEMailContent(ctx)
	qr := suite.extractQR(email)
	keyURL := suite.decodeQRText(qr)
	return keyURL
}

func (suite *RegistrationTest) getEMailContent(ctx context.Context) []byte {
	client := resty.New()

	msgUrl := fmt.Sprintf("%s/messages/1.eml", suite.smtpSrvHttpUrl)
	resp, err := client.R().
		SetContext(ctx).
		Get(msgUrl)
	suite.Require().NoError(err, "Ошибка при обращении к smtp-сервису")

	suite.Require().Equalf(http.StatusOK, resp.StatusCode(), "Не удалось получить сообщение с SMPT-серева по адресу %s", msgUrl)
	return resp.Body()

}

func (suite *RegistrationTest) extractQR(emailContent []byte) []byte {

	// Ищем QR в сообщении. QR в base64 и расположен так:
	// .....
	// Content-Type: application/octet-stream;
	//	 	name="QR.png"
	//
	// iVBORw0KGgoAAAANSUhEUgAAAcIAAAHCEAAAAAAJ40
	// .....
	// --555777

	var qrB64 bytes.Buffer
	buf := bytes.NewBuffer(emailContent)
	scanner := bufio.NewScanner(buf)
	state := 0
Loop:
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			suite.Require().NoError(err, "Ошибка разборе email")
		}
		line := scanner.Text()
		switch state {
		case 0:
			if strings.HasPrefix(line, "Content-Type: application/octet-stream;") {
				state++
			}
		case 1:
			if strings.Contains(line, "name=\"QR.png\"") {
				state++
			} else {
				suite.Require().Failf("Ошибка при разборе email - ожидалась строка вида name=\"QR.png\"", "email: %s", string(emailContent))
			}
		case 2:
			if line == "" {
				state++
			} else {
				suite.Require().Failf("Ошибка при разборе email - ожидалась пустая строка", "email: %s", string(emailContent))
			}
		case 3:
			if strings.HasPrefix(line, "--") {
				break Loop
			} else {
				qrB64.WriteString(line)
			}
		}
	}

	qr, err := base64.StdEncoding.DecodeString(qrB64.String())
	suite.Require().NoError(err, "Ошибка декодировании QR - email: %s\n\n qr: %s", string(emailContent), string(qr))
	return qr
}

func (suite *RegistrationTest) decodeQRText(qr []byte) string {

	img, _, err := image.Decode(bytes.NewReader(qr))
	suite.Require().NoError(err, "Ошибка декодировании QR")

	qrCodes, err := goqr.Recognize(img)
	suite.Require().NoError(err, "Ошибка распознавания qr-кода")

	for _, qrCode := range qrCodes {
		return string(qrCode.Payload)
	}

	suite.Require().Fail("Не удалось излечь данные из QR")
	return ""
}
