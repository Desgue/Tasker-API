package main

type Storage interface {
	GetProjects() ([]Project, error)
	GetProjectById(string) (Project, error)
	CreateProject(Project) error
}
