package repository

import (
	"fmt"
	"sync"

	"app.server/internal/config"
	"app.shared/models"
)

type TasksRepository struct {
	tasks map[string]chan models.TaskResponse
	mu    sync.Mutex
}

func NewTasksRepository() *TasksRepository {
	return &TasksRepository{
		tasks: make(map[string]chan models.TaskResponse),
	}
}

func (r *TasksRepository) AddTaskResponse(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.tasks) >= config.Config.Server.QueueSize {
		return fmt.Errorf("too many tasks")
	}
	r.tasks[id] = make(chan models.TaskResponse, 1)
	return nil
}

func (r *TasksRepository) GetTaskResponse(id string) (chan models.TaskResponse, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	ch, ok := r.tasks[id]
	return ch, ok
}
func (r *TasksRepository) DeleteTaskResponse(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.tasks, id)
}

func (r *TasksRepository) Size() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.tasks)
}
