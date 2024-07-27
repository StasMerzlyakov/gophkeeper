package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (tApp *tuiApp) ShowDataAccessView() {
	tApp.app.QueueUpdateDraw(func() {
		tApp.dataMainFlex.Clear()

		tApp.dataMainFlex.
			SetDirection(tview.FlexRow).
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

		tApp.pages.SetTitle("Data page")
		tApp.pages.SwitchToPage(DataPageMain)
	})
}
