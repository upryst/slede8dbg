package debugger

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

func NewRegistersView() *tview.TextView {
	rv := tview.NewTextView()
	rv.SetDynamicColors(true)
	rv.SetBorderPadding(1, 1, 2, 2)
	rv.SetBorder(true).SetTitle(" Registers ").SetTitleAlign(tview.AlignLeft)
	return rv
}

func (ui *UI) UpdateRegisters() {
	var text strings.Builder

	const regColor = "[green:-:b]"
	fmtRegPair := func(i, j int) string {
		ri, rj := ui.vm.GetReg(i), ui.vm.GetReg(j)
		var riChar, rjChar string

		if ri >= ' ' && ri < 0x80 {
			riChar = fmt.Sprintf(" %c", ri)
		}

		if rj >= ' ' && rj < 0x80 {
			rjChar = fmt.Sprintf(" %c", rj)
		}

		var padding string
		if j < 10 {
			padding = " "
		}
		return fmt.Sprintf("%s%sr%d[-:-:-]: %02x%2s %s%sr%d[-:-:-]: %02x%2s",
			padding, regColor, i, ri, riChar,
			padding, regColor, j, rj, rjChar)
	}

	loadStoreStr := fmt.Sprintf(" 0x%03x/%d ", ui.vm.GetLoadStoreOffset(),
		ui.vm.GetLoadStoreOffset())

	var leftLine, rightLine string
	for i := 0; i+len(loadStoreStr) < 14; i++ {
		if i%2 == 0 {
			leftLine += string(tview.Borders.Horizontal)
		} else {
			rightLine += string(tview.Borders.Horizontal)
		}
	}

	text.WriteString(fmt.Sprintf(" [gray]%c%s[-:-:-]%s[gray:-:-]%s%c[-:-:-]\n",
		tview.Borders.TopLeft, leftLine,
		loadStoreStr,
		rightLine, tview.Borders.TopRight))
	text.WriteString(fmtRegPair(0, 1))
	text.WriteByte('\n')

	text.WriteString(fmtRegPair(2, 3))
	text.WriteByte('\n')

	text.WriteString(fmtRegPair(4, 5))
	text.WriteByte('\n')

	text.WriteString(fmtRegPair(6, 7))
	text.WriteByte('\n')

	text.WriteString(fmtRegPair(8, 9))
	text.WriteByte('\n')

	text.WriteString(fmtRegPair(10, 11))
	text.WriteByte('\n')

	text.WriteString(fmtRegPair(12, 13))
	text.WriteByte('\n')

	text.WriteString(fmtRegPair(14, 15))
	text.WriteString(fmt.Sprintf("\n\n %sPC[-:-:-]: %03x %sFlag[-:-:-]: %v\n",
		regColor, ui.vm.PC, regColor, ui.vm.Flag))

	ui.registers.SetText(text.String())
}
