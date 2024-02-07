package main

import (
	"log"
)

const (
	connStr = "user=postgres password=postgres dbname=postgres sslmode=disable"
)

func main() {
	// Database initialization
	postgress, err := NewPostgresStore(connStr)
	if err != nil {
		log.Fatalln(err)
	}
	if err = postgress.Ping(); err != nil {
		log.Fatalln(err)
	}

	// Project initialization
	projectStore := NewPostgresProjectStore(postgress.db)
	if err := projectStore.Init(); err != nil {
		log.Fatalln(err)
	}
	projectService := NewProjectService(projectStore)

	// Task initialization
	taskStore := NewPostgresTaskStore(postgress.db)
	if err := taskStore.Init(); err != nil {
		log.Fatalln(err)
	}
	taskService := NewTaskService(taskStore)

	// User initialization
	userStore := NewPostgresUserStore(postgress.db)
	userService := NewUserService(userStore)

	// API server initialization
	server := NewApiServer(":8000", taskService, projectService, userService)
	server.Run()

}
