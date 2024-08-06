package tui

import (
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (tApp *tuiApp) ShowRegOTPView() {
	go func() {
		tApp.app.QueueUpdateDraw(func() {
			tApp.regOTPFlex.Clear()
			otpPass := &domain.OTPPass{}

			box := tview.NewBox().SetBorder(true).SetTitle("Registration")
			tApp.regOTPFlex.Box = box
			tApp.regOTPFlex.
				SetDirection(tview.FlexRow).
				AddItem(
					tview.NewForm().
						AddPasswordField("OTPPass", "", 40, '#', func(password string) {
							otpPass.Pass = password
						}).
						AddButton("Enter", func() {
							tApp.controller.RegPassOTP(otpPass)
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
			tApp.app.SetRoot(tApp.pages, true).SetFocus(tApp.pages)
			tApp.pages.SwitchToPage(RegOTPPage)
		})
	}()
}

func (tApp *tuiApp) ShowRegView() {
	go func() {
		tApp.app.QueueUpdateDraw(func() {
			tApp.regFlex.Clear()
			box := tview.NewBox().SetBorder(true).SetTitle("Registration")

			emailData := &domain.EMailData{}

			tApp.regFlex.Box = box

			var checkPass string

			tApp.regFlex.SetDirection(tview.FlexRow).
				AddItem(
					tview.NewForm().
						AddInputField("EMail", "", 40, nil, func(email string) {
							emailData.EMail = email
						}).
						AddPasswordField("Password", "", 40, '#', func(password string) {
							emailData.Password = password
						}).
						AddPasswordField("RetryPassword", "", 40, '#', func(password string) {
							checkPass = password
						}).
						AddButton("Enter", func() {
							if emailData.Password != checkPass {
								tApp.ShowMsg("passwords don't match")
							} else {
								tApp.controller.RegEMail(emailData)
							}
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
			tApp.app.SetRoot(tApp.pages, true).SetFocus(tApp.pages)
			tApp.pages.SwitchToPage(RegEMailPage)
		})
	}()
}

func (tApp *tuiApp) ShowRegMasterKeyView() {
	go func() {
		tApp.app.QueueUpdateDraw(func() {
			tApp.regMKeyFlex.Clear()
			masterKeyData := &domain.UnencryptedMasterKeyData{}

			box := tview.NewBox().SetBorder(true).SetTitle("Registration")
			tApp.regMKeyFlex.Box = box
			var checkKey string

			tApp.regMKeyFlex.
				SetDirection(tview.FlexRow).
				AddItem(
					tview.NewForm().
						AddInputField("MasterKeyHint", "", 40, nil, func(hint string) {
							masterKeyData.MasterPasswordHint = hint
						}).
						AddPasswordField("MasterKey", "", 40, '#', func(password string) {
							masterKeyData.MasterPassword = password
						}).
						AddPasswordField("RetryMasterKey", "", 40, '#', func(chckKey string) {
							checkKey = chckKey
						}).
						AddButton("Enter", func() {
							if masterKeyData.MasterPassword != checkKey {
								tApp.ShowMsg("keys don't match")
							} else {
								tApp.controller.RegInitMasterKey(masterKeyData)
							}

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
			tApp.app.SetRoot(tApp.pages, true).SetFocus(tApp.pages)
			tApp.pages.SwitchToPage(RegMKeyPage)
		})
	}()
}
