package postgres

import (
	"Fitness_REST_API/internal/entity"
	"github.com/jmoiron/sqlx"
)

type TrainerRepository struct {
	db *sqlx.DB
}

func NewTrainerRepository(db *sqlx.DB) *TrainerRepository {
	return &TrainerRepository{db: db}
}

func (r *TrainerRepository) Authorize(login, passwordHash string) error {
	var trainer entity.Trainer

	query := "SELECT * FROM trainers WHERE login = $1 AND password_hash = $2"
	err := r.db.Get(&trainer, query, login, passwordHash)
	return err
}
