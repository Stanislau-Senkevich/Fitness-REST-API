package service

import (
	"Fitness_REST_API/internal/repository"
	"time"
)

const (
	tokenTTL = 15 * 60 * time.Second
)

type Admin interface {
	SignIn(login, passwordHash string) (string, error)
	ParseToken(token string) error
}

type User interface {
	SignIn(login, passwordHash string) (string, error)
	ParseToken(token string) error
}

type Trainer interface {
	SignIn(login, passwordHash string) (string, error)
	ParseToken(token string) error
}

type Workout interface {
}

type Partnership interface {
}

type Services struct {
	User
	Workout
	Admin
	Trainer
	Partnership
}

type Dependencies struct {
}

func NewService(repos *repository.Repository) *Services {
	return &Services{
		Admin:   NewAdminService(repos.Admin, "ergeringeriger", "psgvjviops"),
		Trainer: NewAdminService(repos.Trainer, "ergeringeriger", "dfohkjdfsdf"),
		User:    NewUserService(repos.User, "ergeringeriger", "etiwepirefbjsd"),
	}
}
