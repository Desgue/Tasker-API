package main

import "log"

type IProjectService interface {
	GetProjects() ([]Project, error)
	CreateProject(*CreateProjectRequest) (*CreateProjectRequest, error)
	GetProjectById(string) (Project, error)
}

type ProjectService struct {
	store Storage
}

func NewProjectService(store Storage) *ProjectService {
	return &ProjectService{
		store: store,
	}
}

func (s *ProjectService) GetProjects() ([]Project, error) {
	projects, err := s.store.GetProjects()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return projects, nil
}

func (s *ProjectService) GetProjectById(id string) (Project, error) {
	project, err := s.store.GetProjectById(id)
	if err != nil {
		log.Println(err)
		return Project{}, err
	}
	return project, nil
}

func (s *ProjectService) CreateProject(r *CreateProjectRequest) (*CreateProjectRequest, error) {
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

	if err := s.store.CreateProject(r); err != nil {
		return &CreateProjectRequest{}, err
	}
	return r, nil
}
