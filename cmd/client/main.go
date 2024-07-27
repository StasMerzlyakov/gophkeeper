package main

import (
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/client/adapters/tui"
	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
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

type testView interface {
	ShowError(err string)
}

type testController struct {
	view testView
}

func (tC *testController) Login(dt *domain.EMailData) {
	tC.view.ShowError("Hello")
}

func main() {

	printVersion()

	conf, err := config.LoadClientConf()
	if err != nil {
		panic(err)
	}

	tView := tui.NewApp(conf)

	tCtrl := &testController{
		view: tView,
	}

	tView.SetController(tCtrl)

	tView.Start()
}
