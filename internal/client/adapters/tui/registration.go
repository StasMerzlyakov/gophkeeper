package tui

import (
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (tApp *tuiApp) ShowRegPassOPTView() {
	tApp.app.QueueUpdateDraw(func() {
		tApp.regOTPFlex.Clear()
		otpPass := &domain.OTPPass{}

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

		tApp.pages.SetTitle("Registration")
		tApp.pages.SwitchToPage(RegOTPPage)
	})
}

func (tApp *tuiApp) ShowRegEmailView() {
	tApp.app.QueueUpdateDraw(func() {
		tApp.regFlex.Clear()
		emailData := &domain.EMailData{}

		tApp.regFlex.
			SetDirection(tview.FlexRow).
			AddItem(
				tview.NewForm().
					AddInputField("EMail", "", 40, nil, func(email string) {
						emailData.EMail = email
					}).
					AddPasswordField("Password", "", 40, '#', func(password string) {
						emailData.Password = password
					}).
					AddButton("Enter", func() {
						tApp.controller.RegEMail(emailData)
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

		tApp.pages.SetTitle("Registration")
		tApp.pages.SwitchToPage(RegEMailPage)
	})
}

func (tApp *tuiApp) ShowRegMasterKeyView() {
	tApp.app.QueueUpdateDraw(func() {
		tApp.regMKeyFlex.Clear()
		masterKeyData := &domain.UnencryptedMasterKeyData{}

		tApp.regMKeyFlex.
			SetDirection(tview.FlexRow).
			AddItem(
				tview.NewForm().
					AddInputField("MasterKeyHint", "", 40, nil, func(hint string) {
						masterKeyData.MasterKeyHint = hint
					}).
					AddPasswordField("MasterKey", "", 40, '#', func(password string) {
						masterKeyData.MasterKeyPassword = password
					}).
					AddButton("Enter", func() {
						tApp.controller.RegInitMasterKey(masterKeyData)
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

		tApp.pages.SetTitle("Registration")

		tApp.pages.SwitchToPage(RegMKeyPage)
	})
}
