package repository

import (
	"Fitness_REST_API/internal/entity"
	"Fitness_REST_API/internal/repository/postgres"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	Admin
	User
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Admin: postgres.NewAdminRepository(db),
		User:  postgres.NewUserRepository(db),
	}
}

type Admin interface {
	Authorize(login, passwordHash string) error
}

type User interface {
	Authorize(email, passwordHash, role string) (int64, error)
	CreateUser(user *entity.User) (int64, error)
	GetUser(id int64) (*entity.User, error)
	HasApprovedPartnership(trainerId, userId int64) bool

	CreateWorkoutAsUser(*entity.Workout) (int64, error)
	UpdateWorkout(workoutId, userId int64, update *entity.UpdateWorkout) error
	GetAllUserWorkouts(id int64) ([]*entity.Workout, error)
	GetWorkoutById(workoutId, userId int64) (*entity.Workout, error)
	DeleteWorkout(workoutId, userId int64) error
	GetAllTrainers() ([]*entity.User, error)
	GetTrainerById(id int64) (*entity.User, error)
	SendRequestToTrainer(trainerId, userId int64) (int64, error)
	EndPartnershipWithTrainer(trainerId, userId int64) (int64, error)
	GetUserPartnerships(userId int64) ([]*entity.Partnership, error)
}
