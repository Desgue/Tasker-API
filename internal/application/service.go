package main

import "log"

type ITaskService interface {
	GetTasks(projectId int) ([]Task, error)
	CreateTask(*CreateTaskRequest) (*CreateTaskRequest, error)
	GetTaskById(string) (Task, error)
	UpdateTask(string, *CreateTaskRequest) error
	DeleteTask(string) error
}

type TaskService struct {
	store TaskStorage
}

func NewTaskService(store TaskStorage) *TaskService {
	return &TaskService{
		store: store,
	}
}

func (s *TaskService) GetTasks(projectId int) ([]Task, error) {
	projects, err := s.store.GetTasks(projectId)
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
	case "Pending", "pending", "PENDING":
		r.Status = Pending
	case "InProgress", "inprogress", "Inprogress", "INPROGRESS", "In Progress", "in progress", "IN PROGRESS":
		r.Status = InProgress
	case "Done", "done", "DONE":
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
	switch r.Status {
	case "Pending", "pending", "PENDING":
		r.Status = Pending
	case "InProgress", "inprogress", "Inprogress", "INPROGRESS", "In Progress", "in progress", "IN PROGRESS":
		r.Status = InProgress
	case "Done", "done", "DONE":
		r.Status = Done
	default:
		r.Status = Pending
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

// Project service that handles business logic before inserting project into the database

type IProjectService interface {
	GetProjects(userId string) ([]Project, error)
	CreateProject(*CreateProjectRequest) error
	GetProjectById(projectId, cognitoId string) (Project, error)
	UpdateProject(string, *CreateProjectRequest) error
	DeleteProject(projectId, cognitoId string) error
}

type ProjectService struct {
	store ProjectStorage
}

func NewProjectService(store ProjectStorage) *ProjectService {
	return &ProjectService{
		store: store,
	}
}

func (s *ProjectService) GetProjects(cognitoId string) ([]Project, error) {
	projects, err := s.store.GetProjects(cognitoId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return projects, nil
}

func (s *ProjectService) GetProjectById(projectId, cognitoId string) (Project, error) {
	project, err := s.store.GetProjectById(projectId, cognitoId)
	if err != nil {
		log.Println(err)
		return Project{}, err
	}
	return project, nil
}

func (s *ProjectService) CreateProject(r *CreateProjectRequest) error {
	switch r.Priority {
	case "High", "high", "HIGH":
		r.Priority = High
	case "Medium", "medium", "MEDIUM":
		r.Priority = Medium
	case "Low", "low", "LOW":
		r.Priority = Low
	default:
		r.Priority = Low
	}

	if err := s.store.CreateProject(r); err != nil {

		return err
	}
	return nil
}

func (s *ProjectService) UpdateProject(id string, r *CreateProjectRequest) error {
	switch r.Priority {
	case "High", "high", "HIGH":
		r.Priority = High
	case "Medium", "medium", "MEDIUM":
		r.Priority = Medium
	case "Low", "low", "LOW":
		r.Priority = Low
	default:
		r.Priority = Low
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
	store UserStorage
}

func NewUserService(store UserStorage) *UserService {
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
