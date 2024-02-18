package api

import "github.com/Desgue/ttracker-api/internal/domain"

type TeamController struct {
	service domain.ITeamService
}

func NewTeamController(service domain.ITeamService) *TeamController {
	return &TeamController{
		service: service,
	}
}
