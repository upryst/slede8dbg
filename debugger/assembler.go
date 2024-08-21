package debugger

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/upryst/slede8dbg/assembler"
	"github.com/upryst/slede8dbg/vm"
)

const (
	asmDialogWidth  = 60
	asmDialogHeight = 5
)

type AsmView struct {
	*tview.Form

	edit *tview.InputField
	ui   *UI
}

func (ui *UI) ShowAsm() {
	av := &AsmView{
		Form: tview.NewForm(),
		ui:   ui,
	}

	instr := vm.ParseInstruction(ui.vm.GetWord(ui.code.lastHighlightedPC))

	av.AddInputField("", instr.String(), asmDialogWidth-6, nil, nil)

	// This is ugly, but I can't find a better way of doing it with tview
	ugly := av.GetFormItem(0)
	if edit, ok := ugly.(*tview.InputField); !ok {
		panic("I don't know how tview works")
	} else {
		av.edit = edit
	}

	av.SetFieldBackgroundColor(tcell.ColorBlack)

	av.SetBorder(true).SetTitle(" Rudimentary SLEDE8 assembler [ Esc - exit ] ")
	av.SetTitleAlign(tview.AlignLeft)

	ui.pages.AddPage("asm", makeModal(av, asmDialogWidth, asmDialogHeight), true, true)
	ui.app.SetFocus(av)

	av.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape, tcell.KeyCtrlSpace:
			av.Close()
		case tcell.KeyEnter:
			src := av.edit.GetText()
			if bytecode, err := assembler.AssembleLine(src); err == nil {
				for i := range bytecode {
					ui.vm.SetByte(uint16(i)+ui.code.offset, bytecode[i])
				}
				ui.code.offset = (ui.code.offset + uint16(len(bytecode))) % MemSize
				ui.code.lastHighlightedPC += uint16(len(bytecode))
				instr := vm.ParseInstruction(ui.vm.GetWord(ui.code.lastHighlightedPC))
				av.edit.SetText(instr.String())
				ui.status.ClearErrorText()
			} else {
				ui.status.SetErrorText(err.Error())
			}
		default:
			return event
		}

		return nil
	})
}

func (av *AsmView) Close() {
	av.ui.status.ClearErrorText()
	av.ui.pages.RemovePage("asm")
	av.ui.pages.SwitchToPage("main")

	// TODO: be more flexible
	av.ui.app.SetFocus(av.ui.code)
}

func (av *AsmView) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return av.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		x, y := event.Position()
		if !av.InRect(x, y) && action == tview.MouseLeftClick {
			av.Close()
		}

		consumed = true
		return
	})
}
