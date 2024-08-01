package tui

// Based on the https://github.com/aditya-K2/gomp/blob/master/ui/progressBar.go and tview.Modal

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ProgressBar is a centered message window used to inform
//
// See https://github.com/rivo/tview/wiki/Modal for an example.
type ProgressBar struct {
	*tview.Box

	// The frame embedded in the ProgressBar.
	frame *tview.Frame

	// The form embedded in the ProgressBar's frame.
	form *tview.Form

	// The message text (original, not word-wrapped).
	text string

	// The progress bar text
	progressText string

	// The text color.
	textColor tcell.Color

	percentage float64

	// The optional callback for when the user clicked one of the buttons. It
	// receives the index of the clicked button and the button's label.
	done func()
}

func GetProgressGlyph(width, percentage float64, btext string) string {
	q := "[black:white:b]"
	var a string
	a += strings.Repeat(" ", int(width)-len(btext))
	a = InsertAt(a, btext, int(width/2)-10)
	a = InsertAt(a, "[-:-:-]", int(width*percentage/100))
	q += a
	return q
}

func (m *ProgressBar) SetPercentage(percentage float64) *ProgressBar {
	m.percentage = percentage
	return m
}

func InsertAt(inputString, stringTobeInserted string, index int) string {
	s := inputString[:index] + stringTobeInserted + inputString[index:]
	return s
}

// NewProgressBar returns a new ProgressBar message window.
func NewProgressBar() *ProgressBar {
	m := &ProgressBar{
		Box:       tview.NewBox().SetBorder(true).SetBackgroundColor(tview.Styles.ContrastBackgroundColor),
		textColor: tview.Styles.PrimaryTextColor,
	}
	m.form = tview.NewForm().
		SetButtonsAlign(tview.AlignCenter).
		SetButtonBackgroundColor(tview.Styles.PrimitiveBackgroundColor).
		SetButtonTextColor(tview.Styles.PrimaryTextColor)
	m.form.SetBackgroundColor(tview.Styles.ContrastBackgroundColor).SetBorderPadding(0, 0, 0, 0)
	m.form.SetCancelFunc(func() {
		if m.done != nil {
			m.done()
		}
	})
	m.frame = tview.NewFrame(m.form).SetBorders(0, 0, 1, 0, 0, 0)
	m.frame.SetBackgroundColor(tview.Styles.ContrastBackgroundColor).
		SetBorderPadding(1, 1, 1, 1)
	return m
}

// SetBackgroundColor sets the color of the ProgressBar frame background.
func (m *ProgressBar) SetBackgroundColor(color tcell.Color) *ProgressBar {
	m.form.SetBackgroundColor(color)
	m.frame.SetBackgroundColor(color)
	return m
}

// SetTextColor sets the color of the message text.
func (m *ProgressBar) SetTextColor(color tcell.Color) *ProgressBar {
	m.textColor = color
	return m
}

// SetButtonBackgroundColor sets the background color of the buttons.
func (m *ProgressBar) SetButtonBackgroundColor(color tcell.Color) *ProgressBar {
	m.form.SetButtonBackgroundColor(color)
	return m
}

// SetButtonTextColor sets the color of the button texts.
func (m *ProgressBar) SetButtonTextColor(color tcell.Color) *ProgressBar {
	m.form.SetButtonTextColor(color)
	return m
}

// SetButtonStyle sets the style of the buttons when they are not focused.
func (m *ProgressBar) SetButtonStyle(style tcell.Style) *ProgressBar {
	m.form.SetButtonStyle(style)
	return m
}

// SetButtonActivatedStyle sets the style of the buttons when they are focused.
func (m *ProgressBar) SetButtonActivatedStyle(style tcell.Style) *ProgressBar {
	m.form.SetButtonActivatedStyle(style)
	return m
}

// SetDoneFunc sets a handler which is called when one of the buttons was
// pressed. It receives the index of the button as well as its label text. The
// handler is also called when the user presses the Escape key. The index will
// then be negative and the label text an empty string.
func (m *ProgressBar) SetDoneFunc(handler func()) *ProgressBar {
	m.done = handler
	return m
}

