package assembler

import (
	"bytes"
	"strings"

	"github.com/pkg/errors"

	"github.com/upryst/slede8dbg/vm"
)

func AssembleLine(line string) ([]byte, error) {
	label, mnemonic, args, err := tokenize(line)
	if err != nil {
		return nil, err
	}

	if label != "" {
		return nil, errors.Errorf("Labels are not supported in single line mode")
	}

	return assemble(singleLineMode, nil, mnemonic, args)
}

func Assemble(src string) ([]byte, error) {
	// First pass, collect labels
	labels := make(map[string]uint16)
	var offset uint16
	for i, line := range strings.Split(src, "\n") {
		label, mnemonic, args, err := tokenize(line)
		if err != nil {
			return nil, errors.Errorf("Line %d: %v", i+1, err)
		}

		if label != "" {
			labels[label] = offset
			continue
		}

		bytecode, err := assemble(multilineFirstPass, nil, mnemonic, args)
		if err != nil {
			return nil, errors.Errorf("Line %d: %v", i+1, err)
		}

		offset += uint16(len(bytecode))

		if offset >= vm.MemSize {
			return nil, errors.Errorf("Program doesn't fit %d bytes", vm.MemSize)
		}
	}

	var output bytes.Buffer

	// Second (and final) pass with known label addresses
	for i, line := range strings.Split(src, "\n") {
		label, mnemonic, args, err := tokenize(line)
		if err != nil {
			return nil, err
		}

		if label != "" {
			continue
		}

		bytecode, err := assemble(multilineFinalPass, labels, mnemonic, args)
		if err != nil {
			return nil, errors.Errorf("Line %d: %v", i+1, err)
		}

		output.Write(bytecode)
	}

	return output.Bytes(), nil
}
