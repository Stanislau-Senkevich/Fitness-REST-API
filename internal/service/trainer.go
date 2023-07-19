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
	id, err := s.repos.Authorize(login, s.getPasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		id,
	})

	return token.SignedString(s.signingKey)
}

func (s *TrainerService) ParseToken(token string) (int64, error) {
	t, err := jwt.ParseWithClaims(token, &tokenClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.signingKey, nil
	})

	if err != nil {
		return -1, err
	}

	claims, ok := t.Claims.(*tokenClaims)
	if !ok {
		return -1, fmt.Errorf("error get user claims from token")
	}

	return claims.ID, nil
}

func (s *TrainerService) getPasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(s.hashSalt))

	return fmt.Sprintf("%x", sha1.Sum([]byte(password)))
}
