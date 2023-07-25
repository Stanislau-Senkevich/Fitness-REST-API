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
	GetUserInfoById(id int64) (*entity.User, error)
	CreateWorkoutAsUser(*entity.Workout) (int64, error)
	UpdateWorkout(workoutId, userId int64, update *entity.UpdateWorkout) error
	GetUserWorkouts(id int64) ([]*entity.Workout, error)
	GetWorkoutById(workoutId, userId int64) (*entity.Workout, error)
	DeleteWorkout(workoutId, userId int64) error
	GetTrainers() ([]*entity.User, error)
	GetTrainerById(id int64) (*entity.User, error)
	SendRequestToTrainer(trainerId, userId int64) (int64, error)
	EndPartnershipWithTrainer(trainerId, userId int64) (int64, error)
	GetUserPartnerships(userId int64) ([]*entity.Partnership, error)

	GetTrainerUsers(trainerId int64) ([]*entity.User, error)
	GetTrainerRequests(trainerId int64) ([]*entity.Request, error)
	GetTrainerUserById(trainerId, userId int64) (*entity.User, error)
	GetTrainerRequestById(trainerId, requestId int64) (*entity.Request, error)
	InitPartnershipWithUser(trainerId, userId int64) (int64, error)
	EndPartnershipWithUser(trainerId, userId int64) (int64, error)
	AcceptRequest(trainerId, requestId int64) (int64, error)
	DenyRequest(trainerId, requestId int64) error
	CreateWorkoutAsTrainer(workout *entity.Workout) (int64, error)
	GetTrainerWorkouts(trainerId int64) ([]*entity.Workout, error)
	GetTrainerWorkoutsWithUser(trainerId, userId int64) ([]*entity.Workout, error)
}
