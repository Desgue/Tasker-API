package svc

import (
	"log"

	"github.com/Desgue/ttracker-api/internal/domain"
)

type TaskService struct {
	store domain.TaskStorage
}

func NewTaskService(store domain.TaskStorage) *TaskService {
	return &TaskService{
		store: store,
	}
}

func (s *TaskService) GetTasks(projectId int) ([]domain.Task, error) {
	projects, err := s.store.GetTasks(projectId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return projects, nil
}

func (s *TaskService) GetTaskById(id string) (domain.Task, error) {
	project, err := s.store.GetTaskById(id)
	if err != nil {
		log.Println(err)
		return domain.Task{}, err
	}
	return project, nil
}

func (s *TaskService) CreateTask(r *domain.CreateTaskRequest) (*domain.CreateTaskRequest, error) {
	switch r.Status {
	case "domain.Pending", "pending", "PENDING":
		r.Status = domain.Pending
	case "domain.InProgress", "inprogress", "Inprogress", "INPROGRESS", "In Progress", "in progress", "IN PROGRESS":
		r.Status = domain.InProgress
	case "domain.Done", "done", "DONE":
		r.Status = domain.Done
	default:
		r.Status = domain.Pending
	}

	if err := s.store.CreateTask(r); err != nil {
		return &domain.CreateTaskRequest{}, err
	}
	return r, nil
}

func (s *TaskService) UpdateTask(id string, r *domain.CreateTaskRequest) error {
	switch r.Status {
	case "domain.Pending", "pending", "PENDING":
		r.Status = domain.Pending
	case "domain.InProgress", "inprogress", "Inprogress", "INPROGRESS", "In Progress", "in progress", "IN PROGRESS":
		r.Status = domain.InProgress
	case "domain.Done", "done", "DONE":
		r.Status = domain.Done
	default:
		r.Status = domain.Pending
	}

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
