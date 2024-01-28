package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	createTypeEnumQuery     = `CREATE TYPE  status as ENUM('Pending', 'InProgress', 'Done')`
	createProjectTableQuery = `
	CREATE TABLE IF NOT EXISTS Projects (
	id varchar(255) NOT NULL PRIMARY KEY,
	title varchar(255),
	description text,
	status status,
	createdAt timestamp

);`
)

type Storage interface {
	GetProjects() ([]Project, error)
	GetProjectById(string) (Project, error)
	CreateProject(*Project) error
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
	return Project{}, nil
}

func (store *PostgresStore) CreateProject(p *Project) error {
	_, err := store.db.Exec("INSERT INTO Projects VALUES($1, $2, $3, $4, $5)", p.Id, p.Title, p.Description, p.Status, p.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}
