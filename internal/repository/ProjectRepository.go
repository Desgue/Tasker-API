package repo

import (
	"database/sql"

	"github.com/Desgue/ttracker-api/internal/domain"
	_ "github.com/lib/pq"
)

type PostgresProjectStore struct {
	DB *sql.DB
}

func NewPostgresProjectStore(DB *sql.DB) *PostgresProjectStore {
	return &PostgresProjectStore{
		DB: DB,
	}
}

func (store *PostgresProjectStore) GetProjects(cognitoId string) ([]domain.Project, error) {
	// Perform a joing with the users id to retriev all projects associated with the users cognitoId
	// Then select all projects wich matches the user cognitoId

	rows, err := store.DB.Query(`
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
	var projects []domain.Project
	for rows.Next() {
		project := domain.Project{}
		err = rows.Scan(&project.Id, &project.Title, &project.Description, &project.Priority, &project.CreatedAt)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil

}

func (store *PostgresProjectStore) GetProjectById(projectId, cognitoId string) (domain.Project, error) {
	// Perform a join with the users id to retriev the project associated id and the users cognitoId
	rows, err := store.DB.Query(`
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
		return domain.Project{}, err
	}
	var project domain.Project
	for rows.Next() {
		err = rows.Scan(&project.Id, &project.Title, &project.Description, &project.Priority, &project.CreatedAt)
		if err != nil {
			return domain.Project{}, err
		}
	}
	return project, nil

}

func (store *PostgresProjectStore) CreateProject(p *domain.CreateProjectRequest) error {
	// Create a new project and associate it with the user cognitoId
	// The user cognitoId is used to retrieve the user id from the Users table
	// Then the user id is used to associate the project with the user
	row, err := store.DB.Query("SELECT id from Users where cognitoId=$1", p.UserCognitoId)
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

	_, err = store.DB.Exec(`
	INSERT INTO Projects 
	(title, description, priority, userId) 
	VALUES($1, $2, $3, $4)`,
		p.Title, p.Description, p.Priority, userId)

	if err != nil {
		return err
	}
	return nil
}

func (store *PostgresProjectStore) UpdateProject(id string, p *domain.CreateProjectRequest) error {
	// Update the project with the new values if the project id matches the id of the user
	// The user id is used to associate the project with the user and must be retrieved using the cognitoId
	row, err := store.DB.Query("SELECT id from Users where cognitoId=$1", p.UserCognitoId)
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
	_, err = store.DB.Exec(`
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

	row, err := store.DB.Query("SELECT id from Users where cognitoId=$1", cognitoId)
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

	_, err = store.DB.Exec(`
	DELETE FROM Tasks WHERE projectId=$1`,
		projectId)
	if err != nil {
		return err
	}
	_, err = store.DB.Exec(`
	DELETE FROM Projects
	WHERE id=$1 AND userId=$2`,
		projectId, userId)

	if err != nil {
		return err
	}
	return nil
}
