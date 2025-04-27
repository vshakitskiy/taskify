package service

import (
	"context"
	"fmt"
	"time"

	"app.server/internal/mq"
	"app.server/internal/repository"
	"app.shared/models"
	"github.com/google/uuid"
)

var (
	ErrExecutionFailed         = fmt.Errorf("execution failed")
	InternalErrExecutionFailed = fmt.Errorf("unable to execute task")
	replyTo                    = "results"
)

type TaskService struct {
	mq   *mq.RabbitMQ
	repo *repository.TasksRepository
}

func NewTaskService(mq *mq.RabbitMQ, repo *repository.TasksRepository) *TaskService {
	return &TaskService{
		mq:   mq,
		repo: repo,
	}
}

func (s *TaskService) ExecuteTask(
	ctx context.Context,
	duration int64,
	message string,
) (*ExecutionResult, error) {
	task := models.NewTask(message, int(duration))
	correlationID := uuid.New().String()

	s.repo.AddTaskResponse(correlationID)
	defer s.repo.DeleteTaskResponse(correlationID)

	if err := s.mq.Methods.Publish(
		ctx,
		s.mq.TaskQueue,
		task.ToJSON(),
		&correlationID,
		&replyTo,
	); err != nil {
		return nil, InternalErrExecutionFailed
	}

	resCh, _ := s.repo.GetTaskResponse(correlationID)

	select {
	case res := <-resCh:
		if res.ErrorCode != "" {
			return nil, ErrExecutionFailed
		}

		return &ExecutionResult{
			Message:      res.Message,
			CorelationId: correlationID,
			WorkerNum:    res.WorkerNum,
		}, nil
	case <-time.After(time.Duration(duration+2) * time.Second):
		return nil, InternalErrExecutionFailed
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

type ExecutionResult struct {
	Message      string
	CorelationId string
	WorkerNum    int64
}
