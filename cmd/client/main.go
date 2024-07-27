package main

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	ShowError(err error)
	ShowMsg(msg string)
}

type testController struct {
	exitChan <-chan struct{}
	view     testView
	conf     *config.ClientConf
}

func (tC *testController) Login(dt *domain.EMailData) {
	go func() {
		ctx, cancelFn := context.WithTimeout(context.Background(), tC.conf.InterationTimeout/2)
		defer cancelFn()

		chanDone := make(chan struct{})

		go func() {
			defer func() {
				chanDone <- struct{}{}
			}()

			// function call with context
			select {
			case <-ctx.Done():
				// cancel
			case <-time.After(tC.conf.InterationTimeout):
				tC.view.ShowMsg("Done")
			}

		}()

		select {
		case <-tC.exitChan:
			// application stopped
			return
		case <-chanDone:
			return
		case <-ctx.Done():
			tC.view.ShowError(errors.New("operation timeout"))
		}
	}()
}

func main() {

	printVersion()

	conf, err := config.LoadClientConf()
	if err != nil {
		panic(err)
	}

	tView := tui.NewApp(conf)

	exitChan := make(chan struct{})

	tCtrl := &testController{
		view:     tView,
		exitChan: exitChan,
		conf:     conf,
	}

	tView.SetController(tCtrl)

	tView.Start()

	close(exitChan)
}
