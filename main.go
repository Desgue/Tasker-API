package main

import (
	"log"

	"github.com/Desgue/ttracker-api/internal/api"
	repo "github.com/Desgue/ttracker-api/internal/repository"
	svc "github.com/Desgue/ttracker-api/internal/services"
	"github.com/Desgue/ttracker-api/internal/util"
)

func main() {
	// Load variables from .env file

	util.LoadENV()

	// Database initialization
	postgress, err := repo.NewPostgresStore(util.ConnStr)
	if err != nil {
		log.Fatalln(err)
	}
	if err := postgress.Ping(); err != nil {
		log.Fatalln(err)
	}
	postgress.Init()

	// User initialization
	//userStore := repo.NewPostgresUserStore(postgress.DB)

	// Project initialization
	projectStore := repo.NewPostgresProjectStore(postgress.DB)

	projectService := svc.NewProjectService(projectStore)

	// Task initialization
	taskStore := repo.NewPostgresTaskStore(postgress.DB)

	taskService := svc.NewTaskService(taskStore)

	// API server initialization
	server := api.NewApiServer(util.ListenAddr, taskService, projectService)
	server.Run()

}
