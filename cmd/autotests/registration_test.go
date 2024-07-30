package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/StasMerzlyakov/gophkeeper/internal/client/adapters/grpc/handler"
	"github.com/StasMerzlyakov/gophkeeper/internal/client/app"
	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/StasMerzlyakov/gophkeeper/internal/fork"
	"github.com/stretchr/testify/suite"
)

type RegistrationTest struct {
	suite.Suite
	gophKeeperProcess *fork.BackgroundProcess
	serverPort        string
	client            app.AppServer
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
	envs := suite.getGophKeeperEnv()
	args := []string{} // оставлю на будущее
	suite.gophKeeperProcess = fork.NewBackgroundProcess(context.Background(),
		flagGophKeeperServerBinaryPath,
		fork.WithEnv(envs...),
		fork.WithArgs(args...),
	)
	err := suite.gophKeeperProcess.Start(ctx)
	suite.Require().NoErrorf(err, "Невозможно запустить процесс командой %q: %s. Переменные окружения: %+v, флаги командной строки: %+v", suite.gophKeeperProcess, err, envs, args)

	err = suite.gophKeeperProcess.WaitPort(ctx, "tcp", suite.serverPort)
	suite.Require().NoErrorf(err, "Не удалось дождаться пока порт %s станет доступен для запроса: %s", suite.serverPort, err)
}

func (suite *RegistrationTest) getGophKeeperEnv() []string {

	var envs []string

	smtpHostAddress, smtpPortNumber := "127.0.0.1", smtpServer.PortNumber()
	envs = append(envs, fmt.Sprintf("SMTP_HOST=%v", smtpHostAddress))
	envs = append(envs, fmt.Sprintf("SMTP_PORT=%v", smtpPortNumber))

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
	suite.Run("registration_success", func() {

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

		suite.Require().NoError(err, "Ошибка при регистрации email: %w", err)

		// Провека email
		time.Sleep(3 * time.Second)

		// С mockSmtp возникают проблемы. TODO - заменить на container
		// msgs := smtpServer.MessagesAndPurge()
		// suite.Require().True(len(msgs) == 1, "На smtp сервере не найдено сообщений")

	})

}
