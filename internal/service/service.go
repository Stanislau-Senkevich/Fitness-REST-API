package service

import (
	"Fitness_REST_API/internal/entity"
	"Fitness_REST_API/internal/repository"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	tokenTTL = 150 * 60 * time.Second
)

type Admin interface {
	SignIn(login, passwordHash string) (string, error)
	ParseToken(token string) error
}

type User interface {
	SignIn(email, passwordHash string) (string, error)
	SignUp(user *entity.User) (int64, error)
	ParseToken(token string) (int64, error)
	GetUser(id int64) (*entity.User, error)
	CreateWorkout(workout *entity.UserWorkout) (int64, error)
}

type Trainer interface {
	SignIn(login, passwordHash string) (string, error)
	ParseToken(token string) (int64, error)
}

type Services struct {
	User
	Admin
	Trainer
}

type Dependencies struct {
}

type tokenClaims struct {
	jwt.StandardClaims
	ID int64 `json:"id"`
}

func NewService(repos *repository.Repository) *Services {
	return &Services{
		Admin:   NewAdminService(repos.Admin, "ergeringeriger", "psgvjviops"),
		Trainer: NewTrainerService(repos.Trainer, "ergeringeriger", "dfohkjdfsdf"),
		User:    NewUserService(repos.User, "ergeringeriger", "etiwepirefbjsd"),
	}
}
