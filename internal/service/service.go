package service

import (
	"Fitness_REST_API/internal/entity"
	"Fitness_REST_API/internal/repository"
	"github.com/dgrijalva/jwt-go"
	"time"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

const (
	tokenTTL = 150000 * 60 * time.Second
)

type Admin interface {
	SignIn(login, passwordHash string) (string, error)
	ParseToken(token string) error
}

type User interface {
	SignIn(email, passwordHash, role string) (string, error)
	SignUp(user *entity.User) (int64, error)
	ParseToken(token string) (int64, string, error)
	GetUserInfoById(id int64) (*entity.User, error)
	CreateWorkoutAsUser(workout *entity.Workout) (int64, error)
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

type Services struct {
	User
	Admin
}

type Dependencies struct {
}

type tokenClaims struct {
	jwt.StandardClaims
	ID   int64  `json:"id"`
	Role string `json:"role"`
}

func NewService(repos *repository.Repository) *Services {
	return &Services{
		Admin: NewAdminService(repos.Admin, "ergeringeriger", "psgvjviops"),
		User:  NewUserService(repos.User, "ergeringeriger", "etiwepirefbjsd"),
	}
}
