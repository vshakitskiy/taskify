package config

import (
	"flag"

	"github.com/fatih/color"
)

var (
	numTasks = flag.Int(
		"n",
		10,
		"Number of tasks to execute",
	)
	minDuration = flag.Int(
		"min",
		3,
		"Minimum duration of task in seconds",
	)
	maxDuration = flag.Int(
		"max",
		10,
		"Maximum duration of task in seconds",
	)
	message = flag.String(
		"m",
		"Hello, world n. %d!",
		"Message payload for the task. Use %d to substitute task number",
	)
)

type RequestsConfig struct {
	NumTasks      int
	MinDuration   int
	MaxDuration   int
	MessageFormat string
}

func InitReqConfig() *RequestsConfig {
	flag.Parse()

	if *numTasks <= 0 {
		color.Red("Number of tasks must be greater than 0; adjusting to 10")
	}

	if *minDuration < 1 {
		color.Yellow("Minimum duration must be not less than 1 second; adjusting to 1")
		*minDuration = 1
	}

	if *maxDuration > 30 {
		color.Yellow("Maximum duration must be not greater than 30 seconds; adjusting to 30")
		*maxDuration = 30
	}

	return &RequestsConfig{
		NumTasks:      *numTasks,
		MinDuration:   *minDuration,
		MaxDuration:   *maxDuration,
		MessageFormat: *message,
	}
}
