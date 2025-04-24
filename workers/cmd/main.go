package main

import (
	"os"
	"strconv"

	a "app.workers/internal/app"
	"github.com/fatih/color"
)

func main() {
	args := os.Args[1:]

	amount := 3
	if len(args) > 0 {
		num, err := strconv.Atoi(args[0])
		if err == nil {
			amount = num
		}
	}

	app := a.NewApp(amount, "amqp://admeanie:shabi@localhost:5672/")
	if err := app.Run(); err != nil {
		color.Red("Error starting app: %s", err.Error())
		return
	}
}
