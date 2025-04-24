package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"app.workers/internal/mq"
	"app.workers/internal/workers"
	"github.com/fatih/color"
)

type App struct {
	mq     *mq.RabbitMQ
	w      *workers.Workers
	url    string
	amount int
}

func NewApp(workersAmount int, mqUrl string) *App {
	return &App{
		url:    mqUrl,
		amount: workersAmount,
	}
}

func (a *App) Run() error {
	color.Blue("Connecting to RabbitMQ...")
	mq, err := mq.NewMQ("amqp://admeanie:shabi@localhost:5672/", a.amount)
	if err != nil {
		return err
	}
	a.mq = mq
	color.Green("Connection to RabbitMQ established")
	fmt.Println()

	a.w = workers.NewWorkers(a.mq, a.amount)
	a.w.Run()

	time.Sleep(100 * time.Millisecond)
	color.Cyan("\nApplication started. Press Ctrl+C to exit.")
	fmt.Println()

	return a.gracefulShutdown()
}

func (a *App) gracefulShutdown() error {
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	<-shutdownChan
	color.Red("\n\nShutdown signal received...")
	fmt.Println()

	color.Blue("Waiting for workers to finish...")
	a.w.Stop()
	color.Green("Successfully closed all workers")
	fmt.Println()

	color.Blue("Closing RabbitMQ connections...")
	a.mq.Methods.Close()
	fmt.Println()

	color.Yellow("Application shutdown complete")

	return nil
}
