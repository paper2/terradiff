package main

import "github.com/urfave/cli/v2"

// TODO: https://github.com/suzuki-shunsuke/tfcmt/blob/main/pkg/cli/app.go みたいにする
func NewApp() *cli.App {
	return &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "branch",
				Value: "main",
				Usage: "language for the greeting",
			},
		},
		Action: NewActions().Terradiff,
	}
}
