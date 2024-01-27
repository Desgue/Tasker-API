package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=ttracker sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{
		db: db,
	}, nil
}

func (db *PostgresStore) GetProjects() ([]Project, error) {
	return []Project{}, nil
}

func (db *PostgresStore) GetProjectById(id string) (Project, error) {
	return Project{}, nil
}

func (db *PostgresStore) CreateProject(p Project) error {
	return nil
}
