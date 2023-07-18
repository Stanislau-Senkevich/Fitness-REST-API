package repository

import (
	"Fitness_REST_API/internal/repository/postgres"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	Admin
	Trainer
	User
	Workout
	Partnership
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Admin:       postgres.NewAdminRepository(db),
		User:        postgres.NewUserRepository(db),
		Trainer:     postgres.NewTrainerRepository(db),
		Workout:     postgres.NewWorkoutRepository(db),
		Partnership: postgres.NewPartnershipRepository(db),
	}
}

type Admin interface {
	Authorize(login, passwordHash string) error
}

type User interface {
	Authorize(login, passwordHash string) error
}

type Trainer interface {
	Authorize(login, passwordHash string) error
}

type Workout interface {
}

type Partnership interface {
}
