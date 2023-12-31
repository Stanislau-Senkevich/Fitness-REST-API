package postgres

import (
	"Fitness_REST_API/internal/entity"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type AdminRepository struct {
	db *sqlx.DB
}

func NewAdminRepository(db *sqlx.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) Authorize(login, passwordHash string) error {
	var admin entity.Admin

	query := fmt.Sprintf("SELECT * FROM %s WHERE login =$1 AND password_hash = $2", adminTable)
	err := r.db.Get(&admin, query, login, passwordHash)
	return err
}
