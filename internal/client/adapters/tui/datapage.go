package tui

import (
	"github.com/StasMerzlyakov/gophkeeper/internal/client/app"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (tApp *tuiApp) ShowDataAccessView() {

	go func() {
		tApp.app.QueueUpdateDraw(func() {
			log := app.GetMainLogger()
			log.Debug("DataPageMain start")

			tApp.dataMainFlex.Clear()

			box := tview.NewBox().SetBorder(true).SetTitle("UserData")
			tApp.dataMainFlex.Box = box

			dataTypesList := tview.NewList().ShowSecondaryText(false)

			dataTypes := []string{
				"Bank cards", "UserPasswordData", "Files",
			}
			for index, number := range dataTypes {
				dataTypesList.AddItem(number, "", rune(49+index), nil)
			}

			dataTypesList.SetSelectedFunc(func(index int, name string, second_name string, shortcut rune) {
				switch index {
				case 0:
					tApp.controller.GetBankCardList()
				case 1:
					tApp.controller.GetUserPasswordDataList()
				default:
					tApp.controller.GetFilesInfoList()
				}
			})

			tApp.dataMainFlex.
				SetDirection(tview.FlexRow).
				AddItem(dataTypesList, 0, 1, true).
				AddItem(
					tview.NewTextView().
						SetTextColor(tcell.ColorGreen).
						SetText("(Ctrl-b) to back\n(Ctrl-q) to quit"), 0, 1, false).
				SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
					switch event.Key() {
					case tcell.KeyCtrlQ:
						tApp.app.Stop()
					case tcell.KeyCtrlB:
						tApp.pages.SwitchToPage(DataPageMain)
					}
					return event
				})
			tApp.app.SetRoot(tApp.pages, true).SetFocus(tApp.pages)
			tApp.pages.SwitchToPage(DataPageMain)
			log.Debug("DataPageMain shown")
		})
	}()
}
