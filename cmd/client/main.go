package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/StasMerzlyakov/gophkeeper/internal/client/adapters/grpc/handler"
	"github.com/StasMerzlyakov/gophkeeper/internal/client/adapters/storage"
	"github.com/StasMerzlyakov/gophkeeper/internal/client/adapters/tui"
	"github.com/StasMerzlyakov/gophkeeper/internal/client/app"
	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/sirupsen/logrus"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

var log = logrus.New()

func printVersion() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}

func main() {

	printVersion() // pring to stdout

	flagSet := flag.NewFlagSet("main", flag.ContinueOnError)
	conf, err := config.LoadClientConf(flagSet)
	if err != nil {
		panic(err)
	}

	// logger
	file, err := os.OpenFile(conf.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
	log.SetLevel(logrus.DebugLevel)
	log.SetFormatter(&logrus.JSONFormatter{})

	app.SetMainLogger(log)

	// pring version to log
	log.Infof("Build version: %s", buildVersion)
	log.Infof("Build date: %s", buildDate)
	log.Infof("Build commit: %s", buildCommit)

	// grpc
	helper, err := handler.NewHandler(conf)
	if err != nil {
		panic(err)
	}

	// controller
	appCtrl := app.NewViewController(conf)

	defer func() {
		stopCtx, cnclFn := context.WithTimeout(context.Background(), 5*time.Second)
		defer cnclFn()
		appCtrl.Stop(stopCtx)
	}()

	statusWrapper := app.NewStatusWrapper(conf, helper)

	appCtrl.SetServer(statusWrapper).SetAppStorage(storage.NewStorage())

	// view
	tView := tui.NewApplicationView(conf)
	tView.SetController(appCtrl)

	appCtrl.SetInfoView(tView)

	// start application
	appCtrl.Start()
	if err := tView.Start(); err != nil {
		panic(err)
	}
}
