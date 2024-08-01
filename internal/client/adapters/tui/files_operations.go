package tui

import (
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/client/app"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (tApp *tuiApp) ShowFileInfoListView(filesInfoList []domain.FileInfo) {
	go func() {
		tApp.app.QueueUpdateDraw(func() {
			log := app.GetMainLogger()
			log.Debug("ShowFileInfListView start")
			tApp.fileInfoListFlex.Clear()

			box := tview.NewBox().SetBorder(true).SetTitle("FileInfoList")
			tApp.fileInfoListFlex.Box = box

			cardNumberList := tview.NewList().ShowSecondaryText(false)
			for index, info := range filesInfoList {
				expl := "has local copy"
				if info.Path == "" {
					expl = "only on server"
				}

				cardNumberList.AddItem(info.Name, expl, rune(49+index), nil)
			}

			cardNumberList.SetSelectedFunc(func(index int, name string, second_name string, shortcut rune) {
				tApp.controller.GetFileInfo(name)
			})

			tApp.fileInfoListFlex.
				SetDirection(tview.FlexRow).
				AddItem(cardNumberList, 0, 1, true).
				AddItem(
					tview.NewTextView().
						SetTextColor(tcell.ColorGreen).
						SetText("(Ctrl-n) new\n(Ctrl-b) to back\n(Ctrl-q) to quit"), 0, 1, false).
				SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
					switch event.Key() {
					case tcell.KeyCtrlN:
						tApp.ShowUploadFileView(nil)
					case tcell.KeyCtrlQ:
						tApp.app.Stop()
					case tcell.KeyCtrlB:
						tApp.pages.SwitchToPage(DataPageMain)
					}
					return event
				})
			tApp.pages.SwitchToPage(FileInfoListPage)
			log.Debug("ShowFileInfListView shown")
		})
	}()
}

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

func (tApp *tuiApp) ShowFileInfoView(info *domain.FileInfo) {
	go func() {
		tApp.app.QueueUpdateDraw(func() {
			log := app.GetMainLogger()
			log.Debugf("ShowFileInfoView start")
			tApp.fileInfoFlex.Clear()

			box := tview.NewBox().SetBorder(true).SetTitle(fmt.Sprintf("FileInfo %s", info.Name))
			tApp.fileInfoFlex.Box = box

			if info == nil {
				info = &domain.FileInfo{}
			}

			tApp.fileInfoFlex.
				SetDirection(tview.FlexRow).
				AddItem(
					tview.NewForm().
						AddInputField("Path", info.Path, 40, nil, func(path string) {
							info.Path = path
						}).
						AddButton("Delete", func() {
							tApp.controller.DeleteFile(info.Name)
						}).
						AddButton("Save", func() {
							tApp.controller.SaveFile(info)
						}).
						AddButton("Select directory to save", func() {
							tApp.SelectDirectoryView(info, func(info *domain.FileInfo) {
								tApp.ShowFileInfoView(info)
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
			tApp.pages.SwitchToPage(FileInfoPage)
			log.Debug("ShowFileInfoView complete")
		})
	}()

}
