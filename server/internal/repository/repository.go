package repository

import (
	"sync"

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

func (r *TasksRepository) AddTaskResponse(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks[id] = make(chan models.TaskResponse, 1)
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
