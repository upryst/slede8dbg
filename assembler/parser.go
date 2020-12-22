package assembler

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

var (
	commentLineRe = regexp.MustCompile(`^;.*$`)
	opRe          = regexp.MustCompile(`^(\.?[A-Za-z]+)\s*(.*)$`)

	labelDefRe = regexp.MustCompile(`^\s*([A-Za-zÆØÅæøå_][A-Za-z0-9ÆØÅæøå_]*):\s*$`)
	labelRe    = regexp.MustCompile(`^([A-Za-zÆØÅæøå_][A-Za-z0-9ÆØÅæøå_]*)\s*$`)

	charRe    = regexp.MustCompile(`'.'`)
	hex1Re    = regexp.MustCompile(`^0[xX][0-9A-Fa-f]+$`)
	hex2Re    = regexp.MustCompile(`^[0-9A-Fa-f]+[hH]$`)
	decimalRe = regexp.MustCompile(`^[+-]?[0-9]+$`)

	dataHex1Re    = regexp.MustCompile(`^(0[xX][0-9A-F-a-f]+\s*)`)
	dataHex2Re    = regexp.MustCompile(`^([0-9A-F-a-f]+[hH]\s*)`)
	dataDecimalRe = regexp.MustCompile(`^([+-]?[0-9]+\s*)`)
	dataCharRe    = regexp.MustCompile(`^('.'\s*)`)
)

func tokenize(s string) (label, op, args string, err error) {
	s = strings.TrimSpace(s)

	if s == "" || commentLineRe.MatchString(s) {
		return "", "", "", nil
	}

	labels := labelDefRe.FindStringSubmatch(s)
	if len(labels) == 2 {
		return labels[1], "", "", nil
	}

	tokens := opRe.FindStringSubmatch(s)
	if len(tokens) != 3 {
		return "", "", "", errors.Errorf("Invalid expression")
	}

	op = tokens[1]

	// Hacky af, but .DATA with semicolons in strings makes things go "boom",
	// .DATA parser handles semicolons itself.
	if strings.ToUpper(op) != ".DATA" {
		args = strings.TrimSpace(strings.SplitN(tokens[2], ";", 2)[0])
	} else {
		args = strings.TrimSpace(tokens[2])
	}
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

func parseImm12OrLabel(s string) (uint16, string, error) {
	s = strings.TrimSpace(s)

	var addr uint16
	switch {
	case hex1Re.MatchString(s):
		if _, err := fmt.Sscanf(strings.ToLower(s), "0x%x", &addr); err != nil {
			return 0, "", err
		}
	case hex2Re.MatchString(s):
		if _, err := fmt.Sscanf(strings.ToLower(s), "%xh", &addr); err != nil {
			return 0, "", err
		}
	case labelRe.MatchString(s):
		return 0, s, nil
	default:
		if _, err := fmt.Sscanf(s, "%d", &addr); err != nil {
			return 0, "", err
		}
	}

	if addr > 0xfff {
		return 0, "", errors.Errorf("Address out of range: 0x%x (%d)", addr, addr)
	}

	return addr, "", nil
}

func parseData(args string) ([]byte, error) {
	var output bytes.Buffer

	s := strings.TrimSpace(args)

	for s != "" {
		for s[0] == ' ' || s[0] == '\t' {
			s = s[1:]
		}

		if s[0] == '"' {
			s = s[1:]
			for s != "" && s[0] != '"' {
				if s[0] == '\\' {
					if len(s) == 1 {
						return nil, errors.Errorf("Bad escape sequence")
					}
					switch s[1] {
					case 'n':
						output.WriteByte('\n')
					case 't':
						output.WriteByte('\t')
					case 'r':
						output.WriteByte('\r')
					default:
						output.WriteByte(s[1])
					}
					s = s[2:]
				} else {
					output.WriteByte(s[0])
					s = s[1:]
				}
			}
			// Consume the double quote
			if s == "" {
				return nil, errors.Errorf("Missing closing double quote")
			}
			s = s[1:]
		} else if match := dataHex1Re.FindStringSubmatch(s); len(match) == 2 {
			var b byte
			if _, err := fmt.Sscanf(strings.ToLower(s), "0x%x", &b); err != nil {
				return nil, err
			} else {
				output.WriteByte(b)
				s = s[len(match[1]):]
			}
		} else if match := dataHex2Re.FindStringSubmatch(s); len(match) == 2 {
			var b byte
			if _, err := fmt.Sscanf(strings.ToLower(s), "%xh", &b); err != nil {
				return nil, err
			} else {
				output.WriteByte(b)
				s = s[len(match[1]):]
			}
		} else if match := dataCharRe.FindStringSubmatch(s); len(match) == 2 {
			var b byte
			if _, err := fmt.Sscanf(s, "'%c'", &b); err != nil {
				return nil, err
			} else {
				output.WriteByte(b)
				s = s[len(match[1]):]
			}
		} else if match := dataDecimalRe.FindStringSubmatch(s); len(match) == 2 {
			var b byte
			if _, err := fmt.Sscanf(strings.ToLower(s), "%d", &b); err != nil {
				return nil, err
			} else {
				output.WriteByte(b)
				s = s[len(match[1]):]
			}
		} else {
			return nil, errors.Errorf("Unexpected .DATA sequence: %s", s)
		}

		// Skip whitespace
		for s != "" && (s[0] == ' ' || s[0] == '\t') {
			s = s[1:]
		}

		// Skip comma / handle trailing comments
		if s != "" {
			if s[0] == ';' {
				break
			} else if s[0] != ',' {
				return nil, errors.Errorf("Expected comma")
			}
			s = s[1:]
		}

	}

	return output.Bytes(), nil
}
