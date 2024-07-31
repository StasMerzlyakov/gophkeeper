package tui

import (
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/client/app"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (tApp *tuiApp) ShowUserPasswordDataListView(hintList []string) {
	go func() {
		tApp.app.QueueUpdateDraw(func() {
			log := app.GetMainLogger()
			log.Debug("ShowUserPasswordDataListView start")
			tApp.userPasswordDataListFlex.Clear()

			box := tview.NewBox().SetBorder(true).SetTitle("PasswordDataList")
			tApp.userPasswordDataListFlex.Box = box

			cardNumberList := tview.NewList().ShowSecondaryText(false)
			for index, hint := range hintList {
				cardNumberList.AddItem(hint, "", rune(49+index), nil)
			}

			cardNumberList.SetSelectedFunc(func(index int, hint string, second_name string, shortcut rune) {
				tApp.controller.GetUserPasswordData(hint)
			})

			tApp.userPasswordDataListFlex.
				SetDirection(tview.FlexRow).
				AddItem(cardNumberList, 0, 1, true).
				AddItem(
					tview.NewTextView().
						SetTextColor(tcell.ColorGreen).
						SetText("(Ctrl-n) new\n(Ctrl-b) to back\n(Ctrl-q) to quit"), 0, 1, false).
				SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
					switch event.Key() {
					case tcell.KeyCtrlN:
						tApp.controller.GetUserPasswordData("")
					case tcell.KeyCtrlQ:
						tApp.app.Stop()
					case tcell.KeyCtrlB:
						tApp.pages.SwitchToPage(DataPageMain)
					}
					return event
				})
			tApp.pages.SwitchToPage(UserPasswordDataListPage)
		})
	}()
}

func (tApp *tuiApp) ShowEditUserPasswordDataView(data *domain.UserPasswordData) {
	go func() {
		tApp.app.QueueUpdateDraw(func() {
			log := app.GetMainLogger()
			log.Debugf("UserPasswordData start %v", data.Hint)

			tApp.editUserPasswordDataFlex.Clear()
			box := tview.NewBox().SetBorder(true).SetTitle(fmt.Sprintf("EditPasswordData %v", data.Hint))
			tApp.editUserPasswordDataFlex.Box = box

			tApp.editUserPasswordDataFlex.
				SetDirection(tview.FlexRow).
				AddItem(
					tview.NewForm().
						AddInputField("Login", data.Login, 40, nil, func(login string) {
							data.Login = login
						}).
						AddPasswordField("Password", data.Passwrod, 40, '#', func(pass string) {
							data.Passwrod = pass
						}).
						AddButton("Save", func() {
							tApp.controller.UpdatePasswordData(data)
						}).
						AddButton("Delete", func() {
							tApp.controller.DeleteUpdatePasswordData(data.Hint)
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
						tApp.pages.SwitchToPage(UserPasswordDataListPage)
					}
					return event
				})
			tApp.pages.SwitchToPage(EditUserPasswordDataPage)
		})
	}()
}

func (tApp *tuiApp) ShowNewUserPasswordDataView() {
	go func() {
		tApp.app.QueueUpdateDraw(func() {
			log := app.GetMainLogger()
			log.Debug("NewUserPasswordData")
			tApp.newUserPasswordDataFlex.Clear()
			box := tview.NewBox().SetBorder(true).SetTitle("NewPasswordData")
			tApp.newUserPasswordDataFlex.Box = box

			data := &domain.UserPasswordData{}

			tApp.newUserPasswordDataFlex.
				SetDirection(tview.FlexRow).
				AddItem(
					tview.NewForm().
						AddInputField("Hint", "", 40, nil, func(hint string) {
							data.Hint = hint
						}).
						AddInputField("Login", "", 40, nil, func(login string) {
							data.Login = login
						}).
						AddPasswordField("Password", "", 40, '#', func(pass string) {
							data.Passwrod = pass
						}).
						AddButton("Save", func() {
							tApp.controller.AddUserPasswordData(data)
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
						tApp.pages.SwitchToPage(UserPasswordDataListPage)
					}
					return event
				})
			tApp.pages.SwitchToPage(NewUserPasswordDataPage)
		})
	}()
}
