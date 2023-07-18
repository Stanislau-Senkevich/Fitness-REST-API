package main

import (
	"Fitness_REST_API/internal/config"
	"Fitness_REST_API/internal/handler"
	"Fitness_REST_API/internal/repository"
	"Fitness_REST_API/internal/repository/postgres"
	"Fitness_REST_API/internal/server"
	"Fitness_REST_API/internal/service"
	"log"
)

func main() {

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := postgres.InitPostgresDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	srv := new(server.Server)
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	err = srv.Run(cfg.Port, handlers.InitRoutes())
	if err != nil {
		log.Fatal(err)
	}
}
