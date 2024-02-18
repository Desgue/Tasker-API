package svc

import "github.com/Desgue/ttracker-api/internal/domain"

type TeamService struct {
	store domain.TeamStorage
}

func NewTeamService(store domain.TeamStorage) *TeamService {
	return &TeamService{
		store: store,
	}
}

func (svc *TeamService) GetTeam(id int) (domain.Team, error) {
	team, err := svc.store.GetTeam(id)
	if err != nil {
		return domain.Team{}, err
	}
	return team, nil
}

func (svc *TeamService) CreateTeam(r *domain.CreateTeamRequest) error {
	if err := r.Validate(); err != nil {
		return err
	}
	if err := svc.store.CreateTeam(r); err != nil {
		return err
	}
	return nil
}
