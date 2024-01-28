package main

import "log"

type ITaskService interface {
	GetTasks() ([]Task, error)
	CreateTask(*CreateTaskRequest) (*CreateTaskRequest, error)
	GetTaskById(string) (Task, error)
	UpdateTask(string, *CreateTaskRequest) error
	DeleteTask(string) error
}

type TaskService struct {
	store Storage
}

func NewTaskService(store Storage) *TaskService {
	return &TaskService{
		store: store,
	}
}

func (s *TaskService) GetTasks() ([]Task, error) {
	projects, err := s.store.GetTasks()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return projects, nil
}

func (s *TaskService) GetTaskById(id string) (Task, error) {
	project, err := s.store.GetTaskById(id)
	if err != nil {
		log.Println(err)
		return Task{}, err
	}
	return project, nil
}

func (s *TaskService) CreateTask(r *CreateTaskRequest) (*CreateTaskRequest, error) {
	switch r.Status {
	case "Pending":
		r.Status = Pending
	case "InProgress":
		r.Status = InProgress
	case "Done":
		r.Status = Done
	default:
		r.Status = Pending
	}

	if err := s.store.CreateTask(r); err != nil {
		return &CreateTaskRequest{}, err
	}
	return r, nil
}

func (s *TaskService) UpdateTask(id string, r *CreateTaskRequest) error {
	if err := s.store.UpdateTask(id, r); err != nil {
		return err
	}
	return nil
}

func (s *TaskService) DeleteTask(id string) error {
	if err := s.store.DeleteTask(id); err != nil {
		return err
	}
	return nil
}
