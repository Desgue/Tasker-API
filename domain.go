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

// This type is used to define the priority of a project as a Iota
const (
	High   Priority = "High"
	Medium Priority = "Medium"
	Low    Priority = "Low"
)

type Priority string

// This struct hold the project's tasks received from the database

type Project struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    Priority  `json:"priority"`
	CreatedAt   time.Time `json:"createdAt"`
}

// This struct is used for holding the request data for creating a new project
type CreateProjectRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Priority    Priority `json:"priority"`
}

func NewCreateProjectRequest(title, desc string, priority Priority) *CreateProjectRequest {
	return &CreateProjectRequest{
		Title:       title,
		Description: desc,
		Priority:    priority,
	}
}
