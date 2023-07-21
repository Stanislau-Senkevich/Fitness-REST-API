package postgres

import (
	"Fitness_REST_API/internal/entity"
	"database/sql"
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

func (r *UserRepository) Authorize(email, passwordHash, role string) (int64, error) {
	var user entity.User

	query := fmt.Sprintf("SELECT * FROM %s WHERE email = $1 AND password_hash = $2 AND role = $3", userTable)
	err := r.db.Get(&user, query, email, passwordHash, role)
	return user.Id, err
}

func (r *UserRepository) CreateUser(user *entity.User) (int64, error) {
	var id int64
	if r.HasEmail(user.Email) {
		return -1, errors.New("email has already reserved")
	}

	query := fmt.Sprintf("INSERT INTO %s (email, password_hash, name, surname) values ($1, $2, $4, $5) RETURNING id", userTable)
	row := r.db.QueryRow(query, user.Email, user.PasswordHash, user.Name, user.Surname)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UserRepository) HasEmail(email string) bool {
	var user entity.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE email = $1", userTable)
	_ = r.db.Get(&user, query, email)
	return user.Id > 0
}

func (r *UserRepository) GetUser(id int64) (*entity.User, error) {
	var user entity.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", userTable)
	err := r.db.Get(&user, query, id)
	return &user, err
}

func (r *UserRepository) CreateWorkoutAsUser(workout *entity.Workout) (int64, error) {

	if workout.TrainerId.Int64 > 0 && !r.IsTrainer(workout.TrainerId.Int64) {
		return -1, errors.New("can't set common user as a trainer")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return -1, err
	}

	var id int64

	addQuery := fmt.Sprintf("INSERT INTO %s (title, user_id, trainer_id, description, date) values ($1, $2, $3, $4, $5) RETURNING id", workoutsTable)
	row := tx.QueryRow(addQuery, workout.Title, workout.UserId, workout.TrainerId, workout.Description, workout.Date)
	if err = row.Scan(&id); err != nil {
		_ = tx.Rollback()
		return -1, err
	}
	return id, tx.Commit()
}

func (r *UserRepository) GetPartnership(trainerId, userId int64) (*entity.Partnership, error) {
	var p entity.Partnership
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1 AND trainer_id = $2", partnershipsTable)
	err := r.db.Get(&p, query, userId, trainerId)
	return &p, err
}

func (r *UserRepository) HasApprovedPartnership(trainerId, userId int64) bool {
	p, err := r.GetPartnership(trainerId, userId)
	if err != nil {
		return false
	}
	if p.Status != "approved" {
		return false
	}
	return true
}

func (r *UserRepository) CheckAccessToWorkout(workoutId, userId int64) error {
	var inputId struct {
		user    int64         `db:"user_id"`
		trainer sql.NullInt64 `db:"trainer_id"`
	}
	query := fmt.Sprintf("SELECT user_id, trainer_id FROM %s WHERE id = $1", workoutsTable)
	row := r.db.QueryRow(query, workoutId)
	if err := row.Scan(&inputId.user, &inputId.trainer); err != nil {
		return err
	}
	if inputId.user != userId && (!inputId.trainer.Valid || inputId.trainer.Int64 != userId) {
		return errors.New("no access to this workout")
	}
	return nil
}

func (r *UserRepository) IsTrainer(userId int64) bool {
	var user entity.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", userTable)
	err := r.db.Get(&user, query, userId)
	if err != nil {
		return false
	}
	return user.Role == entity.TrainerRole
}

func (r *UserRepository) UpdateWorkout(workoutId, userId int64, update *entity.UpdateWorkout) error {

	err := r.CheckAccessToWorkout(workoutId, userId)
	if err != nil {
		return err
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	querySample := "UPDATE %s SET %s = $1, %s = $2, %s = $3 WHERE id = $4"
	query := fmt.Sprintf(querySample, workoutsTable, "title", "description", "date")
	_, err = tx.Exec(query, update.Title, update.Description, update.Date, workoutId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (r *UserRepository) GetAllUserWorkouts(id int64) ([]*entity.Workout, error) {
	var workouts []*entity.Workout
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1", workoutsTable)
	err := r.db.Select(&workouts, query, id)
	return workouts, err
}

func (r *UserRepository) GetWorkoutById(workoutId, userId int64) (*entity.Workout, error) {
	err := r.CheckAccessToWorkout(workoutId, userId)
	if err != nil {
		return nil, err
	}

	var workout entity.Workout
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", workoutsTable)
	err = r.db.Get(&workout, query, workoutId)
	return &workout, err
}

func (r *UserRepository) DeleteWorkout(workoutId, userId int64) error {
	err := r.CheckAccessToWorkout(workoutId, userId)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", workoutsTable)
	_, err = r.db.Exec(query, workoutId)
	return err
}

func (r *UserRepository) GetAllTrainers() ([]*entity.User, error) {
	var trainers []*entity.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE role = 'trainer'", userTable)
	err := r.db.Select(&trainers, query)
	return trainers, err
}

func (r *UserRepository) GetTrainerById(id int64) (*entity.User, error) {
	var trainer entity.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE role = 'trainer' AND id = $1", userTable)
	err := r.db.Get(&trainer, query, id)
	return &trainer, err
}

func (r *UserRepository) SendRequestToTrainer(trainerId, userId int64) (int64, error) {
	if !r.IsTrainer(trainerId) {
		return -1, errors.New("can't send request not to trainer")
	}

	p, _ := r.GetPartnership(trainerId, userId)

	if p.Status == entity.StatusApproved {
		return -1, errors.New("there is already approved partnership with trainer")
	}

	if p.Status == entity.StatusRequest {
		return p.Id, nil
	}

	if p.Status == "" {
		var id int64
		status := "'" + entity.StatusRequest + "'"
		query := fmt.Sprintf("INSERT INTO %s (trainer_id, user_id, status) values ($1, $2, %s) RETURNING id", partnershipsTable, status)
		row := r.db.QueryRow(query, trainerId, userId)
		err := row.Scan(&id)
		return id, err
	}

	if p.Status == entity.StatusEnded {
		status := "'" + entity.StatusRequest + "'"
		query := fmt.Sprintf("UPDATE %s SET status = %s WHERE id = $1", partnershipsTable, status)
		_, err := r.db.Exec(query, p.Id)
		if err != nil {
			return 0, err
		}
		return p.Id, nil
	}
	return 0, errors.New("undefined partnership on provided id")
}

func (r *UserRepository) EndPartnershipWithTrainer(trainerId, userId int64) (int64, error) {
	if !r.HasApprovedPartnership(trainerId, userId) {
		return -1, errors.New("no approved partnership to end")
	}

	p, err := r.GetPartnership(trainerId, userId)
	if err != nil {
		return 0, err
	}
	status := "'" + entity.StatusEnded + "'"
	query := fmt.Sprintf("UPDATE %s SET status = %s, ended_at = NOW() WHERE id = $1", partnershipsTable, status)
	_, err = r.db.Exec(query, p.Id)
	if err != nil {
		return 0, err
	}
	return p.Id, nil
}

func (r *UserRepository) GetUserPartnerships(userId int64) ([]*entity.Partnership, error) {
	var partnerships []*entity.Partnership
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1", partnershipsTable)
	err := r.db.Select(&partnerships, query, userId)
	return partnerships, err
}
