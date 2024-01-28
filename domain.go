package main

import (
	"time"

	"github.com/google/uuid"
)

const (
	Pending    status = "Pending"
	InProgress status = "InProgress"
	Done       status = "Done"
)

type status string

type CreateProjectRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      status `json:"status"`
}

type Project struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      status    `json:" status"`
	CreatedAt   time.Time `json:"createdAt"`
}

type Task struct {
	id        string
	title     string
	status    status
	CreatedAt time.Time `json:"createdAt"`
}

func NewCreateProjectRequest(title, desc string, status status) *CreateProjectRequest {
	return &CreateProjectRequest{
		Title:       title,
		Description: desc,
		Status:      status,
	}
}

func NewProject(title, desc string, status status, createdAt time.Time) *Project {

	return &Project{
		Title:       title,
		Description: desc,
		Status:      status,
		CreatedAt:   createdAt,
	}
}
func NewTask(title string, status status) *Task {
	return &Task{
		id:     uuid.NewString(),
		title:  title,
		status: status,
	}
}
