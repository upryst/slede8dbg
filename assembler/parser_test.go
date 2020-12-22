package assembler

import (
	"bytes"
	"testing"
)

func TestParseData(t *testing.T) {
	type testcase struct {
		src      string
		expected []byte
	}

	tests := []testcase{
		{"1,2,3, 4,  5,   6", []byte{1, 2, 3, 4, 5, 6}},
		{"42, 223", []byte{42, 223}},
		{"' '", []byte{' '}},
		{"';'", []byte{';'}},
		{"1,2,3,'A',4,5", []byte{1, 2, 3, 'A', 4, 5}},
		{`"A string", 0`, []byte{'A', ' ', 's', 't', 'r', 'i', 'n', 'g', 0}},
		{`"", 0`, []byte{0}},
		{`"Hi\n"`, []byte{'H', 'i', '\n'}},
		{`1, "A", 2, "BCD"`, []byte{1, 'A', 2, 'B', 'C', 'D'}},
	}

	for _, tc := range tests {
		if output, err := parseData(tc.src); err != nil {
			t.Fatalf("Test: '%s' - %v", tc.src, err)
		} else if bytes.Compare(output, tc.expected) != 0 {
			t.Errorf("For '%s' expected %v, got %v", tc.src, tc.expected, output)
		}
	}
}
