package main

import (
	"log"
	"os"

	"github.com/uga-rosa/umx/internal/adsp"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "umx",
		Usage: "To analyze the gromacs data.",
		Commands: []*cli.Command{
			{
				Name:  "adsp",
				Usage: "Calculate the adsorption probability and output a time transition graph of the number of adsorptions.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "input",
						Aliases:  []string{"f"},
						Value:    "",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "directory",
						Aliases: []string{"d"},
						Value:   "peg_o",
					},
					&cli.Int64Flag{
						Name:    "number",
						Aliases: []string{"n"},
						Value:   1,
					},
				},
				Action: adsp.Cmd,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
