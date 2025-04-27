package api

import (
	"context"
	"fmt"

	"app.server/internal/config"
	"app.server/internal/service"
	"app.server/pkg/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TasksServiceHandler struct {
	pb.UnimplementedTasksServiceServer
	tasksService *service.TaskService
}

func NewTasksServiceHandler(tasksService *service.TaskService) *TasksServiceHandler {
	return &TasksServiceHandler{
		tasksService: tasksService,
	}
}

func (h *TasksServiceHandler) ExecuteTask(
	ctx context.Context,
	req *pb.ExecuteTaskRequest,
) (*pb.ExecuteTaskResponse, error) {
	if req.DurationSeconds < int64(config.Config.Task.MinDuration) || req.DurationSeconds > int64(config.Config.Task.MaxDuration) {
		return nil, status.Error(
			codes.InvalidArgument,
			fmt.Sprintf(
				"duration must be in a range of %d-%d seconds",
				config.Config.Task.MinDuration,
				config.Config.Task.MaxDuration,
			),
		)
	}

	if req.Message == "" {
		return nil, status.Error(codes.InvalidArgument, "message must not be empty")
	}

	res, err := h.tasksService.ExecuteTask(ctx, req.DurationSeconds, req.Message)
	if err != nil {
		switch err {
		case service.ErrExecutionFailed:
			return nil, status.Error(codes.FailedPrecondition, "execution failed")
		case service.InternalErrExecutionFailed:
			return nil, status.Error(codes.Internal, "failed to execute task")
		case service.ErrQueueIsFull:
			return nil, status.Error(codes.ResourceExhausted, "queue is full")
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	return &pb.ExecuteTaskResponse{
		Message:       res.Message,
		CorrelationId: res.CorelationId,
		WorkerNum:     res.WorkerNum,
	}, nil
}
