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

	createTaskReq := new(domain.CreateTaskRequest)
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
	task := new(domain.CreateTaskRequest)
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
