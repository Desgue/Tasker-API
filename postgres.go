package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	createStatusEnumQuery = `CREATE TYPE status as ENUM('Pending', 'InProgress', 'Done');`
	createTaskTableQuery  = `
	CREATE TABLE IF NOT EXISTS Tasks (
	id SMALLINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	title varchar(255),
	description text,
	status status DEFAULT 'Pending',
	createdAt TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	projectId SMALLINT NOT NULL REFERENCES Projects(id) 
);`
	createPriorityEnumQuery = `CREATE TYPE priority as ENUM('High', 'Medium', 'Low');`
	createProjectTableQuery = `
	CREATE TABLE IF NOT EXISTS Projects (
	id SMALLINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	title varchar(255),
	description text,
	priority priority DEFAULT 'Low',
	createdAt TIMESTAMPTZ NOT NULL DEFAULT NOW()
);`
)

type TaskStorage interface {
	GetTasks(projectId int) ([]Task, error)
	GetTaskById(string) (Task, error)
	CreateTask(*CreateTaskRequest) error
	UpdateTask(string, *CreateTaskRequest) error
	DeleteTask(string) error
}

type PostgresTaskStore struct {
	db *sql.DB
}

func NewPostgresTaskStore() (*PostgresTaskStore, error) {
	connStr := "user=ttracker dbname=ttracker password=ttracker sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresTaskStore{
		db: db,
	}, nil
}

func (store *PostgresTaskStore) Init() error {
	return store.createTaskTable()
}

func (store *PostgresTaskStore) createTaskTable() error {
	_, err := store.db.Exec(createStatusEnumQuery)
	if err != nil {
		log.Println(err)
		log.Println("Error creating status enum continuing with the program...")
	}
	_, err = store.db.Exec(createTaskTableQuery)
	if err != nil {
		return err
	}
	return nil
}

// SELECT * from Tasks WHERE projectId=$1
func (store *PostgresTaskStore) GetTasks(projectId int) ([]Task, error) {
	rows, err := store.db.Query("SELECT * from Tasks WHERE projectId=$1", projectId)
	if err != nil {
		log.Println("Error getting tasks from database: ", err)
		return nil, err
	}
	var tasks []Task
	for rows.Next() {
		task := Task{}
		err = rows.Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.ProjectId)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)

	}
	fmt.Println(tasks)
	return tasks, nil
}

func (store *PostgresTaskStore) GetTaskById(id string) (Task, error) {
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

func (store *PostgresTaskStore) CreateTask(p *CreateTaskRequest) error {
	_, err := store.db.Exec("INSERT INTO Tasks (title, description, status, projectId) VALUES($1, $2, $3, $4)", p.Title, p.Description, p.Status, p.ProjectId)
	if err != nil {
		return err
	}

	return nil
}

func (store *PostgresTaskStore) UpdateTask(id string, p *CreateTaskRequest) error {
	_, err := store.db.Exec("UPDATE Tasks SET title=$1, description=$2, status=$3 WHERE id=$4", p.Title, p.Description, p.Status, id)
	if err != nil {
		return err
	}
	return nil
}

func (store *PostgresTaskStore) DeleteTask(id string) error {
	_, err := store.db.Exec("DELETE FROM Tasks WHERE id=$1", id)
	if err != nil {
		return err
	}
	return nil
}

// This is the interface that the service will use to interact with the database
type ProjectStorage interface {
	GetProjects() ([]Project, error)
	GetProjectById(string) (Project, error)
	CreateProject(*CreateProjectRequest) error
	UpdateProject(string, *CreateProjectRequest) error
	DeleteProject(string) error
}

type PostgresProjectStore struct {
	db *sql.DB
}

func NewPostgresProjectStore() (*PostgresProjectStore, error) {
	connStr := "user=ttracker dbname=ttracker password=ttracker sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresProjectStore{
		db: db,
	}, nil
}
func (store *PostgresProjectStore) createProjectTable() error {

	_, err := store.db.Exec(createPriorityEnumQuery)
	if err != nil {
		log.Println(err)
		log.Println("Error creating priority enum continuing with the program...")
	}

	_, err = store.db.Exec(createProjectTableQuery)
	if err != nil {
		return err
	}

	return nil
}

func (store *PostgresProjectStore) Init() error {
	return store.createProjectTable()
}

func (store *PostgresProjectStore) GetProjects() ([]Project, error) {
	rows, err := store.db.Query("SELECT * from Projects")
	if err != nil {
		return nil, err
	}
	var projects []Project
	for rows.Next() {
		project := Project{}
		err = rows.Scan(&project.Id, &project.Title, &project.Description, &project.Priority, &project.CreatedAt)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)

	}
	return projects, nil
}

func (store *PostgresProjectStore) GetProjectById(id string) (Project, error) {
	rows, err := store.db.Query("SELECT * from Projects WHERE id=$1", id)
	var project Project
	err = rows.Scan(&project.Id, &project.Title, &project.Description, &project.Priority, &project.CreatedAt)
	if err != nil {
		return Project{}, err
	}
	return project, nil
}

func (store *PostgresProjectStore) CreateProject(p *CreateProjectRequest) error {
	_, err := store.db.Exec("INSERT INTO Projects (title, description, priority) VALUES($1, $2, $3)", p.Title, p.Description, p.Priority)
	if err != nil {
		return err
	}
	return nil
}

func (store *PostgresProjectStore) UpdateProject(id string, p *CreateProjectRequest) error {
	_, err := store.db.Exec("UPDATE Projects SET title=$1, description=$2, priority=$3 WHERE id=$4", p.Title, p.Description, p.Priority, id)
	if err != nil {
		return err
	}
	return nil
}

func (store *PostgresProjectStore) DeleteProject(id string) error {
	_, err := store.db.Exec("DELETE FROM Projects WHERE id=$1", id)
	if err != nil {
		return err
	}
	return nil
}
