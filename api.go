package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	listenAddr string
	service    ITaskService
}

func NewApiServer(addr string, svc ITaskService) *ApiServer {
	return &ApiServer{
		listenAddr: addr,
		service:    svc,
	}
}

func (s *ApiServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/tasks", makeHttpHandler(s.handleTasks))
	router.HandleFunc("/tasks/{id}", makeHttpHandler(s.handleTask))
	log.Println("Server running and listening on port: ", s.listenAddr)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	})
	handler := c.Handler(router)
	http.ListenAndServe(s.listenAddr, handler)

}

func (s *ApiServer) handleTasks(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetTasks(w, r)
	case "POST":
		return s.handleCreateTask(w, r)
	default:
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: "Method not allowed"})
	}

}

func (s *ApiServer) handleTask(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetTaskById(w, r)
	case "PUT":
		return s.handleUpdateTask(w, r)
	case "DELETE":
		return s.handleDeleteTask(w, r)
	default:
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: "Method not allowed"})
	}

}

func (s *ApiServer) handleGetTasks(w http.ResponseWriter, r *http.Request) error {
	log.Println("GET request ")
	tasks, err := s.service.GetTasks()
	if err != nil {
		log.Println(err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}
	return WriteJson(w, http.StatusOK, tasks)
}

func (s *ApiServer) handleGetTaskById(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	log.Printf("GET request at http://localhost:3000/tasks/%s", id)
	task, err := s.service.GetTaskById(id)
	if err != nil {
		log.Println(err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}
	return WriteJson(w, http.StatusOK, &task)
}

func (s *ApiServer) handleCreateTask(w http.ResponseWriter, r *http.Request) error {
	log.Println("POST resquest at http://localhost:3000/tasks")

	createTaskReq := new(CreateTaskRequest)
	if err := json.NewDecoder(r.Body).Decode(createTaskReq); err != nil {
		log.Panicln(err)
		return err
	}
	_, err := s.service.CreateTask(createTaskReq)
	if err != nil {
		log.Println("Error from databade while creating task: ", err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}

	return WriteJson(w, http.StatusOK, ApiLog{StatusCode: http.StatusOK, Msg: "Task created successfully"})
}

func (s *ApiServer) handleUpdateTask(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	log.Printf("PUT request at http://localhost:3000/tasks/%s", id)
	task := new(CreateTaskRequest)
	if err := json.NewDecoder(r.Body).Decode(task); err != nil {
		log.Panicln(err)
	}
	if err := s.service.UpdateTask(id, task); err != nil {
		log.Println(err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}
	return WriteJson(w, http.StatusOK, ApiLog{StatusCode: http.StatusOK, Msg: fmt.Sprintf("Task with id %s updated successfully", id)})
}
func (s *ApiServer) handleDeleteTask(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	log.Printf("DELETE request at http://localhost:3000/tasks/%s", id)

	err := s.service.DeleteTask(id)
	if err != nil {
		log.Println(err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}

	return WriteJson(w, http.StatusOK, ApiLog{StatusCode: http.StatusOK, Msg: fmt.Sprintf("Task with id %s deleted successfully", id)})
}
