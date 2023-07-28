package postgres

import (
	"Fitness_REST_API/internal/entity"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Authorize(email, passwordHash string, role entity.Role) (int64, error) {
	var user entity.User

	query := fmt.Sprintf("SELECT * FROM %s WHERE email = $1 AND password_hash = $2 AND role = $3", userTable)
	err := r.db.Get(&user, query, email, passwordHash, role)
	return user.Id, err
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

func (r *UserRepository) IsUser(id int64) bool {
	var user entity.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE id = $1", userTable)
	err := r.db.Get(&user, query, id)
	return err == nil
}

func (r *UserRepository) CreateUser(user *entity.User, role entity.Role) (int64, error) {
	var id int64
	if r.HasEmail(user.Email) {
		return -1, errors.New("email has already reserved")
	}

	query := fmt.Sprintf("INSERT INTO %s (email, password_hash, role, name, surname)"+
		" values ($1, $2, '%s', $3, $4) RETURNING id",
		userTable, role)
	row := r.db.QueryRow(query, user.Email, user.PasswordHash, user.Name, user.Surname)

	logrus.Debugf("creating user query: %s\nargs: %s, %s, %s, %s",
		query, user.Email, user.PasswordHash, user.Name, user.Surname)

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

func (r *UserRepository) GetUserInfoById(id int64) (*entity.User, error) {
	var user entity.User
	query := fmt.Sprintf("SELECT id, email, password_hash, name, surname, role, created_at "+
		"FROM %s WHERE id = $1", userTable)
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

	addQuery := fmt.Sprintf("INSERT INTO %s (title, user_id, trainer_id, description, date) "+
		"values ($1, $2, $3, $4, $5) RETURNING id", workoutsTable)
	row := tx.QueryRow(addQuery, workout.Title, workout.UserId, workout.TrainerId, workout.Description, workout.Date)
	if err = row.Scan(&id); err != nil {
		_ = tx.Rollback()
		return -1, err
	}
	return id, tx.Commit()
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

func (r *UserRepository) GetUserWorkouts(userid int64) ([]*entity.Workout, error) {
	workouts := make([]*entity.Workout, 0)
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1 ORDER BY date DESC", workoutsTable)
	err := r.db.Select(&workouts, query, userid)
	if err != nil {
		return nil, err
	}
	return workouts, nil
}

func (r *UserRepository) GetWorkoutById(workoutId, userId int64) (*entity.Workout, error) {
	err := r.CheckAccessToWorkout(workoutId, userId)
	if err != nil {
		return nil, err
	}

	var workout entity.Workout
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", workoutsTable)
	err = r.db.Get(&workout, query, workoutId)
	if err != nil {
		return nil, err
	}
	return &workout, nil
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
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
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

func (r *UserRepository) GetTrainers() ([]*entity.User, error) {
	trainers := make([]*entity.User, 0)
	query := fmt.Sprintf("SELECT id, email, name, surname FROM %s WHERE role = $1 ORDER BY surname", userTable)
	err := r.db.Select(&trainers, query, entity.TrainerRole)
	if err != nil {
		return nil, err
	}
	return trainers, nil
}

func (r *UserRepository) GetTrainerById(id int64) (*entity.User, error) {
	var trainer entity.User
	query := fmt.Sprintf("SELECT id, email, name, surname FROM %s WHERE role = 'trainer' AND id = $1", userTable)
	err := r.db.Get(&trainer, query, id)
	if err != nil {
		return nil, err
	}
	return &trainer, nil
}

func (r *UserRepository) GetUserPartnerships(userId int64) ([]*entity.Partnership, error) {
	partnerships := make([]*entity.Partnership, 0)
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1 ORDER BY created_at DESC", partnershipsTable)
	err := r.db.Select(&partnerships, query, userId)
	if err != nil {
		return nil, err
	}
	return partnerships, nil
}

func (r *UserRepository) GetPartnership(trainerId, userId int64) (*entity.Partnership, error) {
	var p entity.Partnership
	query := fmt.Sprintf("SELECT * FROM %s WHERE trainer_id = $1 AND user_id = $2", partnershipsTable)
	err := r.db.Get(&p, query, trainerId, userId)
	return &p, err
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
		query := fmt.Sprintf("INSERT INTO %s (trainer_id, user_id, status) values "+
			"($1, $2, %s) RETURNING id", partnershipsTable, status)
		row := r.db.QueryRow(query, trainerId, userId)
		err := row.Scan(&id)
		return id, err
	}

	if p.Status == entity.StatusEndedByTrainer || p.Status == entity.StatusEndedByUser {
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
	p, err := r.GetPartnership(trainerId, userId)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return -1, err
		} else {
			return 0, err
		}
	}
	if !hasApprovedPartnership(p) {
		return -1, errors.New("no approved partnership to end")
	}

	status := "'" + entity.StatusEndedByUser + "'"
	query := fmt.Sprintf("UPDATE %s SET status = %s, ended_at = NOW() WHERE id = $1", partnershipsTable, status)
	_, err = r.db.Exec(query, p.Id)
	if err != nil {
		return 0, err
	}
	return p.Id, nil
}

func (r *UserRepository) GetTrainerPartnerships(userId int64) ([]*entity.Partnership, error) {
	partnerships := make([]*entity.Partnership, 0)
	query := fmt.Sprintf("SELECT * FROM %s WHERE trainer_id = $1 ORDER BY created_at DESC", partnershipsTable)
	err := r.db.Select(&partnerships, query, userId)
	if err != nil {
		return nil, err
	}
	return partnerships, nil
}

func (r *UserRepository) GetTrainerUsers(trainerId int64) ([]*entity.User, error) {
	if !r.IsTrainer(trainerId) {
		return nil, errors.New("not a trainer was provided")
	}

	users := make([]*entity.User, 0)

	query := fmt.Sprintf("SELECT %s.id, email, name, surname, %s.created_at "+
		"FROM %s "+
		"JOIN %s "+
		"ON %s.id = %s.user_id "+
		"WHERE %s.trainer_id =$1 "+
		"AND status = %s "+
		"ORDER BY surname;",
		userTable, partnershipsTable,
		userTable, partnershipsTable, userTable,
		partnershipsTable, partnershipsTable,
		"'"+entity.StatusApproved+"'")
	err := r.db.Select(&users, query, trainerId)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) GetTrainerRequests(trainerId int64) ([]*entity.Request, error) {
	if !r.IsTrainer(trainerId) {
		return nil, errors.New("not a trainer was provided")
	}

	requests := make([]*entity.Request, 0)

	query := fmt.Sprintf("SELECT %s.id AS user_id, %s.id AS request_id, email, name, surname, %s.created_at AS send_at "+
		"FROM %s "+
		"JOIN %s "+
		"ON %s.id = %s.user_id "+
		"WHERE %s.trainer_id =$1 "+
		"AND status = %s "+
		"ORDER BY send_at DESC;",
		userTable, partnershipsTable, partnershipsTable,
		userTable, partnershipsTable, userTable,
		partnershipsTable, partnershipsTable,
		"'"+entity.StatusRequest+"'")
	err := r.db.Select(&requests, query, trainerId)
	if err != nil {
		return nil, err
	}
	return requests, nil
}

func (r *UserRepository) GetTrainerUserById(trainerId, userId int64) (*entity.User, error) {
	if !r.IsTrainer(trainerId) {
		return nil, errors.New("not a trainer was provided")
	}

	p, _ := r.GetPartnership(trainerId, userId)
	if !hasApprovedPartnership(p) {
		return nil, errors.New("approved partnership with user was not found")
	}

	var user entity.User
	query := fmt.Sprintf("SELECT id, email, name, surname FROM %s WHERE id = $1", userTable)
	err := r.db.Get(&user, query, userId)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetTrainerRequestById(trainerId, requestId int64) (*entity.Request, error) {
	var p entity.Partnership
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", partnershipsTable)
	err := r.db.Get(&p, query, requestId)
	if err != nil {
		return nil, err
	}
	if !hasRequestOnPartnership(&p) {
		return nil, errors.New("no request for you on that id")
	}
	if p.TrainerId != trainerId {
		return nil, errors.New("no access to request")
	}

	var req entity.Request
	query = fmt.Sprintf("SELECT %s.id AS user_id, %s.id AS request_id, email, name, surname, %s.created_at AS send_at "+
		"FROM %s "+
		"JOIN %s "+
		"ON %s.id = %s.user_id "+
		"WHERE %s.id = $1;",
		userTable, partnershipsTable, partnershipsTable,
		userTable, partnershipsTable, userTable,
		partnershipsTable, partnershipsTable)
	err = r.db.Get(&req, query, requestId)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *UserRepository) InitPartnershipWithUser(trainerId, userId int64) (int64, error) {
	if !r.IsUser(userId) {
		return -1, errors.New("invalid userId")
	}

	var id int64
	p, err := r.GetPartnership(trainerId, userId)
	if err != nil {
		query := fmt.Sprintf("INSERT INTO %s (trainer_id, user_id, status) values ($1, $2, %s) RETURNING id",
			partnershipsTable, "'"+entity.StatusApproved+"'")
		row := r.db.QueryRow(query, trainerId, userId)
		if err := row.Scan(&id); err != nil {
			return 0, err
		}
		return id, nil
	}
	switch p.Status {
	case entity.StatusEndedByUser:
		return p.Id, errors.New("partnership was ended by user, it can be resumed only by request from user")
	case entity.StatusApproved:
		return p.Id, nil
	case entity.StatusEndedByTrainer, entity.StatusRequest:
		query := fmt.Sprintf("UPDATE %s SET status = %s, ended_at = null WHERE id = $1",
			partnershipsTable, "'"+entity.StatusApproved+"'")
		_, err := r.db.Exec(query, p.Id)
		if err != nil {
			return 0, err
		}
		return p.Id, nil
	}
	return 0, errors.New("undefined status of partnership")
}

func (r *UserRepository) EndPartnershipWithUser(trainerId, userId int64) (int64, error) {
	p, err := r.GetPartnership(trainerId, userId)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return -1, err
		} else {
			return 0, err
		}
	}
	if !hasApprovedPartnership(p) {
		return -1, errors.New("no approved partnership to end")
	}

	status := "'" + entity.StatusEndedByTrainer + "'"
	query := fmt.Sprintf("UPDATE %s SET status = %s, ended_at = NOW() WHERE id = $1", partnershipsTable, status)
	_, err = r.db.Exec(query, p.Id)
	if err != nil {
		return 0, err
	}
	return p.Id, nil
}

func (r *UserRepository) AcceptRequest(trainerId, requestId int64) (int64, error) {
	var id int64
	query := fmt.Sprintf("UPDATE %s SET status = %s WHERE status = %s AND trainer_id = $1 AND id = $2 RETURNING id",
		partnershipsTable, "'"+entity.StatusApproved+"'", "'"+entity.StatusRequest+"'")
	row := r.db.QueryRow(query, trainerId, requestId)
	if err := row.Scan(&id); err != nil {
		return -1, errors.New("no request to accept")
	}
	return id, nil
}

func (r *UserRepository) DenyRequest(trainerId, requestId int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE trainer_id = $1 AND id = $2 AND status = %s",
		partnershipsTable, "'"+entity.StatusRequest+"'")
	res, err := r.db.Exec(query, trainerId, requestId)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows != 1 {
		return errors.New("no request to deny")
	}
	return nil
}

func (r *UserRepository) CreateWorkoutAsTrainer(workout *entity.Workout) (int64, error) {
	p, err := r.GetPartnership(workout.TrainerId.Int64, workout.UserId)
	if err != nil || !hasApprovedPartnership(p) {
		return -1, errors.New("no rights to create workout with this user")
	}

	var id int64
	query := fmt.Sprintf("INSERT INTO %s (title, trainer_id, user_id, description, date) values "+
		"($1, $2, $3 ,$4, $5) RETURNING id", workoutsTable)
	row := r.db.QueryRow(query, workout.Title, workout.TrainerId, workout.UserId, workout.Description, workout.Date)
	if err = row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UserRepository) GetTrainerWorkouts(trainerId int64) ([]*entity.Workout, error) {
	workouts := make([]*entity.Workout, 0)
	query := fmt.Sprintf("SELECT * FROM %s WHERE trainer_id = $1 ORDER BY date DESC", workoutsTable)
	err := r.db.Select(&workouts, query, trainerId)
	if err != nil {
		return nil, err
	}
	return workouts, nil
}

func (r *UserRepository) GetTrainerWorkoutsWithUser(trainerId, userId int64) ([]*entity.Workout, error) {
	workouts := make([]*entity.Workout, 0)
	query := fmt.Sprintf("SELECT * FROM %s WHERE trainer_id = $1 AND user_id = $2 ORDER BY date DESC",
		workoutsTable)
	err := r.db.Select(&workouts, query, trainerId, userId)
	if err != nil {
		return nil, err
	}
	return workouts, nil
}

func (r *UserRepository) GetUsersId(role entity.Role) ([]int64, error) {
	idSlice := make([]int64, 0)
	query := fmt.Sprintf("SELECT id FROM %s WHERE role = '%s'", userTable, role)
	err := r.db.Select(&idSlice, query)
	if err != nil {
		return nil, err
	}
	return idSlice, nil
}

func (r *UserRepository) GetUserFullInfoById(userId int64) (*entity.UserInfo, error) {
	user, err := r.GetUserInfoById(userId)
	if err != nil {
		return nil, err
	}

	var userInfo entity.UserInfo
	userInfo.Id = user.Id
	userInfo.Email = user.Email
	userInfo.Role = user.Role
	userInfo.Name = user.Name
	userInfo.Surname = user.Surname
	userInfo.CreatedAt = user.CreatedAt

	switch user.Role {
	case entity.UserRole:
		partnerships, err := r.GetUserPartnerships(userId)
		if err != nil {
			return nil, err
		}
		userInfo.Partnerships = partnerships
		workouts, err := r.GetUserWorkouts(userId)
		if err != nil {
			return nil, err
		}
		userInfo.Workouts = workouts
	case entity.TrainerRole:
		partnerships, err := r.GetTrainerPartnerships(userId)
		if err != nil {
			return nil, err
		}
		userInfo.Partnerships = partnerships
		workouts, err := r.GetTrainerWorkouts(userId)
		if err != nil {
			return nil, err
		}
		userInfo.Workouts = workouts
	default:
		return nil, errors.New("undefined user role")
	}
	return &userInfo, nil
}

func (r *UserRepository) UpdateUser(userId int64, update *entity.UserUpdate) error {
	user, err := r.GetUserInfoById(userId)
	if err != nil {
		return errors.New("invalid userId")
	}
	if user.Email != update.Email && r.HasEmail(update.Email) {
		return errors.New("provided email has already been reserved")
	}

	query := fmt.Sprintf("UPDATE %s SET email = $1, password_hash = $2, role = $3, "+
		"name = $4, surname = $5 WHERE id = $6",
		userTable)
	_, err = r.db.Exec(query, update.Email, update.Password, update.Role, update.Name, update.Surname, userId)
	return err
}

func (r *UserRepository) DeleteUser(userId int64) error {
	if !r.IsUser(userId) {
		return errors.New("no user to delete")
	}
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1", partnershipsTable)
	_, err = tx.Exec(query, userId)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	query = fmt.Sprintf("DELETE FROM %s WHERE user_id = $1", workoutsTable)
	_, err = tx.Exec(query, userId)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	query = fmt.Sprintf("DELETE FROM %s WHERE id = $1", userTable)
	_, err = tx.Exec(query, userId)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func hasApprovedPartnership(p *entity.Partnership) bool {
	if p == nil || p.Status != entity.StatusApproved {
		return false
	}
	return true
}

func hasRequestOnPartnership(p *entity.Partnership) bool {
	if p == nil || p.Status != entity.StatusRequest {
		return false
	}
	return true
}
