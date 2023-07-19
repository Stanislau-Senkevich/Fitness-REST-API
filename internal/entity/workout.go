package entity

import "time"

type UserWorkout struct {
	Id          int64     `db:"id" json:"id"`
	Title       string    `db:"title" json:"title" binding:"required"`
	UserId      int64     `db:"user_id" json:"user_id"`
	TrainerId   int64     `db:"trainer_id" json:"trainer_id"`
	Description string    `db:"description" json:"description"`
	Date        time.Time `db:"date" json:"date"`
}

type TrainerWorkout struct {
	Id          int64     `db:"id" json:"id"`
	Title       string    `db:"title" json:"title" binding:"required"`
	UserId      int64     `db:"user_id" json:"user_id" binding:"required"`
	TrainerId   int64     `db:"trainer_id" json:"trainer_id"`
	Description string    `db:"description" json:"description"`
	Date        time.Time `db:"date" json:"date"`
}
