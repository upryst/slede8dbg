package debugger

import (
	"github.com/gdamore/tcell/v2"
)

func (ui *UI) HandleKeyboard(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyF1:
		ui.ShowHelp()
	case tcell.KeyF10:
		ui.StepVM()
	case tcell.KeyF4:
		ui.code.ui.ToggleASCIIMode()
	case tcell.KeyF5:
		if event.Modifiers()&(tcell.ModCtrl|tcell.ModShift) != 0 {
			ui.RestartVM()
		} else {
			ui.RunVM()
		}
	case tcell.KeyF9:
		ui.ToggleBreakpoint()
	case tcell.KeyEnter:
		ui.ShowAsm()
	}

	if event.Modifiers()&tcell.ModAlt != 0 {
		switch event.Rune() {
		case '0':
			ui.app.SetFocus(ui.input)
		case '1':
			ui.app.SetFocus(ui.code)
		case '2':
			ui.app.SetFocus(ui.memory)
		case '3':
			ui.app.SetFocus(ui.registers)
		case '4':
			ui.app.SetFocus(ui.output)
		}
	}

	ui.Refresh()

	return event
}
