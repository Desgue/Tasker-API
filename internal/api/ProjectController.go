package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Desgue/ttracker-api/internal/domain"
	"github.com/gorilla/mux"
)

// Handler for calls to /projects

func (s *ApiServer) handleProjects(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetProjects(w, r)
	case "POST":
		return s.handleCreateProject(w, r)
	default:
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: "Method not allowed on /projects"})
	}

}
func (s *ApiServer) handleGetProjects(w http.ResponseWriter, r *http.Request) error {

	cognitoId := r.Header.Get("CognitoId")
	projects, err := s.projectService.GetProjects(cognitoId)
	if err != nil {
		log.Println("Err fetching projects: ", err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}
	return WriteJson(w, http.StatusOK, projects)
}
func (s *ApiServer) handleCreateProject(w http.ResponseWriter, r *http.Request) error {

	createProjectReq := new(domain.CreateProjectRequest)
	if err := json.NewDecoder(r.Body).Decode(createProjectReq); err != nil {
		log.Panicln("Error decoding request body, terminating program: ", err)
		return err
	}
	createProjectReq.UserCognitoId = r.Header.Get("CognitoId")
	if err := s.projectService.CreateProject(createProjectReq); err != nil {
		log.Println("Error creating project: ", err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}
	return WriteJson(w, http.StatusOK, ApiLog{StatusCode: http.StatusOK, Msg: "Project created successfully"})
}

// Handler for calls to /projects/{projectId}

func (s *ApiServer) handleProject(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetProjectById(w, r)
	case
		"PUT":
		return s.handleUpdateProject(w, r)
	case "DELETE":
		return s.handleDeleteProject(w, r)
	default:
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: "Method not allowed on /projects/{id}"})
	}
}

func (s *ApiServer) handleGetProjectById(w http.ResponseWriter, r *http.Request) error {
	projectId := mux.Vars(r)["projectId"]
	cognitoId := r.Header.Get("CognitoId")

	project, err := s.projectService.GetProjectById(projectId, cognitoId)
	if err != nil {
		log.Println("Err fetching project: ", err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}
	return WriteJson(w, http.StatusOK, &project)
}

func (s *ApiServer) handleUpdateProject(w http.ResponseWriter, r *http.Request) error {
	projectId := mux.Vars(r)["projectId"]
	cognitoId := r.Header.Get("CognitoId")

	project := new(domain.CreateProjectRequest)
	if err := json.NewDecoder(r.Body).Decode(project); err != nil {
		log.Panicln("Error decoding request body ", err)
	}
	project.UserCognitoId = cognitoId

	if err := s.projectService.UpdateProject(projectId, project); err != nil {
		log.Println("Err updating project: ", err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}

	return WriteJson(w, http.StatusOK, ApiLog{StatusCode: http.StatusOK, Msg: fmt.Sprintf("Project with id %s updated successfully", projectId)})
}

func (s *ApiServer) handleDeleteProject(w http.ResponseWriter, r *http.Request) error {
	projectId := mux.Vars(r)["projectId"]
	cognitoId := r.Header.Get("CognitoId")

	err := s.projectService.DeleteProject(projectId, cognitoId)
	if err != nil {
		log.Println("Err deleting project: ", err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}

	return WriteJson(w, http.StatusOK, ApiLog{StatusCode: http.StatusOK, Msg: fmt.Sprintf("Project with id %s deleted successfully", projectId)})
}
