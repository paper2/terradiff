package main

import "github.com/urfave/cli/v2"

const (
	workDirFlag      = "work-dir"
	srcBranchFlag    = "source-branch"
	dstBranchFlag    = "destination-branch"
	repoURLFlag      = "repository-url"
	debugFlag        = "debug"
	jsonFlag         = "json-log"
	terraformDirFlag = "terraform-dir"
)

func NewApp() *cli.App {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    workDirFlag,
			Value:   "./terradiff-workspace",
			Usage:   "work directory",
			Aliases: []string{"w"},
		},
		&cli.StringFlag{
			Required: true,
			Name:     srcBranchFlag,
			Usage:    "source branch name",
			Aliases:  []string{"s"},
		},
		&cli.StringFlag{
			Name:    dstBranchFlag,
			Value:   "main",
			Usage:   "destination branch name",
			Aliases: []string{"d"},
		},
		&cli.StringFlag{
			Required: true,
			Name:     repoURLFlag,
			Usage:    "repository url",
			Aliases:  []string{"r"},
		},
		&cli.BoolFlag{
			Name:  debugFlag,
			Value: false,
			Usage: "debug log",
		},
		&cli.BoolFlag{
			Name:  jsonFlag,
			Value: false,
			Usage: "json log",
		},
		&cli.StringFlag{
			Name:  terraformDirFlag,
			Value: ".",
			Usage: "terraform root directory path on repository",
		},
	}
	app.Action = TerradiffAction
	return app
}
