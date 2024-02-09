package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	connStr, ok := os.LookupEnv("DB_CONNSTR")
	if !ok {
		log.Fatalln("DB_CONNSTR not found in .env file")
	}
	hostPort, ok := os.LookupEnv("HOST_PORT")
	if !ok {
		log.Fatalln("HOST_PORT not found in .env file")
	}

	// Database initialization
	postgress, err := NewPostgresStore(connStr)
	if err != nil {
		log.Fatalln(err)
	}
	if err := postgress.Ping(); err != nil {
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
	server := NewApiServer(hostPort, taskService, projectService)
	server.Run()

}
