package main

import (
	"Fitness_REST_API/internal/config"
	"Fitness_REST_API/internal/handler"
	"Fitness_REST_API/internal/repository"
	"Fitness_REST_API/internal/repository/postgres"
	"Fitness_REST_API/internal/server"
	"Fitness_REST_API/internal/service"
	"context"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

// @title Fitness REST API
// @version 1.0
// @description API Server for Fitness application

// @host droplet.senkevichdev.work:8001
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

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

	go func() {
		err = srv.Run(cfg.Port, handlers.InitRoutes())
		if err != nil {
			logrus.Fatalf("error due running server: %s", err.Error())
		}
	}()

	closeChan := make(chan os.Signal, 1)
	signal.Notify(closeChan, syscall.SIGTERM, syscall.SIGINT)
	<-closeChan

	if err = srv.ShutDown(context.Background()); err != nil {
		logrus.Fatalf("Error due shutdown: %s", err.Error())
	}
}
