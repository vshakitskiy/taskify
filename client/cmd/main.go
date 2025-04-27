package main

import (
	"app.client/internal/app"
	"app.client/internal/config"
	"github.com/fatih/color"
)

var reqConf *config.RequestsConfig = config.InitReqConfig()

func main() {
	a := app.NewApp()

	if err := a.Run(); err != nil {
		color.Red("Unable to run app: %s", err.Error())
		return
	}

	a.PrintResults()
}
