package vm

type OpClass int

const (
	OpClassHalt      OpClass = 0x00
	OpClassMovImm            = 0x01
	OpClassMovReg            = 0x02
	OpClassFinn              = 0x03
	OpClassLoadStore         = 0x04
	OpClassALU               = 0x05
	OpClassIO                = 0x06
	OpClassCmp               = 0x07
	OpClassJmp               = 0x08
	OpClassCondJmp           = 0x09
	OpClassCall              = 0x0a
	OpClassRet               = 0x0b
	OpClassNop               = 0x0c
)

type Instruction struct {
	Class OpClass
	Op    int
	Addr  uint16
	Val   byte
	Arg1  int
	Arg2  int
	Raw   uint16
}

func ParseInstruction(w uint16) *Instruction {
	return &Instruction{
		Raw: w,

		Class: OpClass(w & 0xf),
		Op:    int((w >> 4) & 0xf),
		Addr:  uint16(w >> 4),
		Val:   byte(w >> 8),
		Arg1:  int((w >> 8) & 0xf),
		Arg2:  int((w >> 12) & 0xf),
	}
}
