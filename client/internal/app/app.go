package app

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"app.server/pkg/pb"

	"app.client/internal/config"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type App struct {
	reqConf *config.RequestsConfig
	resCh   chan string
	errorCh chan error
}

func NewApp() (*App, error) {
	if err := config.InitConfig(); err != nil {
		return nil, err
	}

	reqConfig := config.InitReqConfig()
	return &App{
		reqConf: reqConfig,
		resCh:   make(chan string, reqConfig.NumTasks),
		errorCh: make(chan error, reqConfig.NumTasks),
	}, nil
}

func (a *App) Run() error {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	color.Blue("Creating gRPC server's connection...")
	fmt.Println()

	conn, err := grpc.NewClient(
		"localhost:"+string(config.Config.Server.GRPCPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewTasksServiceClient(conn)

	wg := sync.WaitGroup{}
	for i := 0; i < a.reqConf.NumTasks; i++ {
		wg.Add(1)
		go a.execTaskRequest(client, r, &wg, i+1, a.resCh, a.errorCh)
	}

	color.Cyan("Waiting for results...")
	wg.Wait()
	color.Green("All tasks goroutines finished")
	fmt.Println()

	close(a.resCh)
	close(a.errorCh)

	return nil
}

func (a *App) PrintResults() {
	successCount := 0
	errorCount := 0

	color.Yellow("Results:")

	for res := range a.resCh {
		color.Green("OK: %s", res)
		successCount++
	}

	for err := range a.errorCh {
		color.Red("ERR: %s", err.Error())
		errorCount++
	}

	color.Yellow("Summary: %d tasks executed, %d successful, %d failed", a.reqConf.NumTasks, successCount, errorCount)
}

func (a *App) execTaskRequest(
	client pb.TasksServiceClient,
	r *rand.Rand,
	wg *sync.WaitGroup,
	taskNum int,
	resCh chan<- string,
	errCh chan<- error,
) {
	defer wg.Done()

	duration := int64(r.Intn(a.reqConf.MaxDuration-a.reqConf.MinDuration+1) + a.reqConf.MinDuration)

	req := &pb.ExecuteTaskRequest{
		DurationSeconds: duration,
		Message:         strings.Replace(a.reqConf.MessageFormat, "%d", fmt.Sprint(taskNum), 1),
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(a.reqConf.MaxDuration*a.reqConf.NumTasks)*time.Second,
	)
	defer cancel()

	res, err := client.ExecuteTask(ctx, req)
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			errCh <- fmt.Errorf(
				"(Task %d) gRPC error: code=%s, message=%s",
				taskNum,
				s.Code(),
				s.Message(),
			)
			return
		}

		errCh <- fmt.Errorf("(Task %d) non-gRPC error: %w", taskNum, err)
		return
	}

	resCh <- fmt.Sprintf(
		"(Task %d) Completed: correlationID=%s, workerNum=%d, duration=%d",
		taskNum,
		res.CorrelationId,
		res.WorkerNum,
		duration,
	)
}
