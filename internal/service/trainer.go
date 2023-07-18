package service

import (
	"Fitness_REST_API/internal/repository"
	"crypto/sha1"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type TrainerService struct {
	repos      repository.Trainer
	hashSalt   string
	signingKey []byte
}

func NewTrainerService(repos repository.Trainer, hashSalt string, signingKey string) *TrainerService {
	return &TrainerService{repos: repos, hashSalt: hashSalt, signingKey: []byte(signingKey)}
}

func (s *TrainerService) SignIn(login, password string) (string, error) {
	err := s.repos.Authorize(login, s.getPasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(tokenTTL).Unix(),
		IssuedAt:  time.Now().Unix(),
	})

	return token.SignedString(s.signingKey)
}

func (s *TrainerService) ParseToken(token string) error {
	t, err := jwt.Parse(token, func(token *jwt.Token) (i interface{}, err error) {
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

func (s *TrainerService) getPasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(s.hashSalt))

	return fmt.Sprintf("%x", sha1.Sum([]byte(password)))
}