package entity

import "time"

const (
	UserRole    = "user"
	TrainerRole = "trainer"
)

type User struct {
	Id           int64     `db:"id" json:"id"`
	Email        string    `db:"email" json:"email" binding:"required"`
	PasswordHash string    `db:"password_hash" json:"password,omitempty" binding:"required"`
	Role         string    `db:"role" json:"role,omitempty"`
	Name         string    `db:"name" json:"name" binding:"required"`
	Surname      string    `db:"surname" json:"surname" binding:"required"`
	CreatedAt    time.Time `db:"created_at" json:"created_at,omitempty"`
}
