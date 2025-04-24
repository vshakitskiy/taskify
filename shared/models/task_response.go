package models

import (
	"encoding/json"
)

type TaskResponse struct {
	Message   string `json:"message,omitempty"`
	Suceess   bool   `json:"success,required"`
	ErrorCode string `json:"error_code,omitempty"`
}

func NewSuccessTaskResponse(message string) *TaskResponse {
	return &TaskResponse{
		Message: message,
		Suceess: true,
	}
}

func NewErrorTaskResponse(errCode string) *TaskResponse {
	return &TaskResponse{
		Suceess:   false,
		ErrorCode: errCode,
	}
}

func (t *TaskResponse) ToJSON() []byte {
	b, _ := json.Marshal(t)
	return b
}
