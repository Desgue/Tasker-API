package main

import (
	"database/sql"
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
	userId SMALLINT NOT NULL REFERENCES Users(id),
	title varchar(255),
	description text,
	priority priority DEFAULT 'Low',
	createdAt TIMESTAMPTZ NOT NULL DEFAULT NOW()
);`
	createUserTableQuery = `
	CREATE TABLE IF NOT EXISTS Users (
	id SMALLINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	cognitoId varchar(255) NOT NULL,
	createdAt TIMESTAMPTZ NOT NULL DEFAULT NOW()
);`
)

type PostgresStore struct {
	connStr string
	db      *sql.DB
}

func (store *PostgresStore) Ping() error {
	if err := store.db.Ping(); err != nil {
		log.Println("Error pinging the database: ", err)
		return err
	}
	return nil
}

func NewPostgresStore(connStr string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &PostgresStore{
		connStr: connStr,
		db:      db,
	}, nil
}

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

func NewPostgresTaskStore(db *sql.DB) *PostgresTaskStore {
	return &PostgresTaskStore{
		db: db,
	}
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
	GetProjects(userId string) ([]Project, error)
	GetProjectById(projectId, cognitoId string) (Project, error)
	CreateProject(*CreateProjectRequest) error
	UpdateProject(string, *CreateProjectRequest) error
	DeleteProject(projectId, cognitoId string) error
}

type PostgresProjectStore struct {
	db *sql.DB
}

func NewPostgresProjectStore(db *sql.DB) *PostgresProjectStore {
	return &PostgresProjectStore{
		db: db,
	}
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

func (store *PostgresProjectStore) GetProjects(cognitoId string) ([]Project, error) {
	// Perform a joing with the users id to retriev all projects associated with the users cognitoId
	// Then select all projects wich matches the user cognitoId

	rows, err := store.db.Query(`
	SELECT 
	Projects.id, 
	Projects.title, 
	Projects.description, 
	Projects.priority, 
	Projects.createdAt 
	FROM 
	Projects 
	INNER JOIN Users ON Projects.userId=Users.id 
	WHERE Users.cognitoId=$1`,
		cognitoId)
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

func (store *PostgresProjectStore) GetProjectById(projectId, cognitoId string) (Project, error) {
	// Perform a join with the users id to retriev the project associated id and the users cognitoId
	rows, err := store.db.Query(`
	SELECT
	Projects.id,
	Users.cognitoId,
	Projects.title,
	Projects.description,
	Projects.priority,
	Projects.createdAt
	FROM
	Projects
	INNER JOIN Users ON Projects.userId=Users.id
	WHERE Projects.id=$1 AND Users.cognitoId=$2`, projectId, cognitoId)
	if err != nil {
		return Project{}, err
	}
	var project Project
	for rows.Next() {
		err = rows.Scan(&project.Id, &project.Title, &project.Description, &project.Priority, &project.CreatedAt)
		if err != nil {
			return Project{}, err
		}
	}
	return project, nil

}

func (store *PostgresProjectStore) CreateProject(p *CreateProjectRequest) error {
	// Create a new project and associate it with the user cognitoId
	// The user cognitoId is used to retrieve the user id from the Users table
	// Then the user id is used to associate the project with the user
	row, err := store.db.Query("SELECT id from Users where cognitoId=$1", p.UserCognitoId)
	if err != nil {
		return err
	}
	var userId int
	for row.Next() {
		err = row.Scan(&userId)
		if err != nil {
			return err
		}
	}

	_, err = store.db.Exec(`
	INSERT INTO Projects 
	(title, description, priority, userId) 
	VALUES($1, $2, $3, $4)`,
		p.Title, p.Description, p.Priority, userId)

	if err != nil {
		return err
	}
	return nil
}

func (store *PostgresProjectStore) UpdateProject(id string, p *CreateProjectRequest) error {
	// Update the project with the new values if the project id matches the id of the user
	// The user id is used to associate the project with the user and must be retrieved using the cognitoId
	row, err := store.db.Query("SELECT id from Users where cognitoId=$1", p.UserCognitoId)
	if err != nil {
		return err
	}
	var userId int
	for row.Next() {
		err = row.Scan(&userId)
		if err != nil {
			return err
		}
	}
	_, err = store.db.Exec(`
	UPDATE Projects 
	SET title=$1, description=$2, priority=$3 
	WHERE id=$4 AND userId=$5`,
		p.Title, p.Description, p.Priority, id, userId)

	if err != nil {
		return err
	}
	return nil

}

func (store *PostgresProjectStore) DeleteProject(projectId, cognitoId string) error {
	// Delete all tasks asscoiated with the project id
	// Then delete the project with the project id where the user id matches the cognito id

	row, err := store.db.Query("SELECT id from Users where cognitoId=$1", cognitoId)
	if err != nil {
		return err
	}
	var userId int
	for row.Next() {
		err = row.Scan(&userId)
		if err != nil {
			return err
		}
	}

	_, err = store.db.Exec(`
	DELETE FROM Tasks WHERE projectId=$1`,
		projectId)
	if err != nil {
		return err
	}
	_, err = store.db.Exec(`
	DELETE FROM Projects
	WHERE id=$1 AND userId=$2`,
		projectId, userId)

	if err != nil {
		return err
	}
	return nil
}

// This is the interface that that will define the behavior to interact with the database
type UserStorage interface {
	CheckUser(string) (bool, error)
	CreateUser(string) error
}

// This is the struct that will hold the database connection
type PostgresUserStore struct {
	db *sql.DB
}

func (store *PostgresUserStore) Init() error {
	_, err := store.db.Exec(createUserTableQuery)
	if err != nil {
		return err
	}
	return nil
}
func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{
		db: db,
	}
}

func (store *PostgresUserStore) CheckUser(cognitoId string) (bool, error) {
	var id string
	err := store.db.QueryRow("SELECT cognitoId from Users where cognitoId=$1", cognitoId).Scan(&id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (store *PostgresUserStore) CreateUser(cognitoId string) error {
	_, err := store.db.Exec("INSERT INTO Users (cognitoId) VALUES($1)", cognitoId)
	if err != nil {
		return err
	}
	return nil
}
