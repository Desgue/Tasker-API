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
	connStr, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		connStr = os.Getenv("LOCAL_DB")
	}
	hostPort, ok := os.LookupEnv("PORT")
	if !ok {
		hostPort = "8000"
	}

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
	server := NewApiServer(hostPort, taskService, projectService)
	server.Run()

}
