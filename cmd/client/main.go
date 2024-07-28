package main

import (
	"flag"
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/client/adapters/grpc/handler"
	"github.com/StasMerzlyakov/gophkeeper/internal/client/adapters/tui"
	"github.com/StasMerzlyakov/gophkeeper/internal/client/app"
	"github.com/StasMerzlyakov/gophkeeper/internal/config"
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
	conf, err := config.LoadClientConf(flagSet)
	if err != nil {
		panic(err)
	}

	// grpc
	helper, err := handler.NewHandler(conf)
	if err != nil {
		panic(err)
	}

	// controller
	appCtrl := app.NewAppController(conf)
	defer appCtrl.Stop()
	appCtrl.SetServer(helper)

	// view
	tView := tui.NewApp(conf)
	tView.SetController(appCtrl)

	appCtrl.SetInfoView(tView)

	tView.Start()

}
