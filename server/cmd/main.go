package main

import (
	"app.server/internal/app"
	"app.server/internal/config"
	"github.com/fatih/color"
)

func main() {
	if err := config.InitConfig(); err != nil {
		color.Red("Unable to init config: %s", err)
		return
	}

	a := app.NewApp()

	if err := a.Run(); err != nil {
		color.Red("Unable to start app: %s", err)
	}
}
