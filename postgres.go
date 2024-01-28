package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	createTypeEnumQuery     = `CREATE OR REPLACE TYPE  status as ENUM('Pending', 'InProgress', 'Done')`
	createProjectTableQuery = `
	CREATE TABLE IF NOT EXISTS Projects (
	id SMALLINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	title varchar(255),
	description text,
	status status,
	createdAt TIMESTAMPTZ NOT NULL DEFAULT NOW()
);`
)

type Storage interface {
	GetProjects() ([]Project, error)
	GetProjectById(string) (Project, error)
	CreateProject(*CreateProjectRequest) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=ttracker dbname=ttracker password=ttracker sslmode=disable"
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

func (store *PostgresStore) Init() error {
	return store.createProjectTable()
}

func (store *PostgresStore) createProjectTable() error {
	_, err := store.db.Exec(createProjectTableQuery)
	if err != nil {
		return err
	}
	return nil
}

func (store *PostgresStore) GetProjects() ([]Project, error) {
	rows, err := store.db.Query("SELECT * from Projects")
	if err != nil {
		return nil, err
	}
	var projects []Project
	for rows.Next() {
		project := Project{}
		err = rows.Scan(&project.Id, &project.Title, &project.Description, &project.Status, &project.CreatedAt)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)

	}
	fmt.Println(projects)
	return projects, nil
}

func (store *PostgresStore) GetProjectById(id string) (Project, error) {
	rows, err := store.db.Query("SELECT * from Projects WHERE id=$1", id)
	log.Println("Fetching project with id: ", id)
	for rows.Next() {
		project := Project{}
		err = rows.Scan(&project.Id, &project.Title, &project.Description, &project.Status, &project.CreatedAt)
		if err != nil {
			return Project{}, err
		}
	}
	return Project{}, nil
}

func (store *PostgresStore) CreateProject(p *CreateProjectRequest) error {
	_, err := store.db.Exec("INSERT INTO Projects (title, description, status) VALUES($1, $2, $3)", p.Title, p.Description, p.Status)
	if err != nil {
		return err
	}

	return nil
}
