package debugger

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	errorDialogWidth  = 50
	errorDialogHeight = 5
)

type ErrorMessage struct {
	*tview.TextView

	ui *UI
}

func (ui *UI) ShowError() {
	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewGrid().
			SetColumns(0, width, 0).
			SetRows(0, height, 0).
			AddItem(p, 1, 1, 1, 1, 0, 0, true)
	}

	errBox := &ErrorMessage{tview.NewTextView(), ui}
	errBox.SetText(ui.vm.LastError.Error()).
		SetDynamicColors(true).
		SetTitle(" Error ").
		SetBackgroundColor(tcell.ColorRed).
		SetBorder(true).
		SetBorderPadding(1, 1, 1, 1)

	errBox.SetTextAlign(tview.AlignCenter)
	errBox.SetWrap(true)

	ui.pages.AddPage("error", modal(errBox, errorDialogWidth, errorDialogHeight), true, true)

	ui.app.SetFocus(errBox)

	errBox.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape, tcell.KeyEnter, tcell.KeyCtrlSpace:
			errBox.Close()
		}
		return nil
	})
}

func (em *ErrorMessage) Close() {
	em.ui.pages.RemovePage("error")
	em.ui.pages.SwitchToPage("main")

	// TODO: be more flexible
	em.ui.app.SetFocus(em.ui.code)
}

func (em *ErrorMessage) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return em.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		x, y := event.Position()
		if !em.InRect(x, y) && action == tview.MouseLeftClick {
			em.Close()
		}

		consumed = true
		return
	})
}
