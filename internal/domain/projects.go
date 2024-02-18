package domain

import "time"

// This type is used to define the priority of a project as a Iota
const (
	High   Priority = "High"
	Medium Priority = "Medium"
	Low    Priority = "Low"
)

type Priority string

// This is the interface that the service will use to interact with the database
type ProjectStorage interface {
	GetProjects(userId string) ([]Project, error)
	GetProjectById(projectId, cognitoId string) (Project, error)
	CreateProject(*CreateProjectRequest) error
	UpdateProject(string, *CreateProjectRequest) error
	DeleteProject(projectId, cognitoId string) error
}

type IProjectService interface {
	GetProjects(userId string) ([]Project, error)
	CreateProject(*CreateProjectRequest) error
	GetProjectById(projectId, cognitoId string) (Project, error)
	UpdateProject(string, *CreateProjectRequest) error
	DeleteProject(projectId, cognitoId string) error
}

// This struct hold the project's tasks received from the database

type Project struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    Priority  `json:"priority"`
	CreatedAt   time.Time `json:"createdAt"`
	UserId      string    `json:"userId"`
}

// This struct is used for holding the request data for creating a new project
type CreateProjectRequest struct {
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Priority      Priority `json:"priority"`
	UserCognitoId string   `json:"userCognitoId"`
}

func NewCreateProjectRequest(title, desc string, priority Priority) *CreateProjectRequest {
	return &CreateProjectRequest{
		Title:       title,
		Description: desc,
		Priority:    priority,
	}
}
