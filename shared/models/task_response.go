package models

import (
	"encoding/json"
)

type TaskResponse struct {
	Message   string `json:"message,omitempty"`
	Suceess   bool   `json:"success,required"`
	WorkerNum int64  `json:"worker_num,required"`
	ErrorCode string `json:"error_code,omitempty"`
}

func NewSuccessTaskResponse(num int, message string) *TaskResponse {
	return &TaskResponse{
		WorkerNum: int64(num),
		Message:   message,
		Suceess:   true,
	}
}

func NewErrorTaskResponse(num int, errCode string) *TaskResponse {
	return &TaskResponse{
		WorkerNum: int64(num),
		Suceess:   false,
		ErrorCode: errCode,
	}
}

func (t *TaskResponse) ToJSON() []byte {
	b, _ := json.Marshal(t)
	return b
}
