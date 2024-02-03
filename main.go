package main

import (
	"log"
)

func main() {
	store, err := NewPostgresTaskStore()
	service := NewTaskService(store)
	if err != nil {
		log.Fatalln(err)
	}
	if err := store.Init(); err != nil {
		log.Fatalln(err)
	}
	server := NewApiServer(":8000", service)
	server.Run()

}
