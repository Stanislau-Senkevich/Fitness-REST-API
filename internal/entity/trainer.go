package entity

type Trainer struct {
	Id           int64  `db:"id"`
	Login        string `db:"login"`
	PasswordHash string `db:"password_hash"`
	Name         string `db:"name"`
	Surname      string `db:"surname"`
	Description  string `db:"description"`
	CreatedAt    string `db:"created_at"`
}
