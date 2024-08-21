package debugger

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/upryst/slede8dbg/vm"
)

type CodeView struct {
	*tview.TextView
	ui *UI

	offset            uint16
	lastHighlightedPC uint16
}

func NewCodeView(ui *UI) *CodeView {
	cv := &CodeView{
		TextView: tview.NewTextView(),
		ui:       ui,
	}
	cv.SetWrap(false)
	cv.SetDynamicColors(true)
	cv.SetBorder(true).SetTitle(" Code ").SetTitleAlign(tview.AlignLeft)
	return cv
}

func (cv *CodeView) Draw(screen tcell.Screen) {
	x, y, width, height := cv.TextView.GetInnerRect()

	middle := height >> 1

	cv.TextView.DrawForSubclass(screen, cv)
	for i := 0; i < height; i++ {
		offset := (uint16((i-middle)*2) + cv.ui.code.offset + cv.ui.vm.PC) % MemSize
		hasBreakpoint := cv.ui.vm.BreakpointSet(offset)
		instr := vm.ParseInstruction(cv.ui.vm.GetWord(offset))

		var symbol, color string
		if hasBreakpoint {
			symbol = "🛑"
		} else if offset == cv.ui.vm.PC {
			symbol = "👉"
		} else {
			symbol = "  "
		}

		if i == middle {
			// keeping for toggling breakpoints
			cv.lastHighlightedPC = uint16(offset)
		}

		if offset == cv.ui.vm.PC {
			if i == middle {
				color = "[green:gray:b]"
			} else {
				color = "[green::b]"
			}
		} else if i == middle {
			color = "[:gray:b]"
		} else if hasBreakpoint {
			color = "[:red:b]"
		}

		_, printedWidth := tview.Print(screen, fmt.Sprintf("%s%s%03x: %02x%02x  %s",
			color, symbol, offset, instr.Raw&0xff, instr.Raw>>8, instr.String()),
			x, y+i, width, tview.AlignLeft, 0)

		if color != "" {
			// Hack! Fill the rest of the line
			_, _, style, _ := screen.GetContent(x+1, y+i)
			for printedWidth < width {
				screen.SetContent(x+printedWidth, y+i, ' ', nil, style)
				printedWidth++
			}
		}
	}
}

func (cv *CodeView) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return cv.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		_, _, _, height := cv.GetInnerRect()
		switch event.Key() {
		case tcell.KeyUp:
			cv.offset = (cv.offset - 2) % MemSize
		case tcell.KeyDown:
			cv.offset = (cv.offset + 2) % MemSize
		case tcell.KeyPgUp:
			cv.offset = (cv.offset - uint16(height)) % MemSize
		case tcell.KeyPgDn:
			cv.offset = (cv.offset + uint16(height)) % MemSize
		}
	})
}

func (cv *CodeView) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return cv.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		x, y := event.Position()
		if !cv.InRect(x, y) {
			return false, nil
		}

		switch action {
		case tview.MouseLeftClick:
			setFocus(cv)
		case tview.MouseScrollUp:
			cv.offset = (cv.offset - 2) % MemSize
		case tview.MouseScrollDown:
			cv.offset = (cv.offset + 2) % MemSize
		default:
			consumed = false
			return
		}

		consumed = true
		return
	})
}
