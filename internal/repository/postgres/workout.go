package postgres

import "github.com/jmoiron/sqlx"

type WorkoutRepository struct {
	db *sqlx.DB
}

func NewWorkoutRepository(db *sqlx.DB) *WorkoutRepository {
	return &WorkoutRepository{db: db}
}
