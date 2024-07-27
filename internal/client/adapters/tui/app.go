package tui

import (
	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/gdamore/tcell/v2"

	"github.com/rivo/tview"
)

const (
	InitPage  = "InitPage"
	LoginPage = "LoginPage"
)

func NewApp(conf *config.ClientConf) *tuiApp {
	return &tuiApp{}
}

func (tApp *tuiApp) SetController(controller Controller) *tuiApp {
	tApp.controller = controller
	return tApp
}

type tuiApp struct {
	app        *tview.Application
	pages      *tview.Pages
	controller Controller
	loginFlex  *tview.Flex
}

func (tApp *tuiApp) CreateMainPage() {
	var text = tview.NewTextView().
		SetTextColor(tcell.ColorGreen).
		SetText("(l) to login \n(q) to quit")

	tApp.pages.AddPage("Main", text, true, true)
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
		case 'l':
			tApp.ShowLoginView()
		}

		return event
	})
	return flex
}

func (tApp *tuiApp) ShowInitView() {
	tApp.pages.SwitchToPage(InitPage)
}

func (tApp *tuiApp) ShowError(err string) {

	modal := tview.NewModal().
		SetText(err).
		AddButtons([]string{"Ok"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			//tApp.app.Stop()
			tApp.app.SetRoot(tApp.pages, true).SetFocus(tApp.pages)
		})

	tApp.app.SetRoot(modal, true).SetFocus(modal)
}

func (tApp *tuiApp) ShowLoginView() {
	tApp.loginFlex.Clear()
	emailData := &domain.EMailData{}

	tApp.loginFlex.
		SetDirection(tview.FlexRow).
		AddItem(
			tview.NewForm().
				AddInputField("EMail", "", 40, nil, func(email string) {
					emailData.EMail = email
				}).
				AddPasswordField("Password", "", 40, '#', func(password string) {
					emailData.Password = password
				}).
				AddButton("Login", func() {
					tApp.controller.Login(emailData)
				}), 0, 1, true,
		).
		AddItem(
			tview.NewTextView().
				SetTextColor(tcell.ColorGreen).
				SetText("(b) to back\n(q) to quit"), 0, 1, false).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Rune() {
			case 'q':
				tApp.app.Stop()
			case 'b':
				tApp.pages.SwitchToPage(InitPage)
			}
			return event
		})

	tApp.pages.SwitchToPage(LoginPage)
}

func (tApp *tuiApp) Start() {
	tApp.app = tview.NewApplication()

	tApp.pages = tview.NewPages()
	tApp.pages.AddPage(InitPage, tApp.createStartForm(), true, true)
	tApp.loginFlex = tview.NewFlex()
	tApp.pages.AddPage(LoginPage, tApp.loginFlex, true, false)
	if err := tApp.app.SetRoot(tApp.pages, true).EnableMouse(false).Run(); err != nil {
		panic(err)
	}
}
