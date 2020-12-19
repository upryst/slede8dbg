package debugger

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	helpText = `
Keyboard shortcuts:

[green:-:b]Alt-0[-:-:-]  switch to Input
[green:-:b]Alt-1[-:-:-]  switch to Code
[green:-:b]Alt-2[-:-:-]  switch to Memory
[green:-:b]Alt-3[-:-:-]  switch to Registers
[green:-:b]Alt-4[-:-:-]  switch to Output
[green:-:b]Enter[-:-:-]  Assembler mode (beta)

[green:-:b]F1[-:-:-]   Help screen
[green:-:b]F4[-:-:-]   Toggle output window mode (Hex, ASCII)
[green:-:b]F5[-:-:-]   Run
[green:-:b]F9[-:-:-]   Toggle break point
[green:-:b]F10[-:-:-]  Step


[green:-:b]Ctrl-C[-:-:-]         Quit
[green:-:b]Ctrl-Shift-F5[-:-:-]  Restart debugging from scratch

Press [green:-:b]Esc[-:-:-] or [green:-:b]Enter[-:-:-]
`
)

const (
	helpViewWidth  = 50
	helpViewHeight = 24
)

type HelpView struct {
	*tview.TextView

	ui *UI
}

func (ui *UI) ShowHelp() {
	help := &HelpView{tview.NewTextView(), ui}
	help.SetText(helpText).
		SetDynamicColors(true).
		SetTitle(" Help ").
		SetBackgroundColor(tcell.ColorBlack).
		SetBorder(true).
		SetBorderPadding(0, 0, 1, 1)

	ui.pages.AddPage("help", makeModal(help, helpViewWidth, helpViewHeight), true, true)

	ui.app.SetFocus(help)

	help.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape, tcell.KeyEnter:
			help.Close()
		}
		return nil
	})
}

func (hv *HelpView) Close() {
	hv.ui.pages.RemovePage("help")
	hv.ui.pages.SwitchToPage("main")

	// TODO: be more flexible
	hv.ui.app.SetFocus(hv.ui.code)
}

func (hv *HelpView) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return hv.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		x, y := event.Position()
		if !hv.InRect(x, y) && action == tview.MouseLeftClick {
			hv.Close()
		}

		consumed = true
		return
	})
}
