package assembler

import (
	"bytes"
	"testing"
)

func TestAssemble(t *testing.T) {
	type testcase struct {
		src      string
		expected []byte
	}

	tests := []testcase{
		{`
		SETT r11, 1
		SETT r15, 0xff
		`,
			[]byte{0xb1, 0x01, 0xf1, 0xff},
		},
		{`	; a comment
			SETT r11, 1
			XELLER r0, r0
		loop1:
			LES r5
			LIK r5, r10
			BHOPP done
			SKRIV r5
			HOPP loop1

		done:
			STOPP
		`,
			[]byte{
				0xb1, 0x01, 0x25, 0x00, 0x06, 0x05,
				0x07, 0xa5, 0xe9, 0x00, 0x16, 0x05,
				0x48, 0x00, 0x00, 0x00},
		},
	}

	for _, tc := range tests {
		if output, err := Assemble(tc.src); err != nil {
			t.Fatal(err)
		} else if bytes.Compare(output, tc.expected) != 0 {
			t.Errorf("For '%s' expected %v, got %v", tc.src, tc.expected, output)
		}
	}
}

func TestAssembleSingleLine(t *testing.T) {
	type testcase struct {
		src      string
		expected []byte
	}

	tests := []testcase{
		// SETT regN, imm8
		{"SETT r0, 00", []byte{0x01, 0x00}},
		{"SETT r1, 0x11", []byte{0x11, 0x11}},
		{"SETT r2, 22h", []byte{0x21, 0x22}},
		{"SETT r15, 255", []byte{0xf1, 0xff}},
		{"SETT r14, 'A'", []byte{0xe1, 0x41}},
		{"SETT r11, 0x01", []byte{0xb1, 0x01}},
		{"SETT r11, 01h", []byte{0xb1, 0x01}},

		// SETT regN, regN
		{"SETT r0, r1", []byte{0x02, 0x01}},
		{"SETT r5, r10", []byte{0x52, 0x0a}},
		{"SETT r15, r15", []byte{0xf2, 0x0f}},

		// FINN
		{"FINN 1337", []byte{0x93, 0x53}},
		{"FINN 0x123", []byte{0x33, 0x12}},
		{"FINN 123h", []byte{0x33, 0x12}},

		// LAST / LAGR
		{"LAST r0", []byte{0x04, 0x00}},
		{"LAGR r0", []byte{0x14, 0x00}},
		{"LAST r15", []byte{0x04, 0x0f}},
		{"LAGR r15", []byte{0x14, 0x0f}},

		// ALU
		{"OG r0, r1", []byte{0x05, 0x10}},
		{"ELLER r2, r3", []byte{0x15, 0x32}},
		{"XELLER r4, r5", []byte{0x25, 0x54}},
		{"VSKIFT r6, r7", []byte{0x35, 0x76}},
		{"HSKIFT r8, r9", []byte{0x45, 0x98}},
		{"PLUSS r10, r11", []byte{0x55, 0xba}},
		{"MINUS r12, r13", []byte{0x65, 0xdc}},

		// IO
		{"LES r5", []byte{0x06, 0x05}},
		{"SKRIV r7", []byte{0x16, 0x07}},

		// CMP
		{"LIK r15, r14", []byte{0x07, 0xef}},
		{"ULIK r13, r12", []byte{0x17, 0xcd}},
		{"ME r11, r10", []byte{0x27, 0xab}},
		{"MEL r9, r8", []byte{0x37, 0x89}},
		{"SE r7, r6", []byte{0x47, 0x67}},
		{"SEL r5, r4", []byte{0x57, 0x45}},

		// HOPP / BHOPP
		{"HOPP 0xcba", []byte{0xa8, 0xcb}},
		{"BHOPP 105h", []byte{0x59, 0x10}},

		// TUR / RETUR
		{"TUR 25", []byte{0x9a, 0x01}},
		{"RETUR", []byte{0x0b, 0x00}},

		// NOPE / STOPP
		{"NOPE", []byte{0x0c, 0x00}},
		{"STOPP", []byte{0x00, 0x00}},

		{`.DATA "HELLO", 0`, []byte{'H', 'E', 'L', 'L', 'O', 0}},

		// Comments, padded / empty input
		{"ULIK r13, r12 ; a comment", []byte{0x17, 0xcd}},
		{"  RETUR", []byte{0x0b, 0x00}},
		{"BHOPP    105h  ", []byte{0x59, 0x10}},
		{"; ULIK r13, r12 ; a comment", []byte{}},
		{" ", []byte{}},
		{"", []byte{}},

		// Uppercase / lowercase
		{"uLiK R13, r12", []byte{0x17, 0xcd}},
	}

	for _, tc := range tests {
		if output, err := AssembleLine(tc.src); err != nil {
			t.Fatal(err)
		} else if bytes.Compare(output, tc.expected) != 0 {
			t.Errorf("For '%s' expected %v, got %v", tc.src, tc.expected, output)
		}
	}
}
