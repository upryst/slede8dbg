package vga

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/julebokk/slede8dbg/vm"
)

type UI struct {
	vm *vm.VM
	fb *Framebuffer
}

func Main(program []byte) (err error) {
	ui := &UI{
		fb: NewFramebuffer(),
	}

	if ui.vm, err = vm.NewVM(program, nil, 0); err != nil {
		return err
	}

	ui.vm.FramebufferHandler = ui.fb

	go func() {
		w := app.NewWindow(app.Title("SLEDE8++"),
			app.MinSize(unit.Px(256), unit.Px(256)),
			app.Size(unit.Px(256), unit.Px(256)))
		if err := ui.loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()

	return nil
}

func (ui *UI) loop(w *app.Window) error {
	var ops op.Ops
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			for i := 0; i < 1000000; i++ {
				if err := ui.vm.Step(); err != nil {
					return err
				}
				if ui.fb.pendingVSync {
					ui.fb.pendingVSync = false
					break
				}
			}

			gtx := layout.NewContext(&ops, e)

			ui.fb.imageOp.Add(&ops)
			paint.PaintOp{}.Add(&ops)

			e.Frame(gtx.Ops)
			w.Invalidate()
		}
	}
}
