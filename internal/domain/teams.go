package domain

import "errors"

var (
	ErrInvalidTeamName  = errors.New("invalid team name")
	ErrInvalidTeamAdmin = errors.New("invalid team admin")
)

type TeamStorage interface {
	GetTeam(id int) (Team, error)
	CreateTeam(*CreateTeamRequest) error
}

type ITeamService interface {
	GetTeam(id int) (Team, error)
	CreateTeam(*CreateTeamRequest) error
}

type Team struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	AdminId     int    `json:"adminId"`
}

type CreateTeamRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	AdminId     int    `json:"adminId"`
}

func (t *CreateTeamRequest) Validate() error {
	if t.Name == "" {
		return ErrInvalidTeamName
	}
	if t.AdminId == 0 {
		return ErrInvalidTeamAdmin
	}
	return nil
}
