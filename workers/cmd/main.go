package main

import (
	"app.shared/pkg/env"
	a "app.workers/internal/app"
	"app.workers/internal/config"
	"github.com/fatih/color"
)

func main() {
	if err := config.InitConfig(); err != nil {
		color.Red("Unable to init config: %s", err)
		return
	}

	amount := config.Config.Workers.Amount
	app := a.NewApp(
		amount, env.GetDefaultEnv(
			"RABBITMQ_URL",
			"amqp://admeanie:shabi@localhost:5672/",
		))
	if err := app.Run(); err != nil {
		color.Red("Error starting app: %s", err.Error())
		return
	}
}
