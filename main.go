package main

import (
	"encoding/hex"
	"io/ioutil"
	"log"
	"os"

	"github.com/julebokk/slede8dbg/debugger"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "slede8dbg",
		Usage: "A SLEDE8 debugger, etc, etc",
		Commands: []*cli.Command{
			{
				Name:      "debug",
				Usage:     "debug a SLEDE8 binary",
				UsageText: "slede8dbg debug [options] <path to SLEDE8 binary>",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "input, i",
						Usage: "hexadecimal input string (AKA SLEDE8 føde), e.g. CD21",
					},
				},
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						return cli.NewExitError("Binary path is missing", 1)
					}

					input, err := hex.DecodeString(c.String("input"))
					if err != nil {
						return err
					}

					program, err := ioutil.ReadFile(c.Args().First())
					if err != nil {
						return err
					}

					debugger, err := debugger.NewUI(program, input)
					if err != nil {
						return err
					}

					return debugger.MainLoop()
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
