package assembler

import (
	"github.com/julebokk/slede8dbg/vm"
)

type bitField func(uint16) uint16

func Bytecode(oc vm.OpClass, fields ...bitField) []byte {
	word := uint16(oc)

	for _, f := range fields {
		word = f(word)
	}

	return []byte{byte(word & 0xff), byte(word >> 8)}
}

func Addr(addr uint16) bitField {
	return func(w uint16) uint16 {
		return (w & 0xf) | ((addr & 0xfff) << 4)
	}
}

func Op(v byte) bitField {
	return func(w uint16) uint16 {
		return (w & 0xff0f) | (uint16(v&0xf) << 4)
	}
}

func Val(v byte) bitField {
	return func(w uint16) uint16 {
		return (w & 0x00ff) | (uint16(v) << 8)
	}
}

func Arg1(v byte) bitField {
	return func(w uint16) uint16 {
		return (w & 0xf0ff) | (uint16(v&0xf) << 8)
	}
}

func Arg2(v byte) bitField {
	return func(w uint16) uint16 {
		return (w & 0xfff) | (uint16(v&0xf) << 12)
	}
}
