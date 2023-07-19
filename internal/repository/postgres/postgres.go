package postgres

import (
	"Fitness_REST_API/internal/config"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	adminTable                = "admins"
	userTable                 = "users"
	trainerTable              = "trainers"
	usersWorkoutsTable        = "users_workouts"
	trainersWorkoutsTable     = "trainers_workouts"
	usersPartnershipsTable    = "users_partnerships"
	trainersPartnerShipsTable = "trainers_partnerships"
)

func InitPostgresDB(cfg *config.Config) (*sqlx.DB, error) {
	conn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sqlx.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
