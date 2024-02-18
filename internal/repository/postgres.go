package repo

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
	DB      *sql.DB
}

func (store *PostgresStore) Ping() error {
	if err := store.DB.Ping(); err != nil {
		log.Println("Error pinging the database: ", err)
		return err
	}
	return nil
}

func (store *PostgresStore) Init() {
	store.createEnums()
	store.createTables()
}

func (store *PostgresStore) createEnums() {
	_, err := store.DB.Exec(createStatusEnumQuery)
	if err != nil {
		log.Println(err)
		log.Println("Error creating status enum continuing with the program...")
	}
	_, err = store.DB.Exec(createPriorityEnumQuery)
	if err != nil {
		log.Println(err)
		log.Println("Error creating priority enum continuing with the program...")
	}

}

func (store *PostgresStore) createTables() {
	_, err := store.DB.Exec(createUserTableQuery)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = store.DB.Exec(createProjectTableQuery)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = store.DB.Exec(createTaskTableQuery)
	if err != nil {
		log.Fatalln(err)
	}
}

func NewPostgresStore(connStr string) (*PostgresStore, error) {
	DB, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &PostgresStore{
		connStr: connStr,
		DB:      DB,
	}, nil
}
