package postgres

import (
	"Fitness_REST_API/internal/entity"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Authorize(login, passwordHash string) error {
	var user entity.User

	query := "SELECT * FROM users WHERE login = $1 AND password_hash = $2"
	err := r.db.Get(&user, query, login, passwordHash)
	return err
}
