package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Desgue/ttracker-api/internal/domain"
)

type TeamController struct {
	service domain.ITeamService
}

func NewTeamController(service domain.ITeamService) *TeamController {
	return &TeamController{
		service: service,
	}
}

// Handler for calls to /teams
func (c *TeamController) handleTeams(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		return c.handleCreateTeam(w, r)
	default:
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: "Method not allowed on /teams"})
	}

}

func (c *TeamController) handleCreateTeam(w http.ResponseWriter, r *http.Request) error {

	team := new(domain.CreateTeamRequest)

	if err := json.NewDecoder(r.Body).Decode(team); err != nil {
		return err
	}

	if err := c.service.CreateTeam(team); err != nil {
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}
	return WriteJson(w, http.StatusOK, ApiLog{StatusCode: http.StatusOK, Msg: fmt.Sprintf("Team %s created successfully", team.Name)})
}

// Handler for calls to /teams/{teamId}
func (c *TeamController) handleTeam(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return c.handleGetTeam(w, r)
	default:
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: "Method not allowed on /teams/{teamId}"})
	}
}

func (c *TeamController) handleGetTeam(w http.ResponseWriter, r *http.Request) error {
	teamId, err := strconv.Atoi(r.URL.Query().Get("teamId"))
	if err != nil {
		return WriteJson(w, http.StatusInternalServerError, ApiLog{Err: err.Error(), StatusCode: http.StatusInternalServerError})
	}

	team, err := c.service.GetTeam(teamId)
	if err != nil {
		return WriteJson(w, http.StatusBadRequest, ApiLog{Err: err.Error(), StatusCode: http.StatusBadRequest})
	}
	return WriteJson(w, http.StatusOK, team)
}
