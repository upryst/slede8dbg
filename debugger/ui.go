package debugger

import (
	"github.com/rivo/tview"

	"github.com/upryst/slede8dbg/vm"
)

type UI struct {
	app *tview.Application

	code      *CodeView
	input     *tview.TextView
	memory    *MemoryView
	modal     *tview.Modal
	output    *OutputView
	pages     *tview.Pages
	registers *tview.TextView
	status    *StatusBar

	program    []byte
	inputBytes []byte
	cycleLimit int

	vm *vm.VM
}

func (ui *UI) MainLoop() error {
	return ui.app.Run()
}

func NewUI(program, inputBytes []byte, cycleLimit int) (*UI, error) {
	ui := &UI{
		app: tview.NewApplication(),

		input:     NewInputView(),
		modal:     tview.NewModal(),
		output:    NewOutputView(),
		pages:     tview.NewPages(),
		registers: NewRegistersView(),

		program:    program,
		inputBytes: inputBytes,
		cycleLimit: cycleLimit,
	}

	vm, err := vm.NewVM(program, inputBytes, cycleLimit)
	if err != nil {
		return nil, err
	}
	ui.vm = vm

	ui.code = NewCodeView(ui)
	ui.memory = NewMemoryView(ui)
	ui.status = NewStatusBar(ui)

	mainView := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(ui.input, 3, 0, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(ui.code, 0, 3, false).
				AddItem(ui.memory, 0, 2, true), 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(ui.registers, 15, 0, false).
				AddItem(ui.output, 0, 1, false),
				27, 0, false),
			0, 1, false).
		AddItem(ui.status, 1, 0, false)

	mainView.SetInputCapture(ui.HandleKeyboard)

	ui.pages.AddPage("main", mainView, true, true)

	ui.app.SetRoot(ui.pages, true).EnableMouse(true)
	ui.app.SetFocus(ui.code)

	ui.Refresh()

	return ui, nil
}

func (ui *UI) ToggleBreakpoint() {
	ui.vm.ToggleBreakpoint(ui.code.lastHighlightedPC)
}

func (ui *UI) StepVM() {
	previousState := ui.vm.State
	if err := ui.vm.Step(); err == nil {
		ui.code.offset = 0
	} else if previousState != vm.Error {
		ui.ShowError()
	}
}

func (ui *UI) RunVM() {
	previousState := ui.vm.State
	if err := ui.vm.Run(); err == nil {
		ui.code.offset = 0
	} else if previousState != vm.Error {
		ui.ShowError()
	}
}

func (ui *UI) RestartVM() {
	if newVM, err := vm.NewVM(ui.program, ui.inputBytes, ui.cycleLimit); err != nil {
		panic(err)
	} else {
		ui.vm = newVM
		ui.Refresh()
	}
}

func (ui *UI) Refresh() {
	ui.UpdateRegisters()
	ui.UpdateInput()
	ui.UpdateOutput()
}
