package main

import (
	"github.com/urfave/cli/v2"
	"gitlab.kilic.dev/docker/beamer/pipe"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

func main() {
	NewPlumber(
		func(p *Plumber) *cli.App {
			return &cli.App{
				Name:        CLI_NAME,
				Version:     VERSION,
				Usage:       DESCRIPTION,
				Description: DESCRIPTION,
				Flags:       p.AppendFlags(pipe.Flags),
				Action: func(ctx *cli.Context) error {
					return pipe.TL.RunJobs(
						pipe.New(p).SetCliContext(ctx).Job(),
					)
				},
			}
		}).
		SetDocumentationOptions(DocumentationOptions{
			MarkdownOutputFile: "CLI.md",
			MarkdownBehead:     0,
			ExcludeFlags:       true,
			ExcludeHelpCommand: true,
		}).
		Run()
}
