package vm

import (
	"bytes"

	"github.com/pkg/errors"
)

const (
	RegCount = 16
	MemSize  = 4096

	SledeHeader = ".SLEDE8"
)

type VMState int

const (
	Running VMState = iota
	Stopped
	Error
)

var (
	ErrNoMoreInput        = errors.New("No more input available")
	ErrEmptyStack         = errors.New("Stack is empty")
	ErrCycleLimitExceeded = errors.New("Cycle limit exceeded")
)

type VM struct {
	Flag bool
	PC   uint16
	Regs [RegCount]byte

	Mem [MemSize]byte

	Input      []byte
	InputIndex int

	Output []byte
	Stack  []uint16

	CycleCount int
	CycleLimit int

	State     VMState
	LastError error

	// Single breakpoint "covers" the whole word
	BreakPoints [MemSize / 2]bool

	IOHandler          IOHandler
	FramebufferHandler FramebufferHandler
}

func NewVM(program, input []byte, cycleLimit int) (*VM, error) {
	r := bytes.NewReader(program)

	magic := make([]byte, len(SledeHeader))
	if _, err := r.Read(magic); err != nil {
		return nil, errors.WithStack(err)
	}
	if string(magic) != SledeHeader {
		return nil, errors.Errorf("Expected %s header", SledeHeader)
	}

	codeSize := len(program) - len(SledeHeader)
	if codeSize > MemSize {
		return nil, errors.Errorf("Program size (%d) exceeds memory limit (%d)",
			codeSize, MemSize)
	}

	vm := &VM{
		Input:      input,
		CycleLimit: cycleLimit,

		IOHandler:          dummyIOHandler{},
		FramebufferHandler: dummyFrameBuffer{},
	}

	if _, err := r.Read(vm.Mem[:]); err != nil {
		return nil, errors.WithStack(err)
	}

	return vm, nil
}

func (vm *VM) Run() error {
	for vm.State == Running {
		if err := vm.Step(); err != nil {
			return err
		}
		if vm.BreakpointSet(vm.PC) {
			break
		}
	}
	return nil
}

func (vm *VM) ToggleBreakpoint(addr uint16) {
	addr = (addr % MemSize) >> 1
	vm.BreakPoints[addr] = !vm.BreakPoints[addr]
}

func (vm *VM) BreakpointSet(addr uint16) bool {
	return vm.BreakPoints[(addr%MemSize)>>1]
}

