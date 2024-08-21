package debugger

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/upryst/slede8dbg/vm"
)

const (
	MemSize = vm.MemSize
)

type MemoryView struct {
	*tview.TextView
	ui *UI

	offset       uint16
	bytesPerLine int
}

func NewMemoryView(ui *UI) *MemoryView {
	mv := &MemoryView{
		TextView: tview.NewTextView(),
		ui:       ui,
	}
	mv.SetWrap(false)
	mv.SetDynamicColors(true)
	mv.SetBorder(true).SetTitle(" Memory ").SetTitleAlign(tview.AlignLeft)
	return mv
}

func (mv *MemoryView) Scroll(lines int) {
	mv.offset = (mv.offset + uint16(mv.bytesPerLine*lines)) % MemSize
}

func (mv *MemoryView) Draw(screen tcell.Screen) {
	_, _, width, height := mv.TextView.GetInnerRect()

	if width > 16*4+20 {
		mv.bytesPerLine = 16
	} else {
		mv.bytesPerLine = 8
	}

	var text strings.Builder
	for i := 0; i < height; i++ {
		text.WriteString(fmt.Sprintf("  %03x: ", (mv.offset+uint16(i*mv.bytesPerLine))%MemSize))
		for j := 0; j < mv.bytesPerLine; j++ {
			text.WriteString(fmt.Sprintf("%02x ", mv.ui.vm.GetByte((mv.offset+uint16(i*mv.bytesPerLine+j))%MemSize)))
		}
		text.WriteString("  ")
		for j := 0; j < mv.bytesPerLine; j++ {
			b := mv.ui.vm.GetByte((mv.offset + uint16(i*mv.bytesPerLine+j)) % MemSize)
			if b >= ' ' && b < 0x80 {
				text.WriteByte(b)
			} else {
				text.WriteByte('.')
			}
		}
		text.WriteByte('\n')
	}
	mv.SetText(text.String())
	mv.TextView.Draw(screen)
}

func (mv *MemoryView) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return mv.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		_, _, _, height := mv.GetInnerRect()
		switch event.Key() {
		case tcell.KeyUp:
			mv.Scroll(-1)
		case tcell.KeyDown:
			mv.Scroll(1)
		case tcell.KeyPgUp:
			mv.Scroll(-height)
		case tcell.KeyPgDn:
			mv.Scroll(height)
		}
	})
}

func (mv *MemoryView) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return mv.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		x, y := event.Position()
		if !mv.InRect(x, y) {
			return false, nil
		}

		switch action {
		case tview.MouseLeftClick:
			setFocus(mv.TextView)
		case tview.MouseScrollUp:
			mv.Scroll(-1)
		case tview.MouseScrollDown:
			mv.Scroll(1)
		default:
			consumed = false
			return
		}

		consumed = true
		return
	})
}
