package domain

import (
	"time"
)

const (
	Pending    status = "Pending"
	InProgress status = "InProgress"
	Done       status = "Done"
)

type status string

type TaskStorage interface {
	GetTasks(projectId int) ([]Task, error)
	GetTaskById(string) (Task, error)
	CreateTask(*CreateTaskRequest) error
	UpdateTask(string, *CreateTaskRequest) error
	DeleteTask(string) error
}

type ITaskService interface {
	GetTasks(projectId int) ([]Task, error)
	CreateTask(*CreateTaskRequest) (*CreateTaskRequest, error)
	GetTaskById(string) (Task, error)
	UpdateTask(string, *CreateTaskRequest) error
	DeleteTask(string) error
}

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      status `json:"status"`
	ProjectId   int    `json:"projectId"`
}

type Task struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      status    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	ProjectId   int       `json:"projectId"`
}

func NewCreateTaskRequest(title, desc string, status status, projectId int) *CreateTaskRequest {
	return &CreateTaskRequest{
		Title:       title,
		Description: desc,
		Status:      status,
		ProjectId:   projectId,
	}
}

func NewTask(title, desc string, status status, projectId int, createdAt time.Time) *Task {
	return &Task{
		Title:       title,
		Description: desc,
		Status:      status,
		ProjectId:   projectId,
		CreatedAt:   createdAt,
	}
}
