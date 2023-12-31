package service

import (
	"Fitness_REST_API/internal/entity"
	"Fitness_REST_API/internal/repository"
	"crypto/sha1"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type UserService struct {
	repo       repository.User
	hashSalt   string
	signingKey []byte
}

func NewUserService(repos repository.User, hashSalt string, signingKey string) *UserService {
	return &UserService{repo: repos, hashSalt: hashSalt, signingKey: []byte(signingKey)}
}

func (s *UserService) SignIn(email, password string, role entity.Role) (string, error) {
	id, err := s.repo.Authorize(email, s.GetPasswordHash(password), role)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		id,
		role,
	})

	return token.SignedString(s.signingKey)
}

func (s *UserService) SignUp(user *entity.User) (int64, error) {
	user.PasswordHash = s.GetPasswordHash(user.PasswordHash)
	return s.repo.CreateUser(user, entity.UserRole)
}

func (s *UserService) ParseToken(token string) (int64, entity.Role, error) {
	t, err := jwt.ParseWithClaims(token, &tokenClaims{}, func(token *jwt.Token) (i interface{}, err error) { //nolint
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.signingKey, nil
	})

	if err != nil {
		return -1, "", err
	}

	claims, ok := t.Claims.(*tokenClaims)
	if !ok {
		return -1, "", fmt.Errorf("error get user claims from token")
	}

	return claims.ID, claims.Role, nil
}

func (s *UserService) InitUpdateUser(userId int64, update *entity.UserUpdate) error {
	user, err := s.GetUserInfoById(userId)
	if err != nil {
		return err
	}

	if update.Email == "" {
		update.Email = user.Email
	}
	if update.Password == "" {
		update.Password = user.PasswordHash
	} else {
		update.Password = s.GetPasswordHash(update.Password)
	}
	if update.Name == "" {
		update.Name = user.Name
	}
	if update.Surname == "" {
		update.Surname = user.Surname
	}
	if update.Role == "" {
		update.Role = user.Role
	}
	return nil
}

func (s *UserService) FormatUpdateWorkout(input *entity.UpdateWorkout, workoutId, userId int64) error {
	workout, err := s.GetWorkoutById(workoutId, userId)
	if err != nil {
		return err
	}

	if input.Title == "" {
		input.Title = workout.Title
	}

	if input.Date.IsZero() {
		input.Date = workout.Date
	}
	return nil
}

func (s *UserService) GetUserInfoById(id int64) (*entity.User, error) {
	return s.repo.GetUserInfoById(id)
}

func (s *UserService) CreateWorkoutAsUser(workout *entity.Workout) (int64, error) {
	return s.repo.CreateWorkoutAsUser(workout)
}

func (s *UserService) UpdateWorkout(workoutId, userId int64, update *entity.UpdateWorkout) error {
	return s.repo.UpdateWorkout(workoutId, userId, update)
}

func (s *UserService) GetUserWorkouts(id int64) ([]*entity.Workout, error) {
	return s.repo.GetUserWorkouts(id)
}

func (s *UserService) GetWorkoutById(workoutId, userId int64) (*entity.Workout, error) {
	return s.repo.GetWorkoutById(workoutId, userId)
}

func (s *UserService) DeleteWorkout(workoutId, userId int64) error {
	return s.repo.DeleteWorkout(workoutId, userId)
}

func (s *UserService) GetTrainers() ([]*entity.User, error) {
	return s.repo.GetTrainers()
}

func (s *UserService) GetTrainerById(id int64) (*entity.User, error) {
	return s.repo.GetTrainerById(id)
}

func (s *UserService) SendRequestToTrainer(trainerId, userId int64) (int64, error) {
	return s.repo.SendRequestToTrainer(trainerId, userId)
}

func (s *UserService) EndPartnershipWithTrainer(trainerId, userId int64) (int64, error) {
	return s.repo.EndPartnershipWithTrainer(trainerId, userId)
}

func (s *UserService) GetUserPartnerships(userId int64) ([]*entity.Partnership, error) {
	return s.repo.GetUserPartnerships(userId)
}

func (s *UserService) GetTrainerUsers(trainerId int64) ([]*entity.User, error) {
	return s.repo.GetTrainerUsers(trainerId)
}

func (s *UserService) GetTrainerRequests(trainerId int64) ([]*entity.Request, error) {
	return s.repo.GetTrainerRequests(trainerId)
}

func (s *UserService) GetTrainerUserById(trainerId, userId int64) (*entity.User, error) {
	return s.repo.GetTrainerUserById(trainerId, userId)
}

func (s *UserService) GetTrainerRequestById(trainerId, requestId int64) (*entity.Request, error) {
	return s.repo.GetTrainerRequestById(trainerId, requestId)
}

func (s *UserService) InitPartnershipWithUser(trainerId, userId int64) (int64, error) {
	return s.repo.InitPartnershipWithUser(trainerId, userId)
}

func (s *UserService) EndPartnershipWithUser(trainerId, userId int64) (int64, error) {
	return s.repo.EndPartnershipWithUser(trainerId, userId)
}

func (s *UserService) AcceptRequest(trainerId, requestId int64) (int64, error) {
	return s.repo.AcceptRequest(trainerId, requestId)
}

func (s *UserService) DenyRequest(trainerId, requestId int64) error {
	return s.repo.DenyRequest(trainerId, requestId)
}

func (s *UserService) GetTrainerWorkouts(trainerId int64) ([]*entity.Workout, error) {
	return s.repo.GetTrainerWorkouts(trainerId)
}

func (s *UserService) CreateWorkoutAsTrainer(workout *entity.Workout) (int64, error) {
	return s.repo.CreateWorkoutAsTrainer(workout)
}

func (s *UserService) GetTrainerWorkoutsWithUser(trainerId, userId int64) ([]*entity.Workout, error) {
	return s.repo.GetTrainerWorkoutsWithUser(trainerId, userId)
}

func (s *UserService) GetPasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(s.hashSalt))

	return fmt.Sprintf("%x", sha1.Sum([]byte(password)))
}
