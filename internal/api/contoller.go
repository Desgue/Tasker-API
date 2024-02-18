package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Desgue/ttracker-api/internal/domain"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiServer struct {
	listenAddr     string
	taskService    domain.ITaskService
	projectService domain.IProjectService
}

func NewApiServer(addr string, svc domain.ITaskService, psvc domain.IProjectService) *ApiServer {
	return &ApiServer{
		listenAddr:     addr,
		taskService:    svc,
		projectService: psvc,
	}
}

type ApiLog struct {
	Err        string `json:"err"`
	StatusCode int    `json:"statusCode"`
	Msg        string `json:"msg"`
}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	return json.NewEncoder(w).Encode(v)
}

func makeHttpHandler(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
		}
	}
}

func (s *ApiServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/projects/{projectId}/tasks", makeHttpHandler(s.handleTasks))
	router.HandleFunc("/projects/{projectId}/tasks/{taskId}", makeHttpHandler(s.handleTask))

	router.HandleFunc("/projects", makeHttpHandler(s.handleProjects))
	router.HandleFunc("/projects/{projectId}", makeHttpHandler(s.handleProject))

	router.Use(loggingMiddleware)
	router.Use(verifyJwtMiddleware)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
	})
	handler := c.Handler(router)

	log.Println("Server running and listening on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, handler)

}
