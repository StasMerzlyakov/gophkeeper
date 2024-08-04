package tui

import (
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/client/app"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (tApp *tuiApp) ShowBankCardListView(cardsNumber []string) {
	go func() {
		tApp.app.QueueUpdateDraw(func() {
			log := app.GetMainLogger()
			log.Debug("ShowBankCardListView start")
			tApp.bankCardListFlex.Clear()

			box := tview.NewBox().SetBorder(true).SetTitle("BankCardList")
			tApp.bankCardListFlex.Box = box

			cardNumberList := tview.NewList().ShowSecondaryText(false)
			for index, number := range cardsNumber {
				cardNumberList.AddItem(number, "", rune(49+index), nil)
			}

			cardNumberList.SetSelectedFunc(func(index int, name string, second_name string, shortcut rune) {
				tApp.controller.GetBankCard(name)
			})

			tApp.bankCardListFlex.
				SetDirection(tview.FlexRow).
				AddItem(cardNumberList, 0, 1, true).
				AddItem(
					tview.NewTextView().
						SetTextColor(tcell.ColorGreen).
						SetText("(Ctrl-n) new\n(Ctrl-b) to back\n(Ctrl-q) to quit"), 0, 1, false).
				SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
					switch event.Key() {
					case tcell.KeyCtrlN:
						tApp.controller.NewBankCard()
					case tcell.KeyCtrlQ:
						tApp.app.Stop()
					case tcell.KeyCtrlB:
						tApp.ShowDataAccessView()
					}
					return event
				})
			tApp.app.SetRoot(tApp.pages, true).SetFocus(tApp.pages)
			tApp.pages.SwitchToPage(BankCardListPage)
			log.Debug("ShowBankCardListView shown")
		})
	}()
}

func (tApp *tuiApp) ShowNewBankCardView() {
	go func() {
		tApp.app.QueueUpdateDraw(func() {
			log := app.GetMainLogger()
			log.Debug("NewBankCard")
			box := tview.NewBox().SetBorder(true).SetTitle("NewBankCard")
			tApp.newBankCardFlex.Clear()
			tApp.newBankCardFlex.Box = box

			bankCardView := &domain.BankCardView{}

			tApp.newBankCardFlex.
				SetDirection(tview.FlexRow).
				AddItem(
					tview.NewForm().
						AddInputField("Number", "", 40, nil, func(number string) {
							bankCardView.Number = number
						}).
						AddInputField("ExpiryMonth", "", 40, nil, func(month string) {
							bankCardView.ExpiryMonth = month
						}).
						AddInputField("ExpiryYear", "", 40, nil, func(year string) {
							bankCardView.ExpiryYear = year
						}).
						AddPasswordField("CVV", "", 6, '#', func(cvv string) {
							bankCardView.CVV = cvv
						}).
						AddButton("Save", func() {
							tApp.controller.AddBankCard(bankCardView)
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
						tApp.pages.SwitchToPage(BankCardListPage)
					}
					return event
				})
			tApp.app.SetRoot(tApp.pages, true).SetFocus(tApp.pages)
			tApp.pages.SwitchToPage(NewBankCardPage)
			log.Debug("NewBankCard shown")
		})
	}()
}

func (tApp *tuiApp) ShowEditBankCardView(bankCard *domain.BankCard) {
	go func() {
		tApp.app.QueueUpdateDraw(func() {
			log := app.GetMainLogger()
			log.Debugf("ShowBankCard start %v", bankCard.Number)
			tApp.editBankCardFlex.Clear()

			box := tview.NewBox().SetBorder(true).SetTitle(fmt.Sprintf("EditBankCard %v", bankCard.Number))
			tApp.editBankCardFlex.Box = box

			bankCardView := &domain.BankCardView{
				Number:      bankCard.Number,
				ExpiryMonth: fmt.Sprintf("%02v", bankCard.ExpiryMonth),
				ExpiryYear:  fmt.Sprintf("%v", bankCard.ExpiryYear),
				CVV:         bankCard.CVV,
			}

			tApp.editBankCardFlex.
				SetDirection(tview.FlexRow).
				AddItem(
					tview.NewForm().
						AddInputField("ExpiryMonth", bankCardView.ExpiryMonth, 40, nil, func(month string) {
							bankCardView.ExpiryMonth = month
						}).
						AddInputField("ExpiryYear", bankCardView.ExpiryYear, 40, nil, func(year string) {
							bankCardView.ExpiryYear = year
						}).
						AddPasswordField("CVV", bankCardView.CVV, 6, '#', func(cvv string) {
							bankCardView.CVV = cvv
						}).
						AddButton("Save", func() {
							tApp.controller.UpdateBankCard(bankCardView)
						}).
						AddButton("Delete", func() {
							tApp.controller.DeleteBankCard(bankCardView.Number)
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
						tApp.pages.SwitchToPage(BankCardListPage)
					}
					return event
				})
			tApp.app.SetRoot(tApp.pages, true).SetFocus(tApp.pages)
			tApp.pages.SwitchToPage(EditBankCardPage)
			log.Debug("EditBankCard shown")
		})
	}()
}
