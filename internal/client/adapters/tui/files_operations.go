package tui

import (
	"github.com/StasMerzlyakov/gophkeeper/internal/client/app"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (tApp *tuiApp) ShowUploadFileView(info *domain.FileInfo) {
	go func() {
		tApp.app.QueueUpdateDraw(func() {
			log := app.GetMainLogger()
			log.Debugf("ShowUploadFileView start")
			tApp.uploadFilePageFlex.Clear()

			box := tview.NewBox().SetBorder(true).SetTitle("UploadNewFile")
			tApp.uploadFilePageFlex.Box = box

			if info == nil {
				info = &domain.FileInfo{}
			}

			tApp.uploadFilePageFlex.
				SetDirection(tview.FlexRow).
				AddItem(
					tview.NewForm().
						AddInputField("Name", info.Name, 40, nil, func(name string) {
							info.Name = name
						}).
						AddInputField("Path", info.Path, 40, nil, func(path string) {
							info.Path = path
						}).
						AddButton("Upload", func() {
							tApp.controller.UploadFile(info)
						}).
						AddButton("Select File", func() {
							tApp.SelectFileView(info, func(info *domain.FileInfo) {
								tApp.ShowUploadFileView(info)
							})
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
						tApp.pages.SwitchToPage(DataPageMain)
					}
					return event
				})
			tApp.pages.SwitchToPage(UploadFilePage)
			log.Debug("ShowUploadFileView complete")
		})
	}()

}
