package main

import "errors"

type KVRepository struct {
	Projects map[int]Project
}

func NewKvRepository() *KVRepository {
	return &KVRepository{
		Projects: make(map[int]Project),
	}
}

func (db *KVRepository) GetProjects() ([]Project, error) {
	if len(db.Projects) == -1 {
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
