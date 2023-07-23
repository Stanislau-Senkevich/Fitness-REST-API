package entity

import (
	"database/sql"
	"time"
)

type Status string

const (
	StatusApproved       Status = "approved"
	StatusRequest        Status = "request"
	StatusEndedByUser    Status = "ended by user"
	StatusEndedByTrainer Status = "ended by trainer"
)

type Partnership struct {
	Id        int64        `db:"id"`
	UserId    int64        `db:"user_id"`
	TrainerId int64        `db:"trainer_id"`
	Status    Status       `db:"status"`
	CreatedAt time.Time    `db:"created_at"`
	EndedAt   sql.NullTime `db:"ended_at"`
}

type Request struct {
	RequestId int64     `db:"request_id" json:"request_id"`
	UserId    int64     `db:"user_id" json:"user_id"`
	Email     string    `db:"email" json:"email"`
	Name      string    `db:"name" json:"name"`
	Surname   string    `db:"surname" json:"surname"`
	SendAt    time.Time `db:"send_at" json:"send_at"`
}
