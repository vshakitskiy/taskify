package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	color.Blue("Connecting to RabbitMQ...")
	mq, err := mq.NewMQ("amqp://admeanie:shabi@localhost:5672/", 3)
	if err != nil {
		cancel()
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

	go func(rCtx context.Context) {
		color.Blue("Listening for results queue...")
		fmt.Println()

		for {
			select {
			case <-rCtx.Done():
				color.Blue("Closing results queue listener...")
				return
			case del, ok := <-mq.ResultsCh:
				if !ok {
					color.Red("RESULTS QUEUE LISTENER FETCHED FROM CLOSED CHANNEL")
					return
				}

				resCh, found := tasksRepository.GetTaskResponse(del.CorrelationId)
				if !found {
					color.Red("NOT FOUND: %s, %s", del.CorrelationId, string(del.Body))
					del.Ack(false)
					continue
				}

				taskResponse := models.TaskResponse{}
				if err := json.Unmarshal(del.Body, &taskResponse); err != nil {
					fmt.Println("JSON UNMARSHAL ERROR", del.CorrelationId, string(del.Body))
					del.Ack(false)
					continue
				}

				resCh <- taskResponse
				del.Ack(false)
			}
		}
	}(ctx)

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

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)
	<-shutdownChan
	color.Red("\n\nShutdown signal received...")
	fmt.Println()

	color.Blue("Stopping gRPC server...")
	stopped := make(chan struct{})

	go func() {
		a.grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		color.Green("gRPC server stopped")
		fmt.Println()
	case <-time.After(15 * time.Second):
		color.Red("gRPC server shutdown timed out. Forcing stop...")
		fmt.Println()
		a.grpcServer.Stop()
	}

	cancel()

	color.Blue("Closing RabbitMQ connections...")
	a.mq.Methods.Close()
	fmt.Println()

	color.Yellow("Application shutdown complete")
	return nil
}
