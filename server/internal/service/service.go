package service

import (
	"context"
	"fmt"

	"app.server/internal/mq"
	"app.server/internal/repository"
	"app.shared/models"
	"github.com/google/uuid"
)

var (
	ErrInvalidDuration = fmt.Errorf("invalid duration: duration must not be less than 1s")
	ErrInvalidMessage  = fmt.Errorf("invalid message: message must not be empty")
	ErrPublishFailed   = fmt.Errorf("publish failed")
	replyTo            = "results"
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
	fmt.Println(s.mq.TaskQueue.Consumers)

	task := models.NewTask(message, int(duration))
	id := uuid.New().String()
	if err := s.mq.Methods.Publish(
		ctx,
		s.mq.TaskQueue,
		task.ToJSON(),
		&id,
		&replyTo,
	); err != nil {
		return nil, ErrPublishFailed
	}
	s.repo.AddTaskResponse(id)

	ch, _ := s.repo.GetTaskResponse(id)
	res := <-ch

	fmt.Println(res)

	return &ExecutionResult{}, nil
}

type ExecutionResult struct {
	Message      string
	CorelationId string
	WorkerNum    int64
}
