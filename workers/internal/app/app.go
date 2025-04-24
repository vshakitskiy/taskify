package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"app.workers/internal/mq"
	"app.workers/internal/workers"
)

type App struct {
	mq *mq.RabbitMQ
	w  *workers.Workers
}

func NewApp(amount int) (*App, error) {
	mq, err := mq.NewMQ("amqp://admeanie:shabi@localhost:5672/", amount)
	if err != nil {
		return nil, err
	}

	return &App{
		mq: mq,
		w:  workers.NewWorkers(mq, amount),
	}, nil
}

func (a *App) Run() {
	a.w.Run()
	time.Sleep(500 * time.Millisecond)

	fmt.Println("Application started. Press Ctrl+C to exit.")

	a.gracefulShutdown()
}

func (a *App) gracefulShutdown() {
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	<-shutdownChan
	fmt.Println("\nShutdown signal received...")

	fmt.Println("Waiting for workers to finish...")
	a.w.Stop()

	fmt.Println("Closing RabbitMQ connections...")
	a.mq.Methods.Close()

	fmt.Println("Application shutdown complete")
}
