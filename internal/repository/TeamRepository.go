package repo

import (
	"database/sql"

	"github.com/Desgue/ttracker-api/internal/domain"
)

type PostgresTeamStore struct {
	DB *sql.DB
}

func NewPostgresTeamStore(DB *sql.DB) *PostgresTeamStore {
	return &PostgresTeamStore{
		DB: DB,
	}
}

func (store *PostgresTeamStore) GetTeam(id int) (domain.Team, error) {
	var team domain.Team
	err := store.DB.QueryRow("SELECT * from Teams WHERE id=$1", id).Scan(
		&team.Id,
		&team.Name,
		&team.Description,
		&team.AdminId,
	)
	if err != nil {
		return domain.Team{}, err
	}

	return team, nil
}

func (store *PostgresTeamStore) CreateTeam(p *domain.CreateTeamRequest) error {
	_, err := store.DB.Exec(`
	INSERT INTO Teams (name, description, adminId) VALUES($1, $2, $3)`,
		p.Name, p.Description, p.AdminId)
	if err != nil {
		return err
	}
	return nil
}
