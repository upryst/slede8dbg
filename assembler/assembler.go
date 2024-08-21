package assembler

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/upryst/slede8dbg/vm"
)

type assemblerMode uint64

const (
	singleLineMode assemblerMode = iota
	multilineFirstPass
	multilineFinalPass
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

func assemble(mode assemblerMode, labels map[string]uint16,
	mnemonic, args string) ([]byte, error) {

	mnemonic = strings.ToUpper(mnemonic)

	resolveAddress := func(args string) (uint16, error) {
		addr, label, err := parseImm12OrLabel(args)
		if err != nil {
			return 0, err
		}

		if label == "" {
			return addr, nil
		}

		switch mode {
		case singleLineMode:
			return 0, errors.Errorf("Labels are not allowed in single line mode")
		case multilineFirstPass:
			// We are not interested in actual addresses at this stage
			return 0, nil
		case multilineFinalPass:
			if addr, found := labels[label]; !found {
				return 0, errors.Errorf("Label not found: %s", label)
			} else {
				return addr, nil
			}
		default:
			panic(errors.Errorf("Unhandled assemblerMode: %v", mode))
		}
	}

	switch mnemonic {
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
		if addr, err := resolveAddress(args); err != nil {
			return nil, err
		} else {
			return Bytecode(vm.OpClassFinn, Addr(addr)), nil
		}

	case "LAST", "LAGR":
		if opValue, found := loadStoreOps[mnemonic]; !found {
			panic(errors.Errorf("Unsupported Load/store op value for %s", mnemonic))
		} else if reg, err := parseReg(args); err != nil {
			return nil, err
		} else {
			return Bytecode(vm.OpClassLoadStore, Op(opValue), Arg1(reg)), nil
		}

	case "OG", "ELLER", "XELLER", "VSKIFT", "HSKIFT", "PLUSS", "MINUS":
		if opValue, found := aluOps[mnemonic]; !found {
			panic(errors.Errorf("Unsupported ALU op value for %s", mnemonic))
		} else if reg1, reg2, err := parseRegReg(args); err != nil {
			return nil, err
		} else {
			return Bytecode(vm.OpClassALU, Op(opValue), Arg1(reg1), Arg2(reg2)), nil
		}

	case "LES", "SKRIV":
		if opValue, found := ioOps[mnemonic]; !found {
			panic(errors.Errorf("Unsupported IO op value for %s", mnemonic))
		} else if reg, err := parseReg(args); err != nil {
			return nil, err
		} else {
			return Bytecode(vm.OpClassIO, Op(opValue), Arg1(reg)), nil
		}

	case "LIK", "ULIK", "ME", "MEL", "SE", "SEL":
		if opValue, found := cmpOps[mnemonic]; !found {
			panic(errors.Errorf("Unsupported cmp op value for %s", mnemonic))
		} else if reg1, reg2, err := parseRegReg(args); err != nil {
			return nil, err
		} else {
			return Bytecode(vm.OpClassCmp, Op(opValue), Arg1(reg1), Arg2(reg2)), nil
		}

	case "HOPP":
		if addr, err := resolveAddress(args); err != nil {
			return nil, err
		} else {
			return Bytecode(vm.OpClassJmp, Addr(addr)), nil
		}

	case "BHOPP":
		if addr, err := resolveAddress(args); err != nil {
			return nil, err
		} else {
			return Bytecode(vm.OpClassCondJmp, Addr(addr)), nil
		}

	case "TUR":
		if addr, err := resolveAddress(args); err != nil {
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

	case ".DATA":
		if data, err := parseData(args); err != nil {
			return nil, err
		} else {
			return data, nil
		}

	}

	if mnemonic != "" {
		return nil, errors.Errorf("Unrecognized mnemonic '%s'", mnemonic)
	}

	return nil, nil
}
