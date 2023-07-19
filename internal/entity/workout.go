package entity

import (
	"database/sql"
	"time"
)

type Workout struct {
	Id          int64         `db:"id" json:"id"`
	Title       string        `db:"title" json:"title" binding:"required"`
	UserId      int64         `db:"user_id" json:"user_id"`
	TrainerId   sql.NullInt64 `db:"trainer_id" json:"trainer_id"`
	Description string        `db:"description" json:"description"`
	Date        time.Time     `db:"date" json:"date"`
}

type UpdateWorkout struct {
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Date        time.Time `db:"date" json:"date"`
}
