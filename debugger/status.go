package debugger

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/julebokk/slede8dbg/vm"
	"github.com/rivo/tview"
)

type StatusBar struct {
	*tview.Box

	ui *UI
}

func NewStatusBar(ui *UI) *StatusBar {
	return &StatusBar{
		Box: tview.NewBox(),
		ui:  ui,
	}
}

func (sb *StatusBar) Draw(screen tcell.Screen) {
	sb.Box.DrawForSubclass(screen, sb)
	x, y, width, _ := sb.GetInnerRect()

	var stateStr string
	switch sb.ui.vm.State {
	case vm.Running:
		stateStr = "running"
	case vm.Stopped:
		stateStr = "[yellow]stopped[-:-:-]"
	case vm.Error:
		stateStr = "[red]error[-:-:-]"
	default:
		stateStr = fmt.Sprintf("Uknown (%d)", sb.ui.vm.State)
	}

	line := fmt.Sprintf("[ [green:-:b]State[-:-:-] %s ] [ [green:-:b]Cycles:[-:-:-] %d / %d ]",
		stateStr, sb.ui.vm.CycleCount, sb.ui.vm.CycleLimit)

	tview.Print(screen, line, x, y, width, tview.AlignRight, 0)
	tview.Print(screen, "Press [green:-:b]F1[-:-:-] for help", x, y, width, tview.AlignLeft, 0)
}
