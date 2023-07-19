package entity

type User struct {
	Id           int64  `db:"id"`
	Email        string `db:"email" json:"email" binding:"required"`
	PasswordHash string `db:"password_hash" json:"password" binding:"required"`
	Name         string `db:"name" json:"name" binding:"required"`
	Surname      string `db:"surname" json:"surname" binding:"required"`
	CreatedAt    string `db:"created_at"`
}
