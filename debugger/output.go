package debugger

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

type OutputView struct {
	*tview.TextView

	ascii bool
}

func (ov *OutputView) UpdateTitle() {
	if ov.ascii {
		ov.SetTitle(" Output (ASCII) ")
	} else {
		ov.SetTitle(" Output (hex) ")
	}
}

func NewOutputView() *OutputView {
	ov := &OutputView{
		TextView: tview.NewTextView(),
	}
	ov.SetDynamicColors(true)
	ov.SetTitleAlign(tview.AlignLeft)
	ov.SetBorder(true)

	ov.UpdateTitle()

	return ov
}

func (ui *UI) ToggleASCIIMode() {
	ui.output.ascii = !ui.output.ascii

	ui.output.UpdateTitle()
	ui.UpdateOutput()
}

func (ui *UI) UpdateOutput() {
	var text strings.Builder

	for _, b := range ui.vm.Output {
		if ui.output.ascii {
			if b >= ' ' && b < 0x80 {
				text.WriteByte(b)
			} else {
				text.WriteByte('.')
			}
		} else {
			text.WriteString(fmt.Sprintf("%02x", b))
		}
	}

	ui.output.SetText(text.String())
}
