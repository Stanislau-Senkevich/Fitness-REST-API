package main

import (
	"Fitness_REST_API/internal/server"
	"log"
)

func main() {
	srv := new(server.Server)
	err := srv.Run("8000")
	if err != nil {
		log.Fatal(err)
	}
}
