package assembler

import (
	"strings"

	"github.com/julebokk/slede8dbg/vm"
	"github.com/pkg/errors"
)

var aluOps = map[string]byte{
	"OG":     0,
	"ELLER":  1,
	"XELLER": 2,
	"VSKIFT": 3,
	"HSKIFT": 4,
	"PLUSS":  5,
	"MINUS":  6,
}

var cmpOps = map[string]byte{
	"LIK":  0,
	"ULIK": 1,
	"ME":   2,
	"MEL":  3,
	"SE":   4,
	"SEL":  5,
}

var loadStoreOps = map[string]byte{
	"LAST": 0,
	"LAGR": 1,
}

var ioOps = map[string]byte{
	"LES":   0,
	"SKRIV": 1,
}

func Assemble(line string) ([]byte, error) {
	opStr, args, err := tokenize(line)
	if err != nil {
		return nil, err
	}

	opStr = strings.ToUpper(opStr)

	switch opStr {
	case "STOPP":
		if args != "" {
			return nil, errors.Errorf("STOPP can't take arguments")
		}
		return Bytecode(vm.OpClassHalt), nil

	case "SETT":
		if reg1, reg2, imm8, err := parseRegRegOrImm8(args); err != nil {
			return nil, err
		} else if imm8 != nil {
			return Bytecode(vm.OpClassMovImm, Op(reg1), Val(*imm8)), nil
		} else {
			return Bytecode(vm.OpClassMovReg, Op(reg1), Arg1(reg2)), nil
		}

	case "FINN":
		if addr, err := parseImm12(args); err != nil {
			return nil, err
		} else {
			return Bytecode(vm.OpClassFinn, Addr(addr)), nil
		}

	case "LAST", "LAGR":
		if opValue, found := loadStoreOps[opStr]; !found {
			panic(errors.Errorf("Unsupported Load/store op value for %s", opStr))
		} else if reg, err := parseReg(args); err != nil {
			return nil, err
		} else {
			return Bytecode(vm.OpClassLoadStore, Op(opValue), Arg1(reg)), nil
		}

	case "OG", "ELLER", "XELLER", "VSKIFT", "HSKIFT", "PLUSS", "MINUS":
		if opValue, found := aluOps[opStr]; !found {
			panic(errors.Errorf("Unsupported ALU op value for %s", opStr))
		} else if reg1, reg2, err := parseRegReg(args); err != nil {
			return nil, err
		} else {
			return Bytecode(vm.OpClassALU, Op(opValue), Arg1(reg1), Arg2(reg2)), nil
		}

	case "LES", "SKRIV":
		if opValue, found := ioOps[opStr]; !found {
			panic(errors.Errorf("Unsupported IO op value for %s", opStr))
		} else if reg, err := parseReg(args); err != nil {
			return nil, err
		} else {
			return Bytecode(vm.OpClassIO, Op(opValue), Arg1(reg)), nil
		}

	case "LIK", "ULIK", "ME", "MEL", "SE", "SEL":
		if opValue, found := cmpOps[opStr]; !found {
			panic(errors.Errorf("Unsupported cmp op value for %s", opStr))
		} else if reg1, reg2, err := parseRegReg(args); err != nil {
			return nil, err
		} else {
			return Bytecode(vm.OpClassCmp, Op(opValue), Arg1(reg1), Arg2(reg2)), nil
		}

	case "HOPP":
		if addr, err := parseImm12(args); err != nil {
			return nil, err
		} else {
			return Bytecode(vm.OpClassJmp, Addr(addr)), nil
		}

	case "BHOPP":
		if addr, err := parseImm12(args); err != nil {
			return nil, err
		} else {
			return Bytecode(vm.OpClassCondJmp, Addr(addr)), nil
		}

	case "TUR":
		if addr, err := parseImm12(args); err != nil {
			return nil, err
		} else {
			return Bytecode(vm.OpClassCall, Addr(addr)), nil
		}

	case "RETUR":
		if args != "" {
			return nil, errors.Errorf("RETUR can't take arguments")
		}
		return Bytecode(vm.OpClassRet), nil

	case "NOPE":
		if args != "" {
			return nil, errors.Errorf("NOPE can't take arguments")
		}
		return Bytecode(vm.OpClassNop), nil
	}

	return nil, nil
}
