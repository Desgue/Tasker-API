package main

import (
	"log"
)

func main() {

	projectStore, err := NewPostgresProjectStore()
	if err != nil {
		log.Fatalln(err)
	}
	if err := projectStore.Init(); err != nil {
		log.Fatalln(err)
	}
	projectService := NewProjectService(projectStore)

	taskStore, err := NewPostgresTaskStore()
	if err != nil {
		log.Fatalln(err)
	}
	if err := taskStore.Init(); err != nil {
		log.Fatalln(err)
	}
	taskService := NewTaskService(taskStore)

	server := NewApiServer(":8000", taskService, projectService)
	server.Run()

}
