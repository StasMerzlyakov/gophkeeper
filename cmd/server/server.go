package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/server/adapters/grpc/handler"
	"github.com/StasMerzlyakov/gophkeeper/internal/server/domain"
	"go.uber.org/zap"
)

func main() {

	conf, err := config.LoadServConf()
	if err != nil {
		panic(err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("cannot initialize zap")
	}
	defer logger.Sync()

	sugarLog := logger.Sugar()

	domain.SetApplicationLogger(sugarLog)

	srvCtx, cancelFn := context.WithCancel(context.Background())

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	handler := handler.NewGRPCHandler(conf)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-exit
		cancelFn()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		handler.Start(srvCtx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-srvCtx.Done()
		handler.Stop()
	}()

	wg.Wait()

}
