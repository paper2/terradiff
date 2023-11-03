package main

import (
	"fmt"
	"os"
)

func main() {
	app := NewApp()
	if err := app.Run(os.Args); err != nil {
		Logger().Error(fmt.Sprintf("%+v", err))
	}
}
