package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiErr struct {
	Err string `json:"err"`
}

func makeHttpHandler(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			WriteJson(w, http.StatusBadRequest, ApiErr{Err: err.Error()})
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

	log.Println("Server running and listening on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)

}

func (s *ApiServer) handleProjects(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetProjects(w, r)
	case "POST":
		return s.handleCreateProject(w, r)
	case "PUT":
		return s.handleUpdateProject(w, r)
	case "DELETE":
		return s.handleDeleteProject(w, r)
	}

	return nil
}

func (s *ApiServer) handleGetProjects(w http.ResponseWriter, r *http.Request) error {
	log.Println("GET request ")
	projects, err := s.service.GetProjects()
	if err != nil {
		log.Println(err)
		return WriteJson(w, http.StatusBadRequest, ApiErr{Err: err.Error()})
	}
	return WriteJson(w, http.StatusOK, projects)
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
		log.Println(err)
		return WriteJson(w, http.StatusBadRequest, ApiErr{Err: err.Error()})
	}

	return nil
}
func (s *ApiServer) handleUpdateProject(w http.ResponseWriter, r *http.Request) error {
	log.Println("PUT request")
	return nil
}
func (s *ApiServer) handleDeleteProject(w http.ResponseWriter, r *http.Request) error {
	log.Println("DELETE request")
	return nil
}