// SetText sets the message text of the window. The text may contain line
// breaks but style tag states will not transfer to following lines. Note that
// words are wrapped, too, based on the final size of the window.
func (m *ProgressBar) SetText(text string) *ProgressBar {
	m.text = text
	return m
}

// SetProgressText sets the message text of the progress bar window.
func (m *ProgressBar) SetProgressText(progressText string) *ProgressBar {
	m.progressText = progressText
	return m
}

// AddButtons adds buttons to the window. There must be at least one button and
// a "done" handler so the window can be closed again.
func (m *ProgressBar) AddCancelButton(label string) *ProgressBar {

	m.form.AddButton(label, func() {
		if m.done != nil {
			m.done()
		}
	})
	button := m.form.GetButton(m.form.GetButtonCount() - 1)
	button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDown, tcell.KeyRight:
			return tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone)
		case tcell.KeyUp, tcell.KeyLeft:
			return tcell.NewEventKey(tcell.KeyBacktab, 0, tcell.ModNone)
		}
		return event
	})

	return m
}

// ClearButtons removes all buttons from the window.
func (m *ProgressBar) ClearButtons() *ProgressBar {
	m.form.ClearButtons()
	return m
}

// SetFocus shifts the focus to the button with the given index.
func (m *ProgressBar) SetFocus(index int) *ProgressBar {
	m.form.SetFocus(index)
	return m
}

// Focus is called when this primitive receives focus.
func (m *ProgressBar) Focus(delegate func(p tview.Primitive)) {
	delegate(m.form)
}

// HasFocus returns whether or not this primitive has focus.
func (m *ProgressBar) HasFocus() bool {
	return m.form.HasFocus()
}

// Draw draws this primitive onto the screen.
func (m *ProgressBar) Draw(screen tcell.Screen) {
	// Calculate the width of this ProgressBar.
	buttonsWidth := 0
	for i := 0; i < m.form.GetButtonCount(); i++ {
		button := m.form.GetButton(i)
		buttonsWidth += tview.TaggedStringWidth(button.GetTitle()) + 4 + 2
	}

	buttonsWidth -= 2
	screenWidth, screenHeight := screen.Size()
	width := screenWidth / 3
	if width < buttonsWidth {
		width = buttonsWidth
	}
	// width is now without the box border.

	// Reset the text and find out how wide it is.
	m.frame.Clear()
	lines := tview.WordWrap(m.text, width) // added a new line for reserving space for the progress bar
	for _, line := range lines {
		m.frame.AddText(line, true, tview.AlignCenter, m.textColor)
	}

	// Set the ProgressBar's position and size.
	height := len(lines) + 6
	width += 4
	x := (screenWidth - width) / 2
	y := (screenHeight - height) / 2
	m.SetRect(x, y, width, height)

	// Draw the frame.
	m.Box.DrawForSubclass(screen, m)
	x, y, width, height = m.GetInnerRect()
	m.frame.SetRect(x, y, width, height)
	m.frame.Draw(screen)

	var (
		OFFSET int = 1
	)

	tview.Print(screen,
		GetProgressGlyph(float64(width-OFFSET-1),
			m.percentage,
			m.progressText),
		x, y+1, width-OFFSET, tview.AlignRight, tcell.ColorWhite)
}

// MouseHandler returns the mouse handler for this primitive.
func (m *ProgressBar) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return m.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		// Pass mouse events on to the form.
		consumed, capture = m.form.MouseHandler()(action, event, setFocus)
		if !consumed && action == tview.MouseLeftDown && m.InRect(event.Position()) {
			setFocus(m)
			consumed = true
		}
		return
	})
}

// InputHandler returns the handler for this primitive.
func (m *ProgressBar) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return m.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if m.frame.HasFocus() {
			if handler := m.frame.InputHandler(); handler != nil {
				handler(event, setFocus)
				return
			}
		}
	})
}
