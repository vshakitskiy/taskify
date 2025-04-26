package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"app.server/internal/api"
	"app.server/internal/mq"
	"app.server/internal/repository"
	"app.server/internal/service"
	"app.server/pkg/pb"
	"app.shared/models"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	grpcServer *grpc.Server
	mq         *mq.RabbitMQ
	port       string
}

func NewApp() *App {
	return &App{
		port: "50000",
	}
}

func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	color.Blue("Connecting to RabbitMQ...")
	mq, err := mq.NewMQ("amqp://admeanie:shabi@localhost:5672/", 3)
	if err != nil {
		return err
	}
	color.Green("Connection to RabbitMQ established")
	fmt.Println()
	a.mq = mq

	tasksRepository := repository.NewTasksRepository()
	tasksService := service.NewTaskService(mq, tasksRepository)
	tasksHandler := api.NewTasksServiceHandler(tasksService)

	a.grpcServer = grpc.NewServer()

	pb.RegisterTasksServiceServer(a.grpcServer, tasksHandler)

	reflection.Register(a.grpcServer)

	go func() {
		color.Blue("Listening for results queue...")
		fmt.Println()

		select {
		case <-ctx.Done():
			color.Blue("Closing results queue listener...")
			break
		case del := <-mq.ResultsCh:
			resCh, ok := tasksRepository.GetTaskResponse(del.CorrelationId)
			if !ok {
				fmt.Println("!", del.CorrelationId, string(del.Body))
				fmt.Println("not found")
				del.Ack(false)
			} else {
				taskResponse := models.TaskResponse{}
				json.Unmarshal(del.Body, &taskResponse)

				resCh <- taskResponse
				del.Ack(false)
			}
		}
	}()

	go func() {
		color.Blue("Starting grpc server...")
		lis, err := net.Listen("tcp", ":"+a.port)
		if err != nil {
			color.Red("Unable to start tcp listener: %s", err)
			return
		}

		color.Green("Starting application. Press Ctrl+C to exit.")
		fmt.Println()
		if err := a.grpcServer.Serve(lis); err != nil {
			color.Red("Unable to serve grpc: %s", err)
		}
	}()

	return a.gracefulShutdown()
}

func (a *App) gracefulShutdown() error {
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)
	<-shutdownChan
	color.Red("\n\nShutdown signal received...")
	fmt.Println()

	color.Blue("Closing RabbitMQ connections...")
	a.mq.Methods.Close()
	fmt.Println()

	color.Yellow("Application shutdown complete")

	return nil
}
