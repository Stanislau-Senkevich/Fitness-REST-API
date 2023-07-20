package entity

import (
	"database/sql"
	"time"
)

const (
	StatusApproved = "approved"
	StatusRequest  = "request"
	StatusEnded    = "ended"
)

type Partnership struct {
	Id        int64        `db:"id"`
	UserId    int64        `db:"user_id"`
	TrainerId int64        `db:"trainer_id"`
	Status    string       `db:"status"`
	CreatedAt time.Time    `db:"created_at"`
	EndedAt   sql.NullTime `db:"ended_at"`
}
