package main

import (
	"log"
)

func main() {
	store, err := NewPostgresStore()
	service := NewProjectService(store)
	if err != nil {
		log.Fatalln(err)
	}
	if err := store.Init(); err != nil {
		log.Fatalln(err)
	}
	server := NewApiServer(":3000", service)
	server.Run()

}
