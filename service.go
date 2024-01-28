package main

import "log"

type IProjectService interface {
	GetProjects() ([]Project, error)
	CreateProject(*CreateProjectRequest) (*Project, error)
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
		log.Panicln(err)
		return nil, err
	}
	return projects, nil
}
func (s *ProjectService) CreateProject(r *CreateProjectRequest) (*Project, error) {
	var status status
	switch r.Status {
	case "Pending":
		status = Pending
	case "InProgress":
		status = InProgress
	case "Done":
		status = Done
	default:
		status = Pending
	}

	project := NewProject(r.Title, r.Description, status)

	if err := s.store.CreateProject(project); err != nil {
		return &Project{}, err
	}
	return project, nil
}
