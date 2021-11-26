package main

import (
	"log"
	"os"

	"github.com/uga-rosa/umx/internal/adsp"
	"github.com/uga-rosa/umx/internal/hist"
	"github.com/uga-rosa/umx/internal/pp"
	"github.com/uga-rosa/umx/internal/rg"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "umx",
		Usage: "To analyze the gromacs data.",
		Commands: []*cli.Command{
			{
				Name:  "pp",
				Usage: "Putting the data into json",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "input",
						Aliases:  []string{"f"},
						Value:    "",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "pego",
						Aliases: []string{"p"},
						Value:   "peg_o",
					},
					&cli.StringFlag{
						Name:    "rg",
						Aliases: []string{"r"},
						Value:   "rg",
					},
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Value:   "pp.json",
					},
				},
				Action: pp.Cmd,
			},
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
					&cli.Int64Flag{
						Name:    "number",
						Aliases: []string{"n"},
						Value:   1,
					},
				},
				Action: adsp.Cmd,
			},
			{
				Name:  "rg",
				Usage: "Divide the radius of gyration by the presence or absence of adsorption and draw histograms.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "input",
						Aliases:  []string{"f"},
						Value:    "umx/pp.json",
						Required: true,
					},
					&cli.Int64Flag{
						Name:    "number",
						Aliases: []string{"n"},
						Value:   1,
					},
				},
				Action: rg.Cmd,
			},
			{
				Name:  "hist",
				Usage: "Draw a histogram from the merged Rg json file.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "input",
						Aliases:  []string{"f"},
						Value:    "",
						Required: true,
					},
					&cli.Int64Flag{
						Name:    "number",
						Aliases: []string{"n"},
						Value:   1,
					},
				},
				Action: hist.Cmd,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
