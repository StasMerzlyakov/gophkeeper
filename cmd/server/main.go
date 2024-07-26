package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/StasMerzlyakov/gophkeeper/internal/server/adapters/grpc/handler"
	"go.uber.org/zap"
)

func main() {

	conf, err := config.LoadServConf()
	if err != nil {
		panic(err)
	}

	if err := domain.CheckServerSecretKeyComplexityLevel(conf.ServerEncryptionKey); err != nil {
		panic(err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		panic("cannot initialize zap")
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			panic(err)
		}
	}()

	sugarLog := logger.Sugar()

	domain.SetApplicationLogger(sugarLog)

	srvCtx, cancelFn := context.WithCancel(context.Background())

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	handler := handler.NewGRPCHandler(conf)

	handler.Start(srvCtx)

	go func() {
		<-srvCtx.Done()
		handler.Stop()
	}()

	<-exit
	cancelFn()
}
