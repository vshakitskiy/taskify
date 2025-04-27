package main

import (
	"app.client/internal/app"
	"github.com/fatih/color"
)

func main() {
	a, err := app.NewApp()
	if err != nil {
		color.Red("Unable to create app: %s", err.Error())
		return
	}

	if err = a.Run(); err != nil {
		color.Red("Unable to run app: %s", err.Error())
		return
	}

	a.PrintResults()
}
