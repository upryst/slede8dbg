package vm

import "fmt"

func (i *Instruction) String() string {
	var comment string

	switch i.Class {
	case OpClassHalt:
		return "STOPP"

	case OpClassMovImm:
		return fmt.Sprintf("SETT r%d, 0x%02x", i.Op, i.Val)

	case OpClassMovReg:
		return fmt.Sprintf("SETT r%d, r%d", i.Op, i.Arg1)

	case OpClassFinn:
		return fmt.Sprintf("FINN 0x%03x", i.Addr)

	case OpClassLoadStore:
		if i.Op == 0 {
			return fmt.Sprintf("LAST r%d", i.Arg1)
		} else if i.Op == 1 {
			return fmt.Sprintf("LAGR r%d", i.Arg1)
		} else if i.Op == 2 {
			return fmt.Sprintf("VLAST r%d", i.Arg1)
		} else if i.Op == 3 {
			return fmt.Sprintf("VLAGR r%d", i.Arg1)
		} else {
			comment = fmt.Sprintf("unsupported load/store op %d", i.Op)
		}

	case OpClassALU:
		switch i.Op {
		case 0:
			return fmt.Sprintf("OG r%d, r%d", i.Arg1, i.Arg2)
		case 1:
			return fmt.Sprintf("ELLER r%d, r%d", i.Arg1, i.Arg2)
		case 2:
			return fmt.Sprintf("XELLER r%d, r%d", i.Arg1, i.Arg2)
		case 3:
			return fmt.Sprintf("VSKIFT r%d, r%d", i.Arg1, i.Arg2)
		case 4:
			return fmt.Sprintf("HSKIFT r%d, r%d", i.Arg1, i.Arg2)
		case 5:
			return fmt.Sprintf("PLUSS r%d, r%d", i.Arg1, i.Arg2)
		case 6:
			return fmt.Sprintf("MINUS r%d, r%d", i.Arg1, i.Arg2)
		default:
			comment = fmt.Sprintf("unsupported ALU op %d", i.Op)
		}

	case OpClassIO:
		if i.Op == 0 {
			return fmt.Sprintf("LES r%d", i.Arg1)
		} else if i.Op == 1 {
			return fmt.Sprintf("SKRIV r%d", i.Arg1)
		} else if i.Op == 2 {
			return fmt.Sprintf("INN r%d", i.Arg1)
		} else if i.Op == 3 {
			return fmt.Sprintf("UT r%d", i.Arg1)
		} else {
			comment = fmt.Sprintf("unsupported IO op %d", i.Op)
		}

	case OpClassCmp:
		switch i.Op {
		case 0:
			return fmt.Sprintf("LIK r%d, r%d", i.Arg1, i.Arg2)
		case 1:
			return fmt.Sprintf("ULIK r%d, r%d", i.Arg1, i.Arg2)
		case 2:
			return fmt.Sprintf("ME r%d, r%d", i.Arg1, i.Arg2)
		case 3:
			return fmt.Sprintf("MEL r%d, r%d", i.Arg1, i.Arg2)
		case 4:
			return fmt.Sprintf("SE r%d, r%d", i.Arg1, i.Arg2)
		case 5:
			return fmt.Sprintf("SEL r%d, r%d", i.Arg1, i.Arg2)
		default:
			comment = fmt.Sprintf("unsupported Cmp op %d", i.Op)
		}

	case OpClassJmp:
		return fmt.Sprintf("HOPP 0x%03x", i.Addr)

	case OpClassCondJmp:
		return fmt.Sprintf("BHOPP 0x%03x", i.Addr)

	case OpClassCall:
		return fmt.Sprintf("TUR 0x%03x", i.Addr)

	case OpClassRet:
		return "RETUR"

	case OpClassNop:
		return "NOPE"

	default:
		comment = fmt.Sprintf("unsupported instruction class %d", i.Class)
	}

	return fmt.Sprintf(".DATA 0x%02x, 0x%02x ; %s",
		i.Raw&0xff, i.Raw>>8, comment)
}
