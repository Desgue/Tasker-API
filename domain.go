package main

import "github.com/google/uuid"

const (
	Pending status = iota
	InProgress
	Done
)

type status int

type Project struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Tasks       []Task `json:" tasks"`
	Status      status `json:" status"`
}

type Task struct {
	id     string
	title  string
	status status
}

func NewProject(title, desc string, status status) *Project {

	return &Project{
		Id:          uuid.NewString(),
		Title:       title,
		Description: desc,
		Tasks:       []Task{},
		Status:      status,
	}
}
func NewTask(title string, status status) *Task {
	return &Task{
		id:     uuid.NewString(),
		title:  title,
		status: status,
	}
}
