package models

import (
	"encoding/json"
	"fmt"
)

type Task struct {
	Timeout int    `json:"timeout,required"`
	Message string `json:"message,required"`
}

func NewTask(message string, timeout int) *Task {
	return &Task{
		Timeout: timeout,
		Message: message,
	}
}

func TaskFromJSON(b []byte) (*Task, error) {
	task := new(Task)
	err := json.Unmarshal(b, task)
	if err != nil {
		return nil, err
	}

	if task.Timeout < 100 {
		return nil, fmt.Errorf("Invalid timeout")
	}

	if task.Message == "" {
		return nil, fmt.Errorf("Invalid message")
	}

	return task, nil
}
