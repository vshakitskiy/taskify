package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"app.server/pkg/pb"

	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
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

func main() {
	flag.Parse()

	if *numTasks <= 0 {
		color.Red("Number of tasks must be greater than 0")
		return
	}

	if *minDuration < 1 {
		color.Yellow("Minimum duration must be not less than 1 second, adjusting to 1")
		*minDuration = 1
	}

	if *maxDuration > 50 {
		color.Yellow("Maximum duration must be not greater than 50 seconds, adjusting to 50")
		*maxDuration = 50
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	color.Blue("Connecting to gRPC server...")

	conn, err := grpc.NewClient(
		"localhost:50000",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		color.Red("Failed to connect to gRPC server: %v", err)
		return
	}
	defer conn.Close()

	color.Green("Connected to gRPC server")
	fmt.Println()

	client := pb.NewTasksServiceClient(conn)

	wg := sync.WaitGroup{}
	resCh := make(chan string, *numTasks)
	errCh := make(chan error, *numTasks)

	for i := 0; i < *numTasks; i++ {
		wg.Add(1)
		go execTaskRequest(client, r, &wg, i+1, resCh, errCh)
	}

	color.Cyan("Waiting for results...")
	wg.Wait()
	color.Green("All tasks goroutines finished")

	close(resCh)
	close(errCh)
	fmt.Println()

	successCount := 0
	errorCount := 0

	color.Yellow("Results:")

	for res := range resCh {
		color.Green("OK: %s", res)
		successCount++
	}

	for err := range errCh {
		color.Red("ERR: %s", err.Error())
		errorCount++
	}

	color.Yellow("Summary: %d tasks executed, %d successful, %d failed", *numTasks, successCount, errorCount)
}

func execTaskRequest(
	client pb.TasksServiceClient,
	r *rand.Rand,
	wg *sync.WaitGroup,
	taskNum int,
	resCh chan<- string,
	errCh chan<- error,
) {
	defer wg.Done()

	duration := int64(r.Intn(*maxDuration-*minDuration+1) + *minDuration)

	req := &pb.ExecuteTaskRequest{
		DurationSeconds: duration,
		Message:         fmt.Sprintf(*message, taskNum),
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(*maxDuration**numTasks)*time.Second,
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
