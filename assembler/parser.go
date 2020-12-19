package assembler

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

var (
	commentLineRe = regexp.MustCompile(`^;.*$`)
	opRe          = regexp.MustCompile(`^(\.?[A-Za-z]+)\s*(.*)$`)

	charRe = regexp.MustCompile(`'.'`)
	hex1Re = regexp.MustCompile(`^0[xX][0-9A-Fa-f]+$`)
	hex2Re = regexp.MustCompile(`^[0-9A-Fa-f]+[hH]$`)
)

func tokenize(s string) (op, args string, err error) {
	s = strings.TrimSpace(s)

	if s == "" || commentLineRe.MatchString(s) {
		return "", "", nil
	}

	tokens := opRe.FindStringSubmatch(s)
	if len(tokens) != 3 {
		return "", "", errors.Errorf("Invalid expression")
	}

	op = tokens[1]
	args = strings.TrimSpace(strings.SplitN(tokens[2], ";", 2)[0])
	return
}

func parseReg(s string) (reg byte, err error) {
	_, err = fmt.Sscanf(strings.TrimSpace(strings.ToLower(s)), "r%d", &reg)
	if err == nil && (reg < 0 || reg > 15) {
		return 0, errors.Errorf("Bad register: r%d", reg)
	}
	return
}

func parseRegReg(args string) (reg1, reg2 byte, err error) {
	tokens := strings.Split(args, ",")
	if len(tokens) != 2 {
		return 0, 0, errors.Errorf("Expected 2 arguments, ")
	}
	if reg1, err = parseReg(tokens[0]); err != nil {
		return 0, 0, err
	}
	if reg2, err = parseReg(tokens[1]); err != nil {
		return 0, 0, err
	}
	return
}

func parseRegRegOrImm8(args string) (reg1, reg2 byte, imm8 *byte, err error) {
	tokens := strings.SplitN(args, ",", 2)
	if len(tokens) != 2 {
		return 0, 0, nil, errors.Errorf("Expected 2 arguments, ")
	}
	if reg1, err = parseReg(tokens[0]); err != nil {
		return 0, 0, nil, err
	}

	if reg2, err = parseReg(tokens[1]); err == nil {
		// Success - "reg, reg"
		return
	}

	err = nil

	if b, err := parseImm8(tokens[1]); err != nil {
		return 0, 0, nil, err
	} else {
		// Success - "reg, <imm8>"
		imm8 = &b
	}

	return
}

func parseImm8(s string) (imm8 byte, err error) {
	s = strings.TrimSpace(s)

	switch {
	case charRe.MatchString(s):
		if _, err := fmt.Sscanf(s, "'%c'", &imm8); err != nil {
			return 0, err
		}

	case hex1Re.MatchString(s):
		if _, err := fmt.Sscanf(strings.ToLower(s), "0x%x", &imm8); err != nil {
			return 0, err
		}

	case hex2Re.MatchString(s):
		if _, err := fmt.Sscanf(strings.ToLower(s), "%xh", &imm8); err != nil {
			return 0, err
		}

	default:
		if _, err := fmt.Sscanf(s, "%d", &imm8); err != nil {
			return 0, err
		}
	}

	return
}

func parseImm12(s string) (uint16, error) {
	s = strings.TrimSpace(s)

	var addr uint16
	switch {
	case hex1Re.MatchString(s):
		if _, err := fmt.Sscanf(strings.ToLower(s), "0x%x", &addr); err != nil {
			return 0, err
		}
	case hex2Re.MatchString(s):
		if _, err := fmt.Sscanf(strings.ToLower(s), "%xh", &addr); err != nil {
			return 0, err
		}
	default:
		if _, err := fmt.Sscanf(s, "%d", &addr); err != nil {
			return 0, err
		}
	}

	if addr > 0xfff {
		return 0, errors.Errorf("Address out of range: 0x%x (%d)", addr, addr)
	}

	return addr, nil
}
