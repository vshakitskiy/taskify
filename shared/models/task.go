package models

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	Id      uuid.UUID     `json:"id"`
	Timeout time.Duration `json:"timeout"`
	Message string        `json:"message"`
}

func NewTask(message string, timeout time.Duration) *Task {
	return &Task{
		Id:      uuid.New(),
		Timeout: timeout,
		Message: message,
	}
}
