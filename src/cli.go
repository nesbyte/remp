package src

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func SetupCli(ver string) *cli.App {

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "v",
		Aliases: []string{"V", "version"},
		Usage:   "Shows version",
	}

	cli.VersionPrinter = func(cCtx *cli.Context) {
		fmt.Printf("%s\n", cCtx.App.Version)
	}

	// Change the order of how h and help messages
	// are shown (from --help -h to -h --help)
	cli.HelpFlag = &cli.BoolFlag{
		Name:    "h",
		Aliases: []string{"help"},
		Usage:   "Shows this help message",
	}

	var i = Input{}

	// var flagRegExFile string
	app := &cli.App{
		Name:                   "remp",
		Usage:                  "find directory from provided path matching the regex pattern",
		HideHelpCommand:        true,
		EnableBashCompletion:   true,
		UseShortOptionHandling: true,
		Version:                ver,

		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:        "e",
				Aliases:     []string{"regexp"},
				Usage:       "The regex pattern to search for",
				Destination: &i.RegexFlag,
			},
			&cli.StringSliceFlag{
				Name:        "X",
				Aliases:     []string{"line-strings"},
				Usage:       "File/directory name must exactly match string",
				Destination: &i.LineStrings,
			},
			&cli.StringSliceFlag{
				Name:        "f",
				Aliases:     []string{"file"},
				Usage:       "Provide regex pattern(s) from file(s). One per line",
				Destination: &i.RegexFiles,
			},
			&cli.BoolFlag{
				Name:        "l",
				Aliases:     []string{"left"},
				Value:       false,
				Usage:       "Search path from left to right instead of right to left",
				Destination: &i.SearchFromLeft,
			},
			&cli.BoolFlag{
				Name:        "b",
				Aliases:     []string{"base-directory"},
				Value:       false,
				Usage:       "Show the base directory without the matched file/directory",
				Destination: &i.ShowBaseDirectory,
			},
			&cli.BoolFlag{
				Name:        "a",
				Aliases:     []string{"match-all"},
				Value:       false,
				Usage:       "Search entire path instead of exiting on the first match",
				Destination: &i.MatchAll,
			},
			&cli.StringFlag{
				Name:        "O",
				Aliases:     []string{"no-match"},
				Value:       "",
				Usage:       "Custom stdout output when no match is found",
				Destination: &i.OutputWhenNoMatch,
			},
			&cli.BoolFlag{
				Name:        "color",
				Value:       false,
				Usage:       "Highlight matched term(s) with colour",
				Destination: &i.ShowColour,
			},
		},
		Action: func(ctx *cli.Context) error {
			i.RegexArg = ctx.Args().First()

			var b = make([]byte, 4096)
			n, err := os.Stdin.Read(b)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			i.PathFromStdin = filepath.ToSlash(string(b[:n]))

			res, err := FindMatchesInPath(i)
			if err != nil {
				if err == ErrNoMatch {
					if i.OutputWhenNoMatch != "" {
						fmt.Println(i.OutputWhenNoMatch)
					}
					return cli.Exit(ErrNoMatch, 1)
				}
				return cli.Exit(err.Error(), 1)
			}

			for _, r := range res {
				fmt.Println(filepath.FromSlash(r))
			}

			return nil
		},
	}

	cli.AppHelpTemplate = `USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
{{if .VisibleFlags}}
OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
`

	return app
}
