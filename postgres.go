package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	createTypeEnumQuery  = `CREATE TYPE status as ENUM('Pending', 'InProgress', 'Done')`
	createTaskTableQuery = `
	CREATE TYPE status as ENUM('Pending', 'InProgress', 'Done');
	CREATE TABLE IF NOT EXISTS Tasks (
	id SMALLINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	title varchar(255),
	description text,
	status status DEFAULT 'Pending',
	createdAt TIMESTAMPTZ NOT NULL DEFAULT NOW()
);`
)

type Storage interface {
	GetTasks() ([]Task, error)
	GetTaskById(string) (Task, error)
	CreateTask(*CreateTaskRequest) error
	UpdateTask(string, *CreateTaskRequest) error
	DeleteTask(string) error
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
	return store.createTaskTable()
}

func (store *PostgresStore) createTaskTable() error {
	_, err := store.db.Exec(createTaskTableQuery)
	if err != nil {
		return err
	}
	return nil
}

func (store *PostgresStore) GetTasks() ([]Task, error) {
	rows, err := store.db.Query("SELECT * from Tasks")
	if err != nil {
		return nil, err
	}
	var tasks []Task
	for rows.Next() {
		task := Task{}
		err = rows.Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.CreatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)

	}
	fmt.Println(tasks)
	return tasks, nil
}

func (store *PostgresStore) GetTaskById(id string) (Task, error) {
	rows, err := store.db.Query("SELECT * from Tasks WHERE id=$1", id)

	var task Task
	for rows.Next() {
		task = Task{}
		err = rows.Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.CreatedAt)
		if err != nil {
			return Task{}, err
		}
	}
	return task, nil
}

func (store *PostgresStore) CreateTask(p *CreateTaskRequest) error {
	_, err := store.db.Exec("INSERT INTO Tasks (title, description, status) VALUES($1, $2, $3)", p.Title, p.Description, p.Status)
	if err != nil {
		return err
	}

	return nil
}

func (store *PostgresStore) UpdateTask(id string, p *CreateTaskRequest) error {
	_, err := store.db.Exec("UPDATE Tasks SET title=$1, description=$2, status=$3 WHERE id=$4", p.Title, p.Description, p.Status, id)
	if err != nil {
		return err
	}
	return nil
}

func (store *PostgresStore) DeleteTask(id string) error {
	_, err := store.db.Exec("DELETE FROM Tasks WHERE id=$1", id)
	if err != nil {
		return err
	}
	return nil
}
