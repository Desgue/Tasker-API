package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

type Server struct {
	addr       string
	controller *Controllers
}

func NewServer(addr string, controllers *Controllers) *Server {
	return &Server{
		addr:       addr,
		controller: controllers,
	}
}

type Controllers struct {
	Project *ProjectController
	Task    *TaskController
	Team    *TeamController
	User    *UserController
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
func (s *Server) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/projects/{projectId}/tasks", makeHttpHandler(s.controller.Task.handleTasks))
	router.HandleFunc("/projects/{projectId}/tasks/{taskId}", makeHttpHandler(s.controller.Task.handleTask))

	router.HandleFunc("/projects", makeHttpHandler(s.controller.Project.handleProjects))
	router.HandleFunc("/projects/{projectId}", makeHttpHandler(s.controller.Project.handleProject))

	router.HandleFunc("/teams", makeHttpHandler(s.controller.Team.handleTeams))
	router.HandleFunc("/teams/{teamId}", makeHttpHandler(s.controller.Team.handleTeam))

	router.HandleFunc("/users", makeHttpHandler(s.controller.User.handleUsers))

	router.Use(loggingMiddleware)
	router.Use(verifyJwtMiddleware)
	router.Use(s.controller.User.verifyUserMiddleware)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
	})
	handler := c.Handler(router)

	log.Println("Server running and listening on port: ", s.addr)
	http.ListenAndServe(s.addr, handler)

}
