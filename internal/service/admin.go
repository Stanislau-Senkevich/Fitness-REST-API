package service

import (
	"Fitness_REST_API/internal/entity"
	"Fitness_REST_API/internal/repository"
	"crypto/sha1"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type AdminService struct {
	adminRepo  repository.Admin
	userRepo   repository.User
	hashSalt   string
	signingKey []byte
}

func NewAdminService(
	adminRepo repository.Admin,
	userRepo repository.User,
	hashSalt string,
	signingKey string) *AdminService {
	return &AdminService{adminRepo: adminRepo, userRepo: userRepo, hashSalt: hashSalt, signingKey: []byte(signingKey)}
}

func (s *AdminService) SignIn(login, password string) (string, error) {
	err := s.adminRepo.Authorize(login, s.getPasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(tokenTTL).Unix(),
		IssuedAt:  time.Now().Unix(),
	})

	return token.SignedString(s.signingKey)
}

func (s *AdminService) ParseToken(token string) error {
	t, err := jwt.Parse(token, func(token *jwt.Token) (i interface{}, err error) { //nolint
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.signingKey, nil
	})

	if err != nil {
		return err
	}

	_, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("error get user claims from token")
	}

	return nil
}

func (s *AdminService) getPasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(s.hashSalt))

	return fmt.Sprintf("%x", sha1.Sum([]byte(password)))
}

func (s *AdminService) GetUsersId(role entity.Role) ([]int64, error) {
	return s.userRepo.GetUsersId(role)
}

func (s *AdminService) GetUserFullInfoById(userId int64) (*entity.UserInfo, error) {
	return s.userRepo.GetUserFullInfoById(userId)
}

func (s *AdminService) CreateUser(user *entity.User) (int64, error) {
	return s.userRepo.CreateUser(user, user.Role)
}

func (s *AdminService) UpdateUser(userId int64, update *entity.UserUpdate) error {
	return s.userRepo.UpdateUser(userId, update)
}
func (s *AdminService) DeleteUser(userId int64) error {
	return s.userRepo.DeleteUser(userId)
}
