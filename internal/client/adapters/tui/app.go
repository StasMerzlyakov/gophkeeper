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

	NewBankCardPage  = "NewBankCardPage"
	EditBankCardPage = "NewBankCardPage"

	BankCardListPage = "BankCardListPage"

	NewUserPasswordDataPage  = "NewUserPasswordDataPage"
	EditUserPasswordDataPage = "EditUserPasswordDataPage"

	UserPasswordDataListPage = "UserPasswordDataListPage"
)

func NewApp(conf *config.ClientConf) *tuiApp {
	return &tuiApp{}
}

var _ app.AppView = (*tuiApp)(nil)

func (tApp *tuiApp) SetController(controller ViewController) *tuiApp {
	tApp.controller = controller
	return tApp
}

type tuiApp struct {
	app        *tview.Application
	pages      *tview.Pages
	controller ViewController

	loginFlex     *tview.Flex
	loginOTPFlex  *tview.Flex
	loginMKeyFlex *tview.Flex

	regFlex     *tview.Flex
	regOTPFlex  *tview.Flex
	regMKeyFlex *tview.Flex

	dataMainFlex *tview.Flex

	bankCardListFlex *tview.Flex
	newBankCardFlex  *tview.Flex
	editBankCardFlex *tview.Flex

	userPasswordDataListFlex *tview.Flex
	newUserPasswordDataFlex  *tview.Flex
	editUserPasswordDataFlex *tview.Flex
}

func (tApp *tuiApp) ShowError(err error) {
	go func() {
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
	}()
}

func (tApp *tuiApp) ShowMsg(msg string) {
	go func() {
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
	}()
}

func (tApp *tuiApp) Start() error {
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
	tApp.bankCardListFlex = tview.NewFlex()
	tApp.newBankCardFlex = tview.NewFlex()
	tApp.editBankCardFlex = tview.NewFlex()

	tApp.userPasswordDataListFlex = tview.NewFlex()
	tApp.newUserPasswordDataFlex = tview.NewFlex()
	tApp.editUserPasswordDataFlex = tview.NewFlex()

	tApp.pages.AddPage(LoginEMailPage, tApp.loginFlex, true, false)
	tApp.pages.AddPage(LoginOTPPage, tApp.loginOTPFlex, true, false)
	tApp.pages.AddPage(LoginMKeyPage, tApp.loginMKeyFlex, true, false)

	tApp.pages.AddPage(RegEMailPage, tApp.regFlex, true, false)
	tApp.pages.AddPage(RegOTPPage, tApp.regOTPFlex, true, false)
	tApp.pages.AddPage(RegMKeyPage, tApp.regMKeyFlex, true, false)

	tApp.pages.AddPage(DataPageMain, tApp.dataMainFlex, true, false)

	tApp.pages.AddPage(BankCardListPage, tApp.bankCardListFlex, true, false)
	tApp.pages.AddPage(NewBankCardPage, tApp.newBankCardFlex, true, false)
	tApp.pages.AddPage(EditBankCardPage, tApp.editBankCardFlex, true, false)

	tApp.pages.AddPage(UserPasswordDataListPage, tApp.userPasswordDataListFlex, true, false)
	tApp.pages.AddPage(NewUserPasswordDataPage, tApp.newUserPasswordDataFlex, true, false)
	tApp.pages.AddPage(EditUserPasswordDataPage, tApp.editUserPasswordDataFlex, true, false)

	if err := tApp.app.SetRoot(tApp.pages, true).EnableMouse(false).Run(); err != nil {
		log := app.GetMainLogger()
		log.Error(err)
		return err
	}

	return nil
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
