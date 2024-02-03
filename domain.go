package main

import (
	"time"
)

const (
	Pending    status = "Pending"
	InProgress status = "InProgress"
	Done       status = "Done"
)

type status string

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      status `json:"status"`
}

type Task struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      status    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
}

func NewCreateTaskRequest(title, desc string, status status) *CreateTaskRequest {
	return &CreateTaskRequest{
		Title:       title,
		Description: desc,
		Status:      status,
	}
}

func NewTask(title, desc string, status status, createdAt time.Time) *Task {

	return &Task{
		Title:       title,
		Description: desc,
		Status:      status,
		CreatedAt:   createdAt,
	}
}
