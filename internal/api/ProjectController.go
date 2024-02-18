package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Desgue/ttracker-api/internal/domain"
	"github.com/gorilla/mux"
)

type ProjectController struct {
	service domain.IProjectService
}

func NewProjectController(service domain.IProjectService) *ProjectController {
	return &ProjectController{
		service: service,
	}
}

// Handler for calls to /projects

func (c *ProjectController) handleProjects(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return c.handleGetProjects(w, r)
	case "POST":
		return c.handleCreateProject(w, r)
	default:
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: "Method not allowed on /projects"})
	}

}
func (c *ProjectController) handleGetProjects(w http.ResponseWriter, r *http.Request) error {

	cognitoId := r.Header.Get("CognitoId")
	projects, err := c.service.GetProjects(cognitoId)
	if err != nil {
		log.Println("Err fetching projects: ", err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}
	return WriteJson(w, http.StatusOK, projects)
}
func (c *ProjectController) handleCreateProject(w http.ResponseWriter, r *http.Request) error {

	createProjectReq := new(domain.CreateProjectRequest)
	if err := json.NewDecoder(r.Body).Decode(createProjectReq); err != nil {
		log.Panicln("Error decoding request body, terminating program: ", err)
		return err
	}
	createProjectReq.UserCognitoId = r.Header.Get("CognitoId")
	if err := c.service.CreateProject(createProjectReq); err != nil {
		log.Println("Error creating project: ", err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}
	return WriteJson(w, http.StatusOK, ApiLog{StatusCode: http.StatusOK, Msg: "Project created successfully"})
}

// Handler for calls to /projects/{projectId}

func (c *ProjectController) handleProject(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return c.handleGetProjectById(w, r)
	case
		"PUT":
		return c.handleUpdateProject(w, r)
	case "DELETE":
		return c.handleDeleteProject(w, r)
	default:
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: "Method not allowed on /projects/{id}"})
	}
}

func (c *ProjectController) handleGetProjectById(w http.ResponseWriter, r *http.Request) error {
	projectId := mux.Vars(r)["projectId"]
	cognitoId := r.Header.Get("CognitoId")

	project, err := c.service.GetProjectById(projectId, cognitoId)
	if err != nil {
		log.Println("Err fetching project: ", err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}
	return WriteJson(w, http.StatusOK, &project)
}

func (c *ProjectController) handleUpdateProject(w http.ResponseWriter, r *http.Request) error {
	projectId := mux.Vars(r)["projectId"]
	cognitoId := r.Header.Get("CognitoId")

	project := new(domain.CreateProjectRequest)
	if err := json.NewDecoder(r.Body).Decode(project); err != nil {
		log.Panicln("Error decoding request body ", err)
	}
	project.UserCognitoId = cognitoId

	if err := c.service.UpdateProject(projectId, project); err != nil {
		log.Println("Err updating project: ", err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}

	return WriteJson(w, http.StatusOK, ApiLog{StatusCode: http.StatusOK, Msg: fmt.Sprintf("Project with id %s updated successfully", projectId)})
}

func (c *ProjectController) handleDeleteProject(w http.ResponseWriter, r *http.Request) error {
	projectId := mux.Vars(r)["projectId"]
	cognitoId := r.Header.Get("CognitoId")

	err := c.service.DeleteProject(projectId, cognitoId)
	if err != nil {
		log.Println("Err deleting project: ", err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}

	return WriteJson(w, http.StatusOK, ApiLog{StatusCode: http.StatusOK, Msg: fmt.Sprintf("Project with id %s deleted successfully", projectId)})
}
