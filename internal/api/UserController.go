package api

import (
	"net/http"

	"github.com/Desgue/ttracker-api/internal/domain"
)

type UserController struct {
	service domain.IUserService
}

func NewUserController(service domain.IUserService) *UserController {
	return &UserController{
		service: service,
	}
}

// Handler for calls to /users

func (c *UserController) handleUsers(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		return c.handleCreateUser(w, r)
	default:
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: "Method not allowed on /users"})
	}

}

func (c *UserController) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("Placeholder for creating a user"))
	return nil
}
