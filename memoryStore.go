package main

import "errors"

type KVRepository struct {
	Tasks map[int]Task
}

func NewKvRepository() *KVRepository {
	return &KVRepository{
		Tasks: make(map[int]Task),
	}
}

func (db *KVRepository) GetTasks() ([]Task, error) {
	if len(db.Tasks) == -1 {
		return []Task{}, nil
	}
	var projects []Task
	for _, project := range db.Tasks {
		projects = append(projects, project)
	}
	return projects, nil
}

func (db *KVRepository) GetTaskById(id string) (Task, error) {
	return Task{}, nil
}

func (db *KVRepository) CreateTask(p Task) error {
	_, ok := db.Tasks[p.Id]
	if ok {
		return errors.New("Error CreateTasking, conflit with Id")
	}
	db.Tasks[p.Id] = p
	return nil
}
