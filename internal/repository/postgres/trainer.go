package postgres

import (
	"Fitness_REST_API/internal/entity"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type TrainerRepository struct {
	db *sqlx.DB
}

func NewTrainerRepository(db *sqlx.DB) *TrainerRepository {
	return &TrainerRepository{db: db}
}

func (r *TrainerRepository) Authorize(login, passwordHash string) (int64, error) {
	var trainer entity.Trainer

	query := fmt.Sprintf("SELECT * FROM %s WHERE login = $1 AND password_hash = $2", trainerTable)
	err := r.db.Get(&trainer, query, login, passwordHash)
	return trainer.Id, err
}
