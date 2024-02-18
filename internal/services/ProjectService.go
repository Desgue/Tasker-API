package svc

import (
	"log"

	"github.com/Desgue/ttracker-api/internal/domain"
)

// domain.Project service that handles business logic before inserting project into the database

type ProjectService struct {
	store domain.ProjectStorage
}

func NewProjectService(store domain.ProjectStorage) *ProjectService {
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
