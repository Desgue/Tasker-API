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

	// User initialization
	userStore := repo.NewPostgresUserStore(postgress.DB)
	if err := userStore.Init(); err != nil {
		log.Fatalln(err)
	}

	// Project initialization
	projectStore := repo.NewPostgresProjectStore(postgress.DB)
	if err := projectStore.Init(); err != nil {
		log.Fatalln(err)
	}
	projectService := svc.NewProjectService(projectStore)

	// Task initialization
	taskStore := repo.NewPostgresTaskStore(postgress.DB)
	if err := taskStore.Init(); err != nil {
		log.Fatalln(err)
	}
	taskService := svc.NewTaskService(taskStore)

	// API server initialization
	server := api.NewApiServer(util.ListenAddr, taskService, projectService)
	server.Run()

}
