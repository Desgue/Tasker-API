package domain

import "errors"

var (
	ErrInvalidTeamName  = errors.New("invalid team name")
	ErrInvalidTeamAdmin = errors.New("invalid team admin")
)

type Team struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Admin       User   `json:"admin"`
}

type CreateTeamRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Admin       User   `json:"admin"`
}

func (t *Team) Validate() error {
	if t.Name == "" {
		return ErrInvalidTeamName
	}
	if t.Admin.Id == 0 {
		return ErrInvalidTeamAdmin
	}

	return nil
}