func (vm *VM) Step() error {
	if vm.CycleLimit > 0 && vm.CycleCount >= vm.CycleLimit {
		return vm.setError(ErrCycleLimitExceeded)
	}

	var nextPC *uint16

	i := ParseInstruction(vm.GetWord(vm.PC))
	switch i.Class {
	case OpClassHalt:
		vm.State = Stopped
		return nil

	case OpClassMovImm:
		vm.SetReg(i.Op, i.Val)

	case OpClassMovReg:
		vm.SetReg(i.Op, vm.GetReg(i.Arg1))

	case OpClassFinn:
		vm.SetReg(0, byte(i.Addr&0xff))
		vm.SetReg(1, byte(i.Addr>>8))

	case OpClassLoadStore:
		if i.Op == 0 {
			vm.SetReg(i.Arg1, vm.GetByte(vm.GetLoadStoreOffset()))
		} else if i.Op == 1 {
			vm.SetByte(vm.GetLoadStoreOffset(), vm.GetReg(i.Arg1))
		} else if i.Op == 2 {
			vm.SetReg(i.Arg1, vm.FramebufferHandler.Read(vm.GetFramebufferAddress()))
		} else if i.Op == 3 {
			vm.FramebufferHandler.Write(vm.GetFramebufferAddress(), vm.GetReg(i.Arg1))
		} else {
			return vm.setError(errors.Errorf("Unsupported load/store op %d (PC %04x)",
				i.Op, vm.PC))
		}

	case OpClassALU:
		switch i.Op {
		case 0:
			vm.SetReg(i.Arg1, vm.GetReg(i.Arg1)&vm.GetReg(i.Arg2))
		case 1:
			vm.SetReg(i.Arg1, vm.GetReg(i.Arg1)|vm.GetReg(i.Arg2))
		case 2:
			vm.SetReg(i.Arg1, vm.GetReg(i.Arg1)^vm.GetReg(i.Arg2))
		case 3:
			vm.SetReg(i.Arg1, vm.GetReg(i.Arg1)<<vm.GetReg(i.Arg2))
		case 4:
			vm.SetReg(i.Arg1, vm.GetReg(i.Arg1)>>vm.GetReg(i.Arg2))
		case 5:
			vm.SetReg(i.Arg1, vm.GetReg(i.Arg1)+vm.GetReg(i.Arg2))
		case 6:
			vm.SetReg(i.Arg1, vm.GetReg(i.Arg1)-vm.GetReg(i.Arg2))
		default:
			return vm.setError(errors.Errorf("Unsupported ALU op %d (PC %04x)",
				i.Op, vm.PC))
		}

	case OpClassIO:
		if i.Op == 0 {
			if value, err := vm.readInput(); err != nil {
				return err
			} else {
				vm.SetReg(i.Arg1, value)
			}
		} else if i.Op == 1 {
			vm.writeOutput(vm.GetReg(i.Arg1))
		} else if i.Op == 2 {
			vm.SetReg(i.Arg1, vm.IOHandler.ReadPort(vm.GetIOAddress()))
		} else if i.Op == 3 {
			vm.IOHandler.WritePort(vm.GetIOAddress(), vm.GetReg(i.Arg1))
		} else if i.Op == 4 {
			vm.FramebufferHandler.VSync()
		} else {
			return vm.setError(errors.Errorf("Unsupported IO op %d (PC %04x)",
				i.Op, vm.PC))
		}

	case OpClassCmp:
		switch i.Op {
		case 0:
			vm.Flag = vm.GetReg(i.Arg1) == vm.GetReg(i.Arg2)
		case 1:
			vm.Flag = vm.GetReg(i.Arg1) != vm.GetReg(i.Arg2)
		case 2:
			vm.Flag = vm.GetReg(i.Arg1) < vm.GetReg(i.Arg2)
		case 3:
			vm.Flag = vm.GetReg(i.Arg1) <= vm.GetReg(i.Arg2)
		case 4:
			vm.Flag = vm.GetReg(i.Arg1) > vm.GetReg(i.Arg2)
		case 5:
			vm.Flag = vm.GetReg(i.Arg1) >= vm.GetReg(i.Arg2)
		default:
			return vm.setError(errors.Errorf("Unsupported Cmp op %d (PC %04x)",
				i.Op, vm.PC))
		}

	case OpClassJmp:
		nextPC = &i.Addr

	case OpClassCondJmp:
		if vm.Flag {
			nextPC = &i.Addr
		}

	case OpClassCall:
		nextPC = &i.Addr
		vm.push(vm.PC + 2)

	case OpClassRet:
		if retAddr, err := vm.pop(); err != nil {
			return err
		} else {
			nextPC = &retAddr
		}

	case OpClassNop:

	default:
		return vm.setError(errors.Errorf("Unsupported instruction class"+
			" %d (PC %04x)", i.Class, vm.PC))
	}

	if nextPC != nil {
		vm.PC = *nextPC
	} else {
		vm.PC = (vm.PC + 2) % MemSize
	}

	vm.CycleCount++

	return nil
}

func regIndexSanityCheck(reg int) {
	if reg < 0 || reg >= RegCount {
		panic(errors.Errorf("Reg index (%d) is out of valid range", reg))
	}
}

func (vm *VM) GetWord(offset uint16) uint16 {
	return uint16(vm.Mem[offset%MemSize]) +
		uint16(vm.Mem[(offset+1)%MemSize])<<8
}

func (vm *VM) GetLoadStoreOffset() uint16 {
	return uint16(vm.GetReg(0)) + uint16(vm.GetReg(1))<<8
}

// Just for code clarity, equivalent to GetLoadStoreOffset()
func (vm *VM) GetIOAddress() uint16 {
	return uint16(vm.GetReg(0)) + uint16(vm.GetReg(1))<<8
}

// Just for code clarity, equivalent to GetLoadStoreOffset()
func (vm *VM) GetFramebufferAddress() uint16 {
	return uint16(vm.GetReg(0)) + uint16(vm.GetReg(1))<<8
}

func (vm *VM) GetByte(offset uint16) byte {
	return vm.Mem[offset%MemSize]
}

func (vm *VM) SetByte(offset uint16, value byte) {
	vm.Mem[offset%MemSize] = value
}

func (vm *VM) SetReg(reg int, value byte) {
	regIndexSanityCheck(reg)
	vm.Regs[reg] = value
}

func (vm *VM) GetReg(reg int) byte {
	regIndexSanityCheck(reg)
	return vm.Regs[reg]
}

func (vm *VM) readInput() (value byte, err error) {
	if vm.InputIndex >= len(vm.Input) {
		return 0, vm.setError(ErrNoMoreInput)
	}

	value = vm.Input[vm.InputIndex]
	vm.InputIndex++

	return
}

func (vm *VM) writeOutput(value byte) {
	vm.Output = append(vm.Output, value)
}

func (vm *VM) push(value uint16) {
	// TODO: consider introducing a limit
	vm.Stack = append(vm.Stack, value)
}

func (vm *VM) pop() (value uint16, err error) {
	if len(vm.Stack) == 0 {
		return 0, vm.setError(ErrEmptyStack)
	}

	value = vm.Stack[len(vm.Stack)-1]
	vm.Stack = vm.Stack[:len(vm.Stack)-1]

	return
}

func (vm *VM) setError(err error) error {
	vm.State = Error
	vm.LastError = err
	return err
}
