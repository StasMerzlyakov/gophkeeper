package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (tApp *tuiApp) ShowDataAccessView() {

	go func() {
		tApp.app.QueueUpdateDraw(func() {
			tApp.dataMainFlex.Clear()

			box := tview.NewBox().SetBorder(true).SetTitle("UserData")
			tApp.dataMainFlex.Box = box

			tApp.dataMainFlex.
				SetDirection(tview.FlexRow).
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
			tApp.pages.SwitchToPage(DataPageMain)
		})
	}()
}
