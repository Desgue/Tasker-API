package main

import "time"

// This type is used to define the priority of a project as a Iota
const (
	High   Priority = "High"
	Medium Priority = "Medium"
	Low    Priority = "Low"
)

type Priority string

// This struct hold the project's tasks received from the database

type Project struct {
	Id            int       `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Priority      Priority  `json:"priority"`
	CreatedAt     time.Time `json:"createdAt"`
	UserCognitoId string    `json:"userCognitoId"`
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
