package main

import (
	"Fitness_REST_API/internal/config"
	"Fitness_REST_API/internal/handler"
	"Fitness_REST_API/internal/repository"
	"Fitness_REST_API/internal/repository/postgres"
	"Fitness_REST_API/internal/server"
	"Fitness_REST_API/internal/service"
	"github.com/sirupsen/logrus"
	"log"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	cfg, err := config.InitConfig()
	if err != nil {
		logrus.Fatalf("error due reading config: %s", err.Error())
	}

	db, err := postgres.InitPostgresDB(cfg)
	if err != nil {
		logrus.Fatalf("error due initializing database: %s", err.Error())
	}
	defer func() {
		err = db.Close()
		if err != nil {
			logrus.Fatalf("error due closing db: %s", err.Error())
		}
	}()

	srv := new(server.Server)
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	err = srv.Run(cfg.Port, handlers.InitRoutes())
	if err != nil {
		log.Fatalf("error due running server: %s", err.Error())
	}
}
