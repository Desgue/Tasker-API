package repo

import (
	"database/sql"
	"log"

	"github.com/Desgue/ttracker-api/internal/domain"
	_ "github.com/lib/pq"
)

type PostgresTaskStore struct {
	DB *sql.DB
}

func NewPostgresTaskStore(DB *sql.DB) *PostgresTaskStore {
	return &PostgresTaskStore{
		DB: DB,
	}
}

// SELECT * from Tasks WHERE projectId=$1
func (store *PostgresTaskStore) GetTasks(projectId int) ([]domain.Task, error) {
	rows, err := store.DB.Query("SELECT * from Tasks WHERE projectId=$1", projectId)
	if err != nil {
		log.Println("Error getting tasks from database: ", err)
		return nil, err
	}
	var tasks []domain.Task
	for rows.Next() {
		task := domain.Task{}
		err = rows.Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.ProjectId)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)

	}

	return tasks, nil
}

func (store *PostgresTaskStore) GetTaskById(id string) (domain.Task, error) {
	rows, err := store.DB.Query("SELECT * from Tasks WHERE id=$1", id)
	if err != nil {
		log.Println("Error getting task from database: ", err)
		return domain.Task{}, err
	}
	var task domain.Task
	for rows.Next() {
		task = domain.Task{}
		err = rows.Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.CreatedAt)
		if err != nil {
			return domain.Task{}, err
		}
	}
	return task, nil
}

func (store *PostgresTaskStore) CreateTask(p *domain.CreateTaskRequest) error {
	_, err := store.DB.Exec("INSERT INTO Tasks (title, description, status, projectId) VALUES($1, $2, $3, $4)", p.Title, p.Description, p.Status, p.ProjectId)
	if err != nil {
		return err
	}

	return nil
}

func (store *PostgresTaskStore) UpdateTask(id string, p *domain.CreateTaskRequest) error {
	_, err := store.DB.Exec("UPDATE Tasks SET title=$1, description=$2, status=$3 WHERE id=$4", p.Title, p.Description, p.Status, id)
	if err != nil {
		return err
	}
	return nil
}

func (store *PostgresTaskStore) DeleteTask(id string) error {
	_, err := store.DB.Exec("DELETE FROM Tasks WHERE id=$1", id)
	if err != nil {
		return err
	}
	return nil
}
