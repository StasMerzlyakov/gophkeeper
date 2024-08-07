package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/StasMerzlyakov/gophkeeper/internal/server/adapters/bigfiles/fs"
	"github.com/StasMerzlyakov/gophkeeper/internal/server/adapters/email"
	"github.com/StasMerzlyakov/gophkeeper/internal/server/adapters/grpc/handler"
	"github.com/StasMerzlyakov/gophkeeper/internal/server/adapters/storage/postgres"
	"github.com/StasMerzlyakov/gophkeeper/internal/server/adapters/ttlstorage"
	"github.com/StasMerzlyakov/gophkeeper/internal/server/usecases"
	"go.uber.org/zap"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func printVersion() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}

func main() {

	printVersion()

	flagSet := flag.NewFlagSet("main", flag.ContinueOnError)
	conf, err := config.LoadServConf(flagSet)
	if err != nil {
		panic(err)
	}

	if err := domain.CheckServerSecretKeyComplexityLevel(conf.ServerSecret); err != nil {
		panic(err)
	}

	logger, err := zap.NewDevelopment()

	if err != nil {
		panic("cannot initialize zap")
	}

	defer func() {
		// Баг в zap  : https://github.com/uber-go/zap/issues/772
		// Sync возвращает ошибку
		// invalid argumentio/fs.PathError {Op: "sync", Path: "/dev/stderr", Err: error(syscall.Errno) EINVAL (22)}
		_ = logger.Sync()
	}()

	sugarLog := logger.Sugar()

	domain.SetApplicationLogger(sugarLog)

	srvCtx, cancelFn := context.WithCancel(context.Background())

	// stateless storage
	memStorage := ttlstorage.NewMemStorage(srvCtx, conf)

	// statefull storage
	pgStorage := postgres.NewStorage(srvCtx, conf)
	defer pgStorage.Close()

	// fileStorage
	fileStorage := fs.NewFileStorage(conf)

	// email sender
	sender := email.NewSender(conf)
	if err := sender.Connect(srvCtx); err != nil {
		panic(err)
	}

	defer func() {
		if err := sender.Close(); err != nil {
			sugarLog.Error(err) // sometimes occurs "write tcp 127.0.0.1:36690->127.0.0.1:25: write: broken pipe". TODO find out reason
		}
	}()

	// application
	helper := usecases.NewRegistrationHelper(conf, rand.Read)
	registrator := usecases.NewRegistrator(conf).
		RegistrationHelper(helper).
		StateFullStorage(pgStorage).
		TemporaryStorage(memStorage).
		EMailSender(sender)
	autHelper := usecases.NewAuth(conf).
		RegistrationHelper(helper).
		StateFullStorage(pgStorage).
		TemporaryStorage(memStorage)
	dataAccess := usecases.NewDataAccessor(conf).
		StateFullStorage(pgStorage)
	fileAcces := usecases.NewFileAccessor(conf).StateFullStorage(pgStorage).FileStorage(fileStorage)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// grpc
	grpcRegHandler := handler.NewRegHandler(registrator)
	grpcAuthService := handler.NewAuthService(autHelper)
	grpcDataAccessor := handler.NewDataAccessor(dataAccess)
	grpcFileAccessor := handler.NewFileAccessor(fileAcces)

	handler := handler.NewGRPCHandler(conf).
		AuthService(grpcAuthService).
		DataAccessor(grpcDataAccessor).
		RegHandler(grpcRegHandler).
		FileAccessor(grpcFileAccessor)

	handler.Start(srvCtx)

	go func() {
		<-srvCtx.Done()
		handler.Stop()
	}()

	<-exit
	cancelFn()
}
