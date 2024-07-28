package tui

import (
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (tApp *tuiApp) ShowLoginPassOPTView() {
	tApp.app.QueueUpdateDraw(func() {
		tApp.loginOTPFlex.Clear()
		otpPass := &domain.OTPPass{}

		tApp.loginOTPFlex.
			SetDirection(tview.FlexRow).
			AddItem(
				tview.NewForm().
					AddPasswordField("OTPPass", "", 40, '#', func(password string) {
						otpPass.Pass = password
					}).
					AddButton("Enter", func() {
						tApp.controller.LoginPassOTP(otpPass)
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

		tApp.pages.SetTitle("Login")

		tApp.pages.SwitchToPage(LoginOTPPage)
	})
}

func (tApp *tuiApp) ShowLoginEmailView() {
	tApp.app.QueueUpdateDraw(func() {
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
						tApp.controller.LoginEMail(emailData)
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

		tApp.pages.SetTitle("Login")

		tApp.pages.SwitchToPage(LoginEMailPage)
	})
}

func (tApp *tuiApp) ShowLoginMasterKeyView(hint string) {
	tApp.app.QueueUpdateDraw(func() {
		tApp.loginMKeyFlex.Clear()
		var masterKey string

		form := tview.NewForm()
		if hint != "" {
			form.AddTextView("MasterKeyHint", hint, 0, 1, false, false)
		}

		form.AddPasswordField("MasterKey", "", 40, '#', func(mKey string) {
			masterKey = mKey
		}).
			AddButton("Enter", func() {
				tApp.controller.LoginCheckMasterKey(masterKey)
			})

		tApp.loginMKeyFlex.
			SetDirection(tview.FlexRow).
			AddItem(
				form, 0, 1, true,
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

		tApp.pages.SetTitle("Login")

		tApp.pages.SwitchToPage(LoginMKeyPage)
	})
}
