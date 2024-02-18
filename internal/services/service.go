package svc

import (
	"log"

	"github.com/Desgue/ttracker-api/internal/domain"
	repo "github.com/Desgue/ttracker-api/internal/repository"
)

type ITaskService interface {
	GetTasks(projectId int) ([]domain.Task, error)
	CreateTask(*domain.CreateTaskRequest) (*domain.CreateTaskRequest, error)
	GetTaskById(string) (domain.Task, error)
	UpdateTask(string, *domain.CreateTaskRequest) error
	DeleteTask(string) error
}

type TaskService struct {
	store repo.TaskStorage
}

func NewTaskService(store repo.TaskStorage) *TaskService {
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

// domain.Project service that handles business logic before inserting project into the database

type IProjectService interface {
	GetProjects(userId string) ([]domain.Project, error)
	CreateProject(*domain.CreateProjectRequest) error
	GetProjectById(projectId, cognitoId string) (domain.Project, error)
	UpdateProject(string, *domain.CreateProjectRequest) error
	DeleteProject(projectId, cognitoId string) error
}

type ProjectService struct {
	store repo.ProjectStorage
}

func NewProjectService(store repo.ProjectStorage) *ProjectService {
	return &ProjectService{
		store: store,
	}
}

func (s *ProjectService) GetProjects(cognitoId string) ([]domain.Project, error) {
	projects, err := s.store.GetProjects(cognitoId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return projects, nil
}

func (s *ProjectService) GetProjectById(projectId, cognitoId string) (domain.Project, error) {
	project, err := s.store.GetProjectById(projectId, cognitoId)
	if err != nil {
		log.Println(err)
		return domain.Project{}, err
	}
	return project, nil
}

func (s *ProjectService) CreateProject(r *domain.CreateProjectRequest) error {
	switch r.Priority {
	case "domain.High", "high", "HIGH":
		r.Priority = domain.High
	case "domain.Medium", "medium", "MEDIUM":
		r.Priority = domain.Medium
	case "domain.Low", "low", "LOW":
		r.Priority = domain.Low
	default:
		r.Priority = domain.Low
	}

	if err := s.store.CreateProject(r); err != nil {

		return err
	}
	return nil
}

func (s *ProjectService) UpdateProject(id string, r *domain.CreateProjectRequest) error {
	switch r.Priority {
	case "domain.High", "high", "HIGH":
		r.Priority = domain.High
	case "domain.Medium", "medium", "MEDIUM":
		r.Priority = domain.Medium
	case "domain.Low", "low", "LOW":
		r.Priority = domain.Low
	default:
		r.Priority = domain.Low
	}

	if err := s.store.UpdateProject(id, r); err != nil {

		return err
	}
	return nil
}

func (s *ProjectService) DeleteProject(projectId, cognitoId string) error {
	if err := s.store.DeleteProject(projectId, cognitoId); err != nil {
		return err
	}
	return nil
}

// User service that handles business logic before inserting user into the database
type IUserService interface {
	CreateUser(string) error
}

type UserService struct {
	store repo.UserStorage
}

func NewUserService(store repo.UserStorage) *UserService {
	return &UserService{
		store: store,
	}
}

func (s *UserService) CreateUser(cognitoId string) error {
	if err := s.store.CreateUser(cognitoId); err != nil {
		return err
	}
	return nil
}
