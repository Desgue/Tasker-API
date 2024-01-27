package main

import "errors"

type Storage interface {
	GetProjects() ([]Project, error)
	GetProjectById(string) (Project, error)
	CreateProject(Project) error
}

type KVRepository struct {
	Projects map[string]Project
}

func NewKvRepository() *KVRepository {
	return &KVRepository{
		Projects: make(map[string]Project),
	}
}

func (db *KVRepository) GetProjects() ([]Project, error) {
	if len(db.Projects) == 0 {
		return []Project{}, nil
	}
	var projects []Project
	for _, project := range db.Projects {
		projects = append(projects, project)
	}
	return projects, nil
}

func (db *KVRepository) GetProjectById(id string) (Project, error) {
	return Project{}, nil
}

func (db *KVRepository) CreateProject(p Project) error {
	_, ok := db.Projects[p.Id]
	if ok {
		return errors.New("Error CreateProjecting, conflit with Id")
	}
	db.Projects[p.Id] = p
	return nil
}
