package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
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
	listenAddr     string
	taskService    ITaskService
	projectService IProjectService
}

func NewApiServer(addr string, svc ITaskService, psvc IProjectService) *ApiServer {
	return &ApiServer{
		listenAddr:     addr,
		taskService:    svc,
		projectService: psvc,
	}
}

func (s *ApiServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/projects/{projectId}/tasks", makeHttpHandler(s.handleTasks))
	router.HandleFunc("/projects/{projectId}/tasks/{taskId}", makeHttpHandler(s.handleTask))

	router.HandleFunc("/projects", makeHttpHandler(s.handleProjects))
	router.HandleFunc("/projects/{projectId}", makeHttpHandler(s.handleProject))

	router.Use(verifyJwtMiddleware)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	})
	handler := c.Handler(router)

	log.Println("Server running and listening on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, handler)

}

// Handler for calls to /projects/{projectId}/tasks

func (s *ApiServer) handleTasks(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetTasks(w, r)
	case "POST":
		return s.handleCreateTask(w, r)
	default:
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: "Method not allowed on /projects/{projectId}/tasks"})
	}

}

func (s *ApiServer) handleGetTasks(w http.ResponseWriter, r *http.Request) error {
	log.Println("GET request at http://localhost:8000/projects/{projectId}/tasks")

	projectId, err := strconv.Atoi(mux.Vars(r)["projectId"])
	if err != nil {
		log.Println(err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}

	tasks, err := s.taskService.GetTasks(projectId)
	if err != nil {
		log.Println(err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}

	return WriteJson(w, http.StatusOK, tasks)
}

func (s *ApiServer) handleCreateTask(w http.ResponseWriter, r *http.Request) error {
	log.Println("POST resquest at http://localhost:8000/projects/{projectId}/tasks")

	createTaskReq := new(CreateTaskRequest)
	if err := json.NewDecoder(r.Body).Decode(createTaskReq); err != nil {
		log.Panicln(err)
		return err
	}
	_, err := s.taskService.CreateTask(createTaskReq)
	if err != nil {
		log.Println("Error from database while creating task: ", err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}

	return WriteJson(w, http.StatusOK, ApiLog{StatusCode: http.StatusOK, Msg: "Task created successfully"})
}

// Handler for calls to /projects/{projectId}/tasks/{taskId}
func (s *ApiServer) handleTask(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetTaskById(w, r)
	case "PUT":
		return s.handleUpdateTask(w, r)
	case "DELETE":
		return s.handleDeleteTask(w, r)
	default:
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: "Method not allowed on /tasks/{id}"})
	}

}

func (s *ApiServer) handleGetTaskById(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["taskId"]
	log.Printf("GET http://localhost:8000/projects/{projectId}/tasks/%s", id)
	task, err := s.taskService.GetTaskById(id)
	if err != nil {
		log.Println(err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}
	return WriteJson(w, http.StatusOK, &task)
}

func (s *ApiServer) handleUpdateTask(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["taskId"]
	log.Printf("PUT http://localhost:8000/projects/{projectId}/tasks/%s", id)
	task := new(CreateTaskRequest)
	if err := json.NewDecoder(r.Body).Decode(task); err != nil {
		log.Panicln(err)
	}
	if err := s.taskService.UpdateTask(id, task); err != nil {
		log.Println(err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}
	return WriteJson(w, http.StatusOK, ApiLog{StatusCode: http.StatusOK, Msg: fmt.Sprintf("Task with id %s updated successfully", id)})
}
func (s *ApiServer) handleDeleteTask(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["taskId"]
	log.Printf("DELETE request at http://localhost:8000/projects/{projectId}/tasks/%s", id)

	err := s.taskService.DeleteTask(id)
	if err != nil {
		log.Println(err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}

	return WriteJson(w, http.StatusOK, ApiLog{StatusCode: http.StatusOK, Msg: fmt.Sprintf("Task with id %s deleted successfully", id)})
}

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
	log.Println("GET request on /projects")
	userId := mux.Vars(r)["userId"]
	projects, err := s.projectService.GetProjects(userId)
	if err != nil {
		log.Println("Err fetching projects: ", err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}
	return WriteJson(w, http.StatusOK, projects)
}
func (s *ApiServer) handleCreateProject(w http.ResponseWriter, r *http.Request) error {
	log.Println("POST request on /projects")
	createProjectReq := new(CreateProjectRequest)
	if err := json.NewDecoder(r.Body).Decode(createProjectReq); err != nil {
		log.Panicln("Error decoding request body, terminating program: ", err)
		return err
	}
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
	id := mux.Vars(r)["projectId"]
	log.Printf("GET request at http://localhost:3000/projects/%s", id)
	project, err := s.projectService.GetProjectById(id)
	if err != nil {
		log.Println("Err fetching project: ", err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}
	return WriteJson(w, http.StatusOK, &project)
}

func (s *ApiServer) handleUpdateProject(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["projectId"]
	log.Printf("PUT request at http://localhost:3000/projects/%s", id)
	project := new(CreateProjectRequest)
	if err := json.NewDecoder(r.Body).Decode(project); err != nil {
		log.Panicln("Error decoding request body ", err)
	}
	if err := s.projectService.UpdateProject(id, project); err != nil {
		log.Println("Err updating project: ", err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}
	return WriteJson(w, http.StatusOK, ApiLog{StatusCode: http.StatusOK, Msg: fmt.Sprintf("Project with id %s updated successfully", id)})
}

func (s *ApiServer) handleDeleteProject(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["projectId"]
	log.Printf("DELETE request at http://localhost:3000/projects/%s", id)
	err := s.projectService.DeleteProject(id)
	if err != nil {
		log.Println("Err deleting project: ", err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}
	return WriteJson(w, http.StatusOK, ApiLog{StatusCode: http.StatusOK, Msg: fmt.Sprintf("Project with id %s deleted successfully", id)})
}
