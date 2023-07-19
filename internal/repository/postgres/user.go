package postgres

import (
	"Fitness_REST_API/internal/entity"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Authorize(email, passwordHash string) (int64, error) {
	var user entity.User

	query := fmt.Sprintf("SELECT * FROM %s WHERE email = $1 AND password_hash = $2", userTable)
	err := r.db.Get(&user, query, email, passwordHash)
	return user.Id, err
}

func (r *UserRepository) CreateUser(user *entity.User) (int64, error) {
	var id int64
	query := fmt.Sprintf("INSERT INTO %s (email, password_hash, name, surname) values ($1, $2, $3, $4) RETURNING id", userTable)
	row := r.db.QueryRow(query, user.Email, user.PasswordHash, user.Name, user.Surname)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UserRepository) GetUser(id int64) (*entity.User, error) {
	var user entity.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", userTable)
	err := r.db.Get(&user, query, id)
	return &user, err
}

func (r *UserRepository) CreateWorkout(workout *entity.UserWorkout) (int64, error) {

	tx, err := r.db.Begin()
	if err != nil {
		return -1, err
	}

	var id int64

	addQuery := fmt.Sprintf("INSERT INTO %s (title, user_id) values ($1, $2) RETURNING id", usersWorkoutsTable)
	row := tx.QueryRow(addQuery, workout.Title, workout.UserId)
	if err = row.Scan(&id); err != nil {
		_ = tx.Rollback()
		return -1, err
	}
	if !workout.Date.IsZero() {
		dateQuery := fmt.Sprintf("UPDATE %s SET date = $1 WHERE id = $2", usersWorkoutsTable)
		_, err = tx.Exec(dateQuery, workout.Date, id)
		if err != nil {
			tx.Rollback()
			return -1, err
		}
	}
	if workout.TrainerId != 0 {
		err = r.CheckPartnership(workout.UserId, workout.TrainerId)
		if err != nil {
			return -1, errors.New("you must have partnership with the trainer to create workout with him as a trainer")
		}
		trainerQuery := fmt.Sprintf("UPDATE %s SET trainer_id = $1 WHERE id = $2", usersWorkoutsTable)
		_, err = tx.Exec(trainerQuery, workout.TrainerId, id)
		if err != nil {
			tx.Rollback()
			return -1, err
		}
	}
	if workout.Description != "" {
		descQuery := fmt.Sprintf("UPDATE %s SET description = $1 WHERE id = $2", usersWorkoutsTable)
		_, err = tx.Exec(descQuery, workout.Description, id)
		if err != nil {
			tx.Rollback()
			return -1, err
		}
	}

	return id, tx.Commit()
}

func (r *UserRepository) CheckPartnership(userId, trainerId int64) error {
	var p entity.Partnership
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1 AND trainer_id = $2", usersPartnershipsTable)
	return r.db.Get(&p, query, userId, trainerId)
}
