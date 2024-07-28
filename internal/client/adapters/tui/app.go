package tui

import (
	"github.com/StasMerzlyakov/gophkeeper/internal/client/app"
	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/gdamore/tcell/v2"

	"github.com/rivo/tview"
)

const (
	InitPage = "InitPage"

	LoginEMailPage = "LoginEMailPage"
	LoginOTPPage   = "LoginOTPPage"
	LoginMKeyPage  = "LoginMKeyPage"

	RegEMailPage = "RegPage"
	RegOTPPage   = "RegOTPPage"
	RegMKeyPage  = "RegMKeyPage"

	DataPageMain = "DataPageMain"
)

func NewApp(conf *config.ClientConf) *tuiApp {
	return &tuiApp{}
}

var _ app.InfoView = (*tuiApp)(nil)
var _ app.LoginView = (*tuiApp)(nil)
var _ app.RegView = (*tuiApp)(nil)

func (tApp *tuiApp) SetController(controller Controller) *tuiApp {
	tApp.controller = controller
	return tApp
}

type tuiApp struct {
	app        *tview.Application
	pages      *tview.Pages
	controller Controller

	loginFlex     *tview.Flex
	loginOTPFlex  *tview.Flex
	loginMKeyFlex *tview.Flex

	regFlex     *tview.Flex
	regOTPFlex  *tview.Flex
	regMKeyFlex *tview.Flex

	dataMainFlex *tview.Flex
}

func (tApp *tuiApp) ShowInitView() {
	tApp.app.QueueUpdateDraw(func() {
		tApp.pages.SwitchToPage(InitPage)
	})
}

func (tApp *tuiApp) ShowError(err error) {
	tApp.app.QueueUpdateDraw(func() {
		modal := tview.NewModal().
			SetText(err.Error()).
			AddButtons([]string{"Quit", "Cancel"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				switch buttonLabel {
				case "Quit":
					tApp.app.Stop()
				case "Cancel":
					tApp.app.SetRoot(tApp.pages, true).SetFocus(tApp.pages)
				}
			})
		modal.SetTitle("Error")

		tApp.app.SetRoot(modal, true).SetFocus(modal)
	})
}

func (tApp *tuiApp) ShowMsg(msg string) {
	tApp.app.QueueUpdateDraw(func() {
		modal := tview.NewModal().
			SetText(msg).
			AddButtons([]string{"Ok"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				tApp.app.SetRoot(tApp.pages, true).SetFocus(tApp.pages)
			})
		modal.SetTitle("Info")
		tApp.app.SetRoot(modal, true).SetFocus(modal)
	})
}

func (tApp *tuiApp) Start() {
	tApp.app = tview.NewApplication()

	tApp.pages = tview.NewPages()
	tApp.pages.AddPage(InitPage, tApp.createStartForm(), true, true)

	tApp.loginFlex = tview.NewFlex()
	tApp.loginMKeyFlex = tview.NewFlex()
	tApp.loginOTPFlex = tview.NewFlex()

	tApp.regFlex = tview.NewFlex()
	tApp.regMKeyFlex = tview.NewFlex()
	tApp.regOTPFlex = tview.NewFlex()

	tApp.dataMainFlex = tview.NewFlex()

	tApp.pages.AddPage(LoginEMailPage, tApp.loginFlex, true, false)
	tApp.pages.AddPage(LoginOTPPage, tApp.loginOTPFlex, true, false)
	tApp.pages.AddPage(LoginMKeyPage, tApp.loginMKeyFlex, true, false)

	tApp.pages.AddPage(RegEMailPage, tApp.regFlex, true, false)
	tApp.pages.AddPage(RegOTPPage, tApp.regOTPFlex, true, false)
	tApp.pages.AddPage(RegMKeyPage, tApp.regMKeyFlex, true, false)

	tApp.pages.AddPage(DataPageMain, tApp.dataMainFlex, true, false)

	if err := tApp.app.SetRoot(tApp.pages, true).EnableMouse(false).Run(); err != nil {
		panic(err)
	}
}

func (tApp *tuiApp) createStartForm() *tview.Flex {
	var flex = tview.NewFlex()
	var text = tview.NewTextView().
		SetTextColor(tcell.ColorGreen).
		SetText("(l) to login \n(r) to registrate\n(q) to quit")
	flex.SetDirection(tview.FlexRow).
		AddItem(text, 0, 1, false)
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'q':
			tApp.app.Stop()
		case 'r':
			tApp.ShowRegView()
		case 'l':
			tApp.ShowLoginView()
		}

		return event
	})
	return flex
}