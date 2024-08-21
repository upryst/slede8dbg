package main

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/upryst/slede8dbg/assembler"
	"github.com/upryst/slede8dbg/debugger"
	"github.com/upryst/slede8dbg/vm"

	"github.com/urfave/cli/v2"
)

const (
	defaultCycleLimit = 50000

	asmExtension = ".asm"
)

func compileAsmFile(path string) ([]byte, error) {
	source, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	bytecode, err := assembler.Assemble(string(source))
	if err != nil {
		return nil, err
	}

	var binary bytes.Buffer
	binary.Write([]byte(vm.SledeHeader))
	binary.Write(bytecode)

	return binary.Bytes(), nil
}

func debug(path, inputStr string, cycleLimit int) error {
	input, err := hex.DecodeString(inputStr)
	if err != nil {
		return err
	}

	var binary []byte

	if filepath.Ext(path) == asmExtension {
		if binary, err = compileAsmFile(path); err != nil {
			return err
		}
	} else if binary, err = ioutil.ReadFile(path); err != nil {
		return err
	}

	debugger, err := debugger.NewUI(binary, input, cycleLimit)
	if err != nil {
		return err
	}

	return debugger.MainLoop()
}

func main() {
	app := &cli.App{
		Name:  "slede8dbg",
		Usage: "A SLEDE8 debugger/assembler",
		UsageText: "slede8dbg command [command options] [arguments...]\n" +
			"   or\n" +
			"   slede8dbg <path to .asm/.s8 path> [<input> [cycle limit]]",
	}

	app.UseShortOptionHandling = true
	app.Commands = []*cli.Command{
		{
			Name:      "debug",
			Aliases:   []string{"d"},
			Usage:     "debug a SLEDE8 binary",
			UsageText: "slede8dbg debug [options] <path to SLEDE8 binary / ASM source>",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "input",
					Aliases: []string{"i"},
					Usage:   "hexadecimal input string (AKA SLEDE8 føde), e.g. CD21",
				},
				&cli.IntFlag{
					Name:    "limit",
					Aliases: []string{"l"},
					Usage:   "cycle (step) limit",
					Value:   defaultCycleLimit,
				},
			},
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					return cli.NewExitError(".s8 / .asm path is missing", 1)
				}

				return debug(c.Args().First(), c.String("input"), c.Int("limit"))
			},
		},
		{
			Name:    "compile",
			Aliases: []string{"c"},
			Usage:   "compile SLEDE8 asm source",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "output",
					Aliases: []string{"o"},
					Usage:   "output file path",
					Value:   "a.s8",
				},
			},
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					return cli.NewExitError("Source path is missing", 1)
				}

				if binary, err := compileAsmFile(c.Args().First()); err != nil {
					return err
				} else {
					return ioutil.WriteFile(c.String("output"), binary, 0644)
				}
			},
		},
	}

	// Alternative syntax (slede8dbg <path> [<input> [cycle limit]])
	app.Action = func(c *cli.Context) error {
		if c.NArg() == 0 {
			cli.ShowAppHelpAndExit(c, 0)
		}

		var err error
		cycleLimit := defaultCycleLimit

		if c.NArg() > 2 {
			if cycleLimit, err = strconv.Atoi(c.Args().Get(2)); err != nil {
				return err
			}
		}

		return debug(c.Args().Get(0), c.Args().Get(1), cycleLimit)
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
