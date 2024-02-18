package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Desgue/ttracker-api/internal/domain"
	"github.com/gorilla/mux"
)

type TaskController struct {
	service domain.ITaskService
}

func NewTaskController(service domain.ITaskService) *TaskController {
	return &TaskController{
		service: service,
	}
}

// Handler for calls to /projects/{projectId}/tasks

func (s *TaskController) handleTasks(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetTasks(w, r)
	case "POST":
		return s.handleCreateTask(w, r)
	default:
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: "Method not allowed on /projects/{projectId}/tasks"})
	}

}

func (s *TaskController) handleGetTasks(w http.ResponseWriter, r *http.Request) error {
	log.Println("GET request at http://localhost:8000/projects/{projectId}/tasks")

	projectId, err := strconv.Atoi(mux.Vars(r)["projectId"])
	if err != nil {
		log.Println(err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}

	tasks, err := s.service.GetTasks(projectId)
	if err != nil {
		log.Println(err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}

	return WriteJson(w, http.StatusOK, tasks)
}

func (s *TaskController) handleCreateTask(w http.ResponseWriter, r *http.Request) error {
	log.Println("POST resquest at http://localhost:8000/projects/{projectId}/tasks")

	createTaskReq := new(domain.CreateTaskRequest)
	if err := json.NewDecoder(r.Body).Decode(createTaskReq); err != nil {
		log.Panicln(err)
		return err
	}
	_, err := s.service.CreateTask(createTaskReq)
	if err != nil {
		log.Println("Error from database while creating task: ", err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}

	return WriteJson(w, http.StatusOK, ApiLog{StatusCode: http.StatusOK, Msg: "Task created successfully"})
}

// Handler for calls to /projects/{projectId}/tasks/{taskId}
func (s *TaskController) handleTask(w http.ResponseWriter, r *http.Request) error {
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

func (s *TaskController) handleGetTaskById(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["taskId"]
	log.Printf("GET http://localhost:8000/projects/{projectId}/tasks/%s", id)
	task, err := s.service.GetTaskById(id)
	if err != nil {
		log.Println(err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}
	return WriteJson(w, http.StatusOK, &task)
}

func (s *TaskController) handleUpdateTask(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["taskId"]
	log.Printf("PUT http://localhost:8000/projects/{projectId}/tasks/%s", id)
	task := new(domain.CreateTaskRequest)
	if err := json.NewDecoder(r.Body).Decode(task); err != nil {
		log.Panicln(err)
	}
	if err := s.service.UpdateTask(id, task); err != nil {
		log.Println(err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}
	return WriteJson(w, http.StatusOK, ApiLog{StatusCode: http.StatusOK, Msg: fmt.Sprintf("Task with id %s updated successfully", id)})
}
func (s *TaskController) handleDeleteTask(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["taskId"]
	log.Printf("DELETE request at http://localhost:8000/projects/{projectId}/tasks/%s", id)

	err := s.service.DeleteTask(id)
	if err != nil {
		log.Println(err)
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}

	return WriteJson(w, http.StatusOK, ApiLog{StatusCode: http.StatusOK, Msg: fmt.Sprintf("Task with id %s deleted successfully", id)})
}
