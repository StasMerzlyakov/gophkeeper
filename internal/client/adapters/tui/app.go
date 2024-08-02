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
	EditBankCardPage = "EditBankCardPage"

	BankCardListPage = "BankCardListPage"

	NewUserPasswordDataPage  = "NewUserPasswordDataPage"
	EditUserPasswordDataPage = "EditUserPasswordDataPage"

	UserPasswordDataListPage = "UserPasswordDataListPage"

	UploadFilePage   = "UploadFilePage"
	FileInfoPage     = "FileInfoPage"
	FileTreePagh     = "FileTreePagh"
	FileInfoListPage = "FileInfoListPath"
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
	app         *tview.Application
	progressBar *ProgressBar
	pages       *tview.Pages
	controller  ViewController

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

	uploadFilePageFlex *tview.Flex
	fileInfoListFlex   *tview.Flex

	fileTreeView *tview.Flex
	fileInfoFlex *tview.Flex
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

func (tApp *tuiApp) ShowProgressBar(title string, progressText string, percentage float64, cancelFn func()) {
	go func() {
		tApp.app.QueueUpdateDraw(func() {
			pBar := NewProgressBar().
				AddCancelButton("Cancel").
				SetProgressText(progressText).
				SetPercentage(percentage).
				SetCancelFunc(cancelFn)
			pBar.SetTitle(title)
			tApp.app.SetRoot(pBar, true).SetFocus(pBar)
		})
	}()
}

func (tApp *tuiApp) Start() error {
	tApp.app = tview.NewApplication()

	tApp.pages = tview.NewPages()

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

	tApp.uploadFilePageFlex = tview.NewFlex()
	tApp.fileTreeView = tview.NewFlex()
	tApp.fileInfoListFlex = tview.NewFlex()
	tApp.fileInfoFlex = tview.NewFlex()
	tApp.progressBar = NewProgressBar()

	tApp.pages.AddPage(InitPage, tApp.createStartForm(), true, true)

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

	tApp.pages.AddPage(FileInfoListPage, tApp.fileInfoListFlex, true, false)
	tApp.pages.AddPage(UploadFilePage, tApp.uploadFilePageFlex, true, false)
	tApp.pages.AddPage(FileTreePagh, tApp.fileTreeView, true, false)
	tApp.pages.AddPage(FileInfoPage, tApp.fileInfoFlex, true, false)

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
