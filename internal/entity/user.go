package entity

import "time"

type Role string

const (
	UserRole    Role = "user"
	TrainerRole Role = "trainer"
)

type User struct {
	Id           int64     `db:"id" json:"id"`
	Email        string    `db:"email" json:"email" binding:"required"`
	PasswordHash string    `db:"password_hash" json:"password_hash,omitempty" binding:"required"`
	Role         Role      `db:"role" json:"role,omitempty"`
	Name         string    `db:"name" json:"name" binding:"required"`
	Surname      string    `db:"surname" json:"surname" binding:"required"`
	CreatedAt    time.Time `db:"created_at" json:"created_at,omitempty"`
}

type UserInfo struct {
	Id           int64          `db:"id" json:"id"`
	Email        string         `db:"email" json:"email" binding:"required"`
	Role         Role           `db:"role" json:"role,omitempty"`
	Name         string         `db:"name" json:"name" binding:"required"`
	Surname      string         `db:"surname" json:"surname" binding:"required"`
	CreatedAt    time.Time      `db:"created_at" json:"created_at,omitempty"`
	Partnerships []*Partnership `json:"partnerships"`
	Workouts     []*Workout     `json:"workouts"`
}

type UserUpdate struct {
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
	Role     Role   `db:"role" json:"role"`
	Name     string `db:"name" json:"name"`
	Surname  string `db:"surname" json:"surname"`
}
