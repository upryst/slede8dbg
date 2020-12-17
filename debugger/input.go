package debugger

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

func NewInputView() *tview.TextView {
	iv := tview.NewTextView()
	iv.SetDynamicColors(true)
	iv.SetBorder(true).SetTitle(" Input ").SetTitleAlign(tview.AlignLeft)
	return iv
}

func (ui *UI) UpdateInput() {
	var text strings.Builder
	for i := range ui.vm.Input {
		if i == ui.vm.InputIndex {
			text.WriteString("[yellow:-:b]")
		}
		text.WriteString(fmt.Sprintf("%02x", ui.vm.Input[i]))
		if i == ui.vm.InputIndex {
			text.WriteString("[-:-:-]")
		}
	}
	ui.input.SetText(text.String())
}
