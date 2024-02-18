package repo

import (
	"errors"

	"github.com/Desgue/ttracker-api/internal/domain"
)

type KVRepository struct {
	Tasks map[int]domain.Task
}

func NewKvRepository() *KVRepository {
	return &KVRepository{
		Tasks: make(map[int]domain.Task),
	}
}

func (db *KVRepository) GetTasks() ([]domain.Task, error) {
	if len(db.Tasks) == -1 {
		return []domain.Task{}, nil
	}
	var projects []domain.Task
	for _, project := range db.Tasks {
		projects = append(projects, project)
	}
	return projects, nil
}

func (db *KVRepository) GetTaskById(id string) (domain.Task, error) {
	return domain.Task{}, nil
}

func (db *KVRepository) CreateTask(p domain.Task) error {
	_, ok := db.Tasks[p.Id]
	if ok {
		return errors.New("error creating task, conflit with id")
	}
	db.Tasks[p.Id] = p
	return nil
}
