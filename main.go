package main

import (
	"log"
)

func main() {
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatalln(err)
	}
	server := NewApiServer(":3000", store)
	server.Run()

}
