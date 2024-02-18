package repo

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// This is the struct that will hold the database connection
type PostgresUserStore struct {
	DB *sql.DB
}

func NewPostgresUserStore(DB *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{
		DB: DB,
	}
}

func (store *PostgresUserStore) CheckUser(cognitoId string) (bool, error) {
	var id string
	err := store.DB.QueryRow("SELECT cognitoId from Users where cognitoId=$1", cognitoId).Scan(&id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (store *PostgresUserStore) CreateUser(cognitoId string) error {
	_, err := store.DB.Exec("INSERT INTO Users (cognitoId) VALUES($1)", cognitoId)
	if err != nil {
		return err
	}
	return nil
}
