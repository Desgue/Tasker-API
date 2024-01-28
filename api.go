package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiLog struct {
	Err        string `json:"err"`
	StatusCode int    `json:"statusCode"`
	Msg        string `json:"msg"`
}

func makeHttpHandler(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
		}
	}
}

type ApiServer struct {
	listenAddr string
	service    IProjectService
}

func NewApiServer(addr string, svc IProjectService) *ApiServer {
	return &ApiServer{
		listenAddr: addr,
		service:    svc,
	}
}

func (s *ApiServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/projects", makeHttpHandler(s.handleProjects))
	router.HandleFunc("/projects/{id}", makeHttpHandler(s.handleProject))
	log.Println("Server running and listening on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)

}

func (s *ApiServer) handleProjects(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetProjects(w, r)
	case "POST":
		return s.handleCreateProject(w, r)
	default:
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: "Method not allowed"})
	}

}

func (s *ApiServer) handleProject(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetProjectById(w, r)
	case "PUT":
		return s.handleUpdateProject(w, r)
	case "DELETE":
		return s.handleDeleteProject(w, r)
	default:
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: "Method not allowed"})
	}

}

func (s *ApiServer) handleGetProjects(w http.ResponseWriter, r *http.Request) error {
	log.Println("GET request ")
	projects, err := s.service.GetProjects()
	if err != nil {
		log.Println(err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}
	return WriteJson(w, http.StatusOK, projects)
}

func (s *ApiServer) handleGetProjectById(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	log.Printf("GET request at http://localhost:3000/projects/%s", id)
	project, err := s.service.GetProjectById(id)
	if err != nil {
		log.Println(err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}
	return WriteJson(w, http.StatusOK, &project)
}

func (s *ApiServer) handleCreateProject(w http.ResponseWriter, r *http.Request) error {
	log.Println("POST resquest at http://localhost:3000/projects")

	createProjectReq := new(CreateProjectRequest)
	if err := json.NewDecoder(r.Body).Decode(createProjectReq); err != nil {
		log.Panicln(err)
		return err
	}
	_, err := s.service.CreateProject(createProjectReq)
	if err != nil {
		log.Println("Error from databade while creating project: ", err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}

	return WriteJson(w, http.StatusOK, ApiLog{StatusCode: http.StatusOK, Msg: "Project created successfully"})
}

func (s *ApiServer) handleUpdateProject(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	log.Printf("PUT request at http://localhost:3000/projects/%s", id)
	project := new(CreateProjectRequest)
	if err := json.NewDecoder(r.Body).Decode(project); err != nil {
		log.Panicln(err)
	}
	if err := s.service.UpdateProject(id, project); err != nil {
		log.Println(err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}
	return WriteJson(w, http.StatusOK, ApiLog{StatusCode: http.StatusOK, Msg: fmt.Sprintf("Project with id %s updated successfully", id)})
}
func (s *ApiServer) handleDeleteProject(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	log.Printf("DELETE request at http://localhost:3000/projects/%s", id)

	err := s.service.DeleteProject(id)
	if err != nil {
		log.Println(err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}

	return WriteJson(w, http.StatusOK, ApiLog{StatusCode: http.StatusOK, Msg: fmt.Sprintf("Project with id %s deleted successfully", id)})
}
