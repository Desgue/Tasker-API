package main

import (
	"log"
)

var (
	connStr         string
	hostPort        string
	listenAddr      string
	ok              bool
	cognito_jwk_url string
	cognito_issuer  string
)

func main() {
	// Load variables from .env file

	loadENV()

	// Database initialization
	postgress, err := NewPostgresStore(connStr)
	if err != nil {
		log.Fatalln(err)
	}
	if err := postgress.Ping(); err != nil {
		log.Fatalln(err)
	}

	// User initialization
	userStore := NewPostgresUserStore(postgress.db)
	if err := userStore.Init(); err != nil {
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

	// API server initialization
	server := NewApiServer(listenAddr, taskService, projectService)
	server.Run()

}
