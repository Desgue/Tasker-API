package main

import "fmt"

func main() {
	fmt.Println("Hello World")

	server := NewApiServer(":3000")
	server.Run()

}
