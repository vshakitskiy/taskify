package main

import (
	"fmt"

	a "app.workers/internal/app"
)

func main() {
	app, err := a.NewApp(3)
	if err != nil {
		fmt.Printf("Error creating app: %s", err.Error())
		return
	}

	app.Run()
}
