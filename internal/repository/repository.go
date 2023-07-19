package repository

import (
	"Fitness_REST_API/internal/entity"
	"Fitness_REST_API/internal/repository/postgres"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	Admin
	Trainer
	User
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Admin:   postgres.NewAdminRepository(db),
		User:    postgres.NewUserRepository(db),
		Trainer: postgres.NewTrainerRepository(db),
	}
}

type Admin interface {
	Authorize(login, passwordHash string) error
}

type User interface {
	Authorize(email, passwordHash string) (int64, error)
	CreateUser(user *entity.User) (int64, error)
	GetUser(id int64) (*entity.User, error)
	CreateWorkout(*entity.UserWorkout) (int64, error)
	CheckPartnership(userId, trainerId int64) error
}

type Trainer interface {
	Authorize(login, passwordHash string) (int64, error)
}
