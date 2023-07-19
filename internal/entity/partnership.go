package entity

import "time"

type Partnership struct {
	Id        int64     `db:"id"`
	UserId    int64     `db:"user_id"`
	TrainerId int64     `db:"trainer_id"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
}
