package entity

type User struct {
	Id           int64  `db:"id"`
	Email        string `db:"email"`
	PasswordHash string `db:"password_hash"`
	Name         string `db:"name"`
	Surname      string `db:"surname"`
	CreatedAt    string `db:"created_at"`
}
