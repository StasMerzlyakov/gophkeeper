package tui

import (
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (tApp *tuiApp) ShowLogOTPView() {
	go func() {
		tApp.app.QueueUpdateDraw(func() {
			tApp.loginOTPFlex.Clear()
			otpPass := &domain.OTPPass{}

			box := tview.NewBox().SetBorder(true).SetTitle("Authorization")
			tApp.loginOTPFlex.Box = box

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
						SetText("(Ctrl-b) to back\n(Ctrl-q) to quit"), 0, 1, false).
				SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
					switch event.Key() {
					case tcell.KeyCtrlQ:
						tApp.app.Stop()
					case tcell.KeyCtrlB:
						tApp.pages.SwitchToPage(InitPage)
					}
					return event
				})
			tApp.pages.SwitchToPage(LoginOTPPage)
		})
	}()
}

func (tApp *tuiApp) ShowLoginView() {
	go func() {
		tApp.app.QueueUpdateDraw(func() {
			tApp.loginFlex.Clear()
			emailData := &domain.EMailData{}

			box := tview.NewBox().SetBorder(true).SetTitle("Authorization")
			tApp.loginFlex.Box = box

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
						SetText("(Ctrl-b) to back\n(Ctrl-q) to quit"), 0, 1, false).
				SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
					switch event.Key() {
					case tcell.KeyCtrlQ:
						tApp.app.Stop()
					case tcell.KeyCtrlB:
						tApp.pages.SwitchToPage(InitPage)
					}
					return event
				})
			tApp.pages.SwitchToPage(LoginEMailPage)
		})
	}()
}

func (tApp *tuiApp) ShowMasterKeyView(hint string) {
	go func() {
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

			box := tview.NewBox().SetBorder(true).SetTitle("Authorization")
			tApp.loginMKeyFlex.Box = box

			tApp.loginMKeyFlex.
				SetDirection(tview.FlexRow).
				AddItem(
					form, 0, 1, true,
				).
				AddItem(
					tview.NewTextView().
						SetTextColor(tcell.ColorGreen).
						SetText("(Ctrl-b) to back\n(Ctrl-q) to quit"), 0, 1, false).
				SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
					switch event.Key() {
					case tcell.KeyCtrlQ:
						tApp.app.Stop()
					case tcell.KeyCtrlB:
						tApp.pages.SwitchToPage(InitPage)
					}
					return event
				})
			tApp.pages.SwitchToPage(LoginMKeyPage)
		})
	}()
}
