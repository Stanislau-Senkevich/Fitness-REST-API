package entity

type Admin struct {
	Id           int64  `db:"id"`
	Login        string `db:"login"`
	PasswordHash string `db:"password_hash"`
}
