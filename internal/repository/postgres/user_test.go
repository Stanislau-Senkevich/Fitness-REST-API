package postgres

import (
	"Fitness_REST_API/internal/entity"
	"database/sql"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
	"time"
)

func TestUserRepository_Authorize(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type mockBehavior func(user entity.User)

	type args struct {
		email    string
		password string
	}

	r := NewUserRepository(db)

	table := []struct {
		name         string
		user         entity.User
		args         args
		mockBehavior mockBehavior
		shouldFail   bool
		shouldReturn int64
	}{
		{
			name: "Ok",
			user: entity.User{Id: 1, Email: "test", PasswordHash: "test"},
			args: args{email: "test", password: "test"},
			mockBehavior: func(user entity.User) {
				rows := sqlmock.NewRows([]string{"id", "email", "password_hash"}).AddRow(user.Id, user.Email, user.PasswordHash)
				mock.ExpectQuery("SELECT (.+) FROM users").WillReturnRows(rows)
			},
			shouldFail:   false,
			shouldReturn: int64(1),
		},
		{
			name: "Empty fields",
			user: entity.User{Id: 1, Email: "test", PasswordHash: "test"},
			args: args{email: "", password: ""},
			mockBehavior: func(user entity.User) {
				mock.ExpectQuery("SELECT (.+) FROM users").WillReturnError(errors.New("no rows selected"))
			},
			shouldFail: true,
		},
	}

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior(test.user)

			got, err := r.Authorize(test.args.email, test.args.password, entity.UserRole)
			if test.shouldFail {
				assert.Error(t, err)
				t.Skip("OK")
			}
			assert.NoError(t, err)
			assert.Equal(t, got, test.shouldReturn)
		})
	}

}

func TestUserRepository_IsTrainer(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type mockBehaviour func(userId int64)

	table := []struct {
		name          string
		userId        int64
		mockBehaviour mockBehaviour
		shouldReturn  bool
	}{
		{
			name:   "Ok",
			userId: 1,
			mockBehaviour: func(userId int64) {
				rows := sqlmock.NewRows([]string{"id", "role"}).AddRow(userId, entity.TrainerRole)
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(userId).WillReturnRows(rows)
			},
			shouldReturn: true,
		},
		{
			name:   "Not a trainer",
			userId: 1,
			mockBehaviour: func(userId int64) {
				rows := sqlmock.NewRows([]string{"id", "role"}).AddRow(userId, entity.UserRole)
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(userId).WillReturnRows(rows)
			},
			shouldReturn: false,
		},
		{
			name:   "No user was found",
			userId: 1,
			mockBehaviour: func(userId int64) {
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(userId).WillReturnError(errors.New("sql: no rows in result set"))
			},
			shouldReturn: false,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour(test.userId)

			r := NewUserRepository(db)
			got := r.IsTrainer(test.userId)

			assert.Equal(t, got, test.shouldReturn)
		})
	}
}

func TestUserRepository_IsUser(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type mockBehaviour func(userId int64)

	table := []struct {
		name          string
		userId        int64
		mockBehaviour mockBehaviour
		shouldReturn  bool
	}{
		{
			name:   "Ok",
			userId: 1,
			mockBehaviour: func(userId int64) {
				rows := sqlmock.NewRows([]string{"id", "role"}).AddRow(userId, entity.UserRole)
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(userId).WillReturnRows(rows)
			},
			shouldReturn: true,
		},
		{
			name:   "No user was found",
			userId: 1,
			mockBehaviour: func(userId int64) {
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(userId).WillReturnError(errors.New("sql: no rows in result set"))
			},
			shouldReturn: false,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour(test.userId)

			r := NewUserRepository(db)
			got := r.IsUser(test.userId)

			assert.Equal(t, got, test.shouldReturn)
		})
	}
}

func TestUserRepository_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	type mockBehaviour func()

	table := []struct {
		name          string
		inputUser     entity.User
		mockBehaviour mockBehaviour
		shouldFail    bool
		shouldReturn  int64
	}{
		{
			name:      "Ok",
			inputUser: entity.User{Email: "testEmail", PasswordHash: "testPassword", Name: "testName", Surname: "testSurname"},
			mockBehaviour: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(int64(1))
				mock.ExpectQuery("INSERT INTO users").WithArgs(
					"testEmail", "testPassword", "testName", "testSurname").WillReturnRows(rows)
			},
			shouldFail:   false,
			shouldReturn: int64(1),
		},
		{
			name:      "Email is reserved",
			inputUser: entity.User{Email: "testEmail", PasswordHash: "testPassword", Name: "testName", Surname: "testSurname"},
			mockBehaviour: func() {
				rowsSelect := sqlmock.NewRows([]string{"id", "email", "password_hash", "name", "surname"}).AddRow(int64(1), "testEmail", "testPassword", "testName", "testSurname")
				mock.ExpectQuery("SELECT (.+) FROM users").WithArgs("testEmail").WillReturnRows(rowsSelect)
			},
			shouldFail:   true,
			shouldReturn: int64(-1),
		},
		{
			name:      "Internal error",
			inputUser: entity.User{Email: "testEmail", PasswordHash: "testPassword", Name: "testName", Surname: "testSurname"},
			mockBehaviour: func() {
				mock.ExpectQuery("INSERT INTO users").WithArgs(
					"testEmail", "testPassword", "testName", "testSurname").WillReturnError(errors.New("error"))
			},
			shouldFail:   true,
			shouldReturn: int64(0),
		},
		{
			name:      "Empty fields",
			inputUser: entity.User{},
			mockBehaviour: func() {
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery("INSERT INTO users").WithArgs().WillReturnRows(rows)
			},
			shouldFail:   true,
			shouldReturn: int64(0),
		},
	}

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			r := NewUserRepository(db)
			test.mockBehaviour()
			got, err := r.CreateUser(&test.inputUser)

			if test.shouldFail {
				assert.Error(t, err)
				assert.Equal(t, got, test.shouldReturn)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, got, test.shouldReturn)
			}

		})
	}
}

func TestUserRepository_HasEmail(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type mockBehaviour func()

	table := []struct {
		name          string
		email         string
		mockBehaviour mockBehaviour
		shouldReturn  bool
	}{
		{
			name:  "No email",
			email: "test",
			mockBehaviour: func() {
				mock.ExpectQuery("SELECT (.+) FROM users").WithArgs("test").WillReturnError(errors.New("no rows selected"))
			},
			shouldReturn: false,
		},
		{
			name:  "Email found",
			email: "test",
			mockBehaviour: func() {
				row := sqlmock.NewRows([]string{"id"}).AddRow(int64(1))
				mock.ExpectQuery("SELECT (.+) FROM users").WithArgs("test").WillReturnRows(row)
			},
			shouldReturn: true,
		},
	}

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			r := NewUserRepository(db)
			test.mockBehaviour()

			got := r.HasEmail(test.email)

			assert.Equal(t, got, test.shouldReturn)
		})
	}
}

func TestUserRepository_GetUser(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type mockBehaviour func(userId int64)

	table := []struct {
		name          string
		userId        int64
		mockBehaviour mockBehaviour
		shouldFail    bool
		shouldReturn  *entity.User
	}{
		{
			name:   "Ok",
			userId: 1,
			mockBehaviour: func(userId int64) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(userId)
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(userId).WillReturnRows(rows)
			},
			shouldFail: false,
			shouldReturn: &entity.User{
				Id: 1,
			},
		},
		{
			name:   "No user was found",
			userId: 1,
			mockBehaviour: func(userId int64) {
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(userId).WillReturnError(errors.New("sql: no rows in result set"))
			},
			shouldFail:   true,
			shouldReturn: &entity.User{},
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			r := NewUserRepository(db)
			test.mockBehaviour(test.userId)

			got, err := r.GetUserInfoById(test.userId)
			if test.shouldFail {
				assert.Error(t, err)
				assert.Equal(t, got, test.shouldReturn)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, got, test.shouldReturn)
			}
		})
	}
}

func TestUserRepository_CreateWorkoutAsUser(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type mockBehaviour func(workout *entity.Workout)

	table := []struct {
		name          string
		workout       entity.Workout
		mockBehaviour mockBehaviour
		shouldFail    bool
		shouldReturn  int64
	}{
		{
			name: "Ok",
			workout: entity.Workout{
				Id:          1,
				Title:       "test",
				Description: "test",
				UserId:      1,
				TrainerId:   sql.NullInt64{Int64: 2, Valid: true},
				Date:        time.Now(),
			},
			mockBehaviour: func(workout *entity.Workout) {
				row := sqlmock.NewRows([]string{"id"}).AddRow(workout.Id)
				rowTrainer := sqlmock.NewRows([]string{"id", "role"}).AddRow(workout.TrainerId, entity.TrainerRole)
				mock.ExpectQuery("SELECT (.+) FROM users").WillReturnRows(rowTrainer)
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO workouts").
					WithArgs(workout.Title, workout.UserId, workout.TrainerId, workout.Description, workout.Date).
					WillReturnRows(row)
				mock.ExpectCommit()
			},
			shouldFail:   false,
			shouldReturn: int64(1),
		},
		{
			name: "Empty title",
			workout: entity.Workout{
				Id:          1,
				Title:       "",
				Description: "test",
				UserId:      1,
				TrainerId:   sql.NullInt64{Int64: 2, Valid: true},
				Date:        time.Now(),
			},
			mockBehaviour: func(workout *entity.Workout) {
				rowTrainer := sqlmock.NewRows([]string{"id", "role"}).AddRow(workout.TrainerId, entity.TrainerRole)
				mock.ExpectQuery("SELECT (.+) FROM users").WillReturnRows(rowTrainer)
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO workouts").
					WithArgs(workout.Title, workout.UserId, workout.TrainerId, workout.Description, workout.Date).
					WillReturnError(errors.New("title must be not null"))
				mock.ExpectRollback()
			},
			shouldFail:   true,
			shouldReturn: int64(-1),
		},
		{
			name: "Not a trainer",
			workout: entity.Workout{
				Id:          1,
				Title:       "",
				Description: "test",
				UserId:      1,
				TrainerId:   sql.NullInt64{Int64: 2, Valid: true},
				Date:        time.Now(),
			},
			mockBehaviour: func(workout *entity.Workout) {
				rowTrainer := sqlmock.NewRows([]string{"id", "role"}).AddRow(workout.TrainerId, entity.UserRole)
				mock.ExpectQuery("SELECT (.+) FROM users").WillReturnRows(rowTrainer)
			},
			shouldFail:   true,
			shouldReturn: int64(-1),
		},
		{
			name: "Invalid userId as trainerId in workout",
			workout: entity.Workout{
				Id:          1,
				Title:       "",
				Description: "test",
				UserId:      1,
				TrainerId:   sql.NullInt64{Int64: 2, Valid: true},
				Date:        time.Now(),
			},
			mockBehaviour: func(workout *entity.Workout) {
				mock.ExpectQuery("SELECT (.+) FROM users").WillReturnError(errors.New("no rows selected"))
			},
			shouldFail:   true,
			shouldReturn: int64(-1),
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			r := NewUserRepository(db)
			test.mockBehaviour(&test.workout)

			got, err := r.CreateWorkoutAsUser(&test.workout)
			if test.shouldFail {
				assert.Error(t, err)
				assert.Equal(t, got, test.shouldReturn)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, got, test.shouldReturn)
			}
		})
	}
}

func TestUserRepository_CheckAccessToWorkout(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type mockBehaviour func(workoutId int64)

	table := []struct {
		name          string
		userId        int64
		workoutId     int64
		mockBehaviour mockBehaviour
		shouldFail    bool
	}{
		{
			name:      "Ok",
			userId:    1,
			workoutId: 1,
			mockBehaviour: func(workoutId int64) {
				rows := sqlmock.NewRows([]string{"trainer_id", "user_id"}).AddRow(int64(0), int64(1))
				mock.ExpectQuery("SELECT user_id, trainer_id FROM workouts").
					WithArgs(workoutId).WillReturnRows(rows)
			},
			shouldFail: false,
		},
		{
			name:      "No workout",
			userId:    1,
			workoutId: 1,
			mockBehaviour: func(workoutId int64) {
				mock.ExpectQuery("SELECT user_id, trainer_id FROM workouts").
					WithArgs(workoutId).WillReturnError(errors.New("sql: no rows in result set"))
			},
			shouldFail: true,
		},
		{
			name:      "Access as trainer",
			userId:    1,
			workoutId: 1,
			mockBehaviour: func(workoutId int64) {
				rows := sqlmock.NewRows([]string{"trainer_id", "user_id"}).AddRow(int64(1), int64(100))
				mock.ExpectQuery("SELECT user_id, trainer_id FROM workouts").
					WithArgs(workoutId).WillReturnRows(rows)
			},
			shouldFail: false,
		},
		{
			name:      "Access as client",
			userId:    10,
			workoutId: 1,
			mockBehaviour: func(workoutId int64) {
				rows := sqlmock.NewRows([]string{"trainer_id", "user_id"}).AddRow(int64(120), int64(10))
				mock.ExpectQuery("SELECT user_id, trainer_id FROM workouts").
					WithArgs(workoutId).WillReturnRows(rows)
			},
			shouldFail: false,
		},
		{
			name:      "No access",
			userId:    1,
			workoutId: 1,
			mockBehaviour: func(workoutId int64) {
				rows := sqlmock.NewRows([]string{"trainer_id", "user_id"}).AddRow(int64(120), int64(10))
				mock.ExpectQuery("SELECT user_id, trainer_id FROM workouts").
					WithArgs(workoutId).WillReturnRows(rows)
			},
			shouldFail: true,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour(test.workoutId)

			r := NewUserRepository(db)
			err = r.CheckAccessToWorkout(test.workoutId, test.userId)
			if test.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserRepository_GetAllUserWorkouts(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type mockBehaviour func(userId int64)

	table := []struct {
		name          string
		userId        int64
		mockBehaviour mockBehaviour
		shouldFail    bool
		shouldReturn  []*entity.Workout
	}{
		{
			name:   "Ok",
			userId: 1,
			mockBehaviour: func(userId int64) {
				rows := sqlmock.NewRows([]string{"id", "trainer_id", "user_id", "title", "description", "date"}).
					AddRow(int64(1), int64(2), int64(1), "test", "test", time.Unix(5, 0))
				mock.ExpectQuery("SELECT (.+) FROM workouts").WithArgs(userId).WillReturnRows(rows)
			},
			shouldFail: false,
			shouldReturn: []*entity.Workout{{
				Id:          1,
				TrainerId:   sql.NullInt64{2, true},
				UserId:      1,
				Title:       "test",
				Description: "test",
				Date:        time.Unix(5, 0),
			}},
		},
		{
			name:   "No workouts",
			userId: 1,
			mockBehaviour: func(userId int64) {
				rows := sqlmock.NewRows([]string{"id", "trainer_id", "user_id", "title", "description", "date"})
				mock.ExpectQuery("SELECT (.+) FROM workouts").WithArgs(userId).WillReturnRows(rows)
			},
			shouldFail:   false,
			shouldReturn: []*entity.Workout{},
		},
		{
			name:   "Internal error",
			userId: 1,
			mockBehaviour: func(userId int64) {
				mock.ExpectQuery("SELECT (.+) FROM workouts").
					WithArgs(userId).WillReturnError(errors.New("internal error"))
			},
			shouldFail:   true,
			shouldReturn: nil,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour(test.userId)

			r := NewUserRepository(db)

			got, err := r.GetUserWorkouts(test.userId)
			if test.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, got, test.shouldReturn)
		})
	}
}

func TestUserRepository_GetWorkoutById(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type mockBehaviour func(workoutId, userId int64)

	table := []struct {
		name          string
		userId        int64
		workoutId     int64
		mockBehaviour mockBehaviour
		shouldFail    bool
		shouldReturn  *entity.Workout
	}{
		{
			name:   "Ok",
			userId: 1,
			mockBehaviour: func(workoutId, userId int64) {
				rows := sqlmock.NewRows([]string{"trainer_id", "user_id"}).AddRow(int64(0), int64(1))
				mock.ExpectQuery("SELECT user_id, trainer_id FROM workouts").
					WithArgs(workoutId).WillReturnRows(rows)
				rows = sqlmock.NewRows([]string{"id", "trainer_id", "user_id", "title", "description", "date"}).
					AddRow(int64(1), int64(2), int64(1), "test", "test", time.Unix(5, 0))
				mock.ExpectQuery("SELECT (.+) FROM workouts").WithArgs(workoutId).WillReturnRows(rows)
			},
			shouldFail: false,
			shouldReturn: &entity.Workout{
				Id:          1,
				TrainerId:   sql.NullInt64{2, true},
				UserId:      1,
				Title:       "test",
				Description: "test",
				Date:        time.Unix(5, 0),
			},
		},
		{
			name:   "No access",
			userId: 1,
			mockBehaviour: func(workoutId, userId int64) {
				rows := sqlmock.NewRows([]string{"trainer_id", "user_id"}).AddRow(int64(100), int64(10))
				mock.ExpectQuery("SELECT user_id, trainer_id FROM workouts").
					WithArgs(workoutId).WillReturnRows(rows)
			},
			shouldFail:   true,
			shouldReturn: nil,
		},
		{
			name:   "No workout",
			userId: 1,
			mockBehaviour: func(workoutId, userId int64) {
				mock.ExpectQuery("SELECT user_id, trainer_id FROM workouts").
					WithArgs(workoutId).WillReturnError(errors.New("sql: no rows in result set"))
			},
			shouldFail:   true,
			shouldReturn: nil,
		},
		{
			name:   "Internal error",
			userId: 1,
			mockBehaviour: func(workoutId, userId int64) {
				rows := sqlmock.NewRows([]string{"trainer_id", "user_id"}).AddRow(int64(0), int64(1))
				mock.ExpectQuery("SELECT user_id, trainer_id FROM workouts").
					WithArgs(workoutId).WillReturnRows(rows)
				mock.ExpectQuery("SELECT (.+) FROM workouts").
					WithArgs(workoutId).WillReturnError(errors.New("internal error"))
			},
			shouldFail:   true,
			shouldReturn: nil,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour(test.workoutId, test.userId)

			r := NewUserRepository(db)

			got, err := r.GetWorkoutById(test.workoutId, test.userId)
			if test.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, got, test.shouldReturn)
		})
	}
}

func TestUserRepository_UpdateWorkout(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type mockBehaviour func(workoutId, userId int64, update *entity.UpdateWorkout)

	table := []struct {
		name          string
		userId        int64
		workoutId     int64
		update        *entity.UpdateWorkout
		mockBehaviour mockBehaviour
		shouldFail    bool
	}{
		{
			name:      "Ok",
			userId:    1,
			workoutId: 1,
			update:    &entity.UpdateWorkout{Title: "newTitle", Description: "newDesc", Date: time.Now()},
			mockBehaviour: func(workoutId, userId int64, update *entity.UpdateWorkout) {
				rowsSelect := sqlmock.NewRows([]string{"user_id", "trainer_id"}).AddRow(int64(1), int64(0))
				mock.ExpectQuery("SELECT (.+) FROM workouts").
					WithArgs(workoutId).WillReturnRows(rowsSelect)
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE workouts SET (.+) WHERE (.+)").
					WithArgs(update.Title, update.Description, update.Date, workoutId).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			shouldFail: false,
		},
		{
			name:      "Empty fields",
			userId:    1,
			workoutId: 1,
			update:    &entity.UpdateWorkout{},
			mockBehaviour: func(workoutId, userId int64, update *entity.UpdateWorkout) {
				rowsSelect := sqlmock.NewRows([]string{"user_id", "trainer_id"}).AddRow(int64(1), int64(0))
				mock.ExpectQuery("SELECT (.+) FROM workouts").
					WithArgs(workoutId).WillReturnRows(rowsSelect)
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE workouts SET (.+) WHERE (.+)").
					WithArgs(update.Title, update.Description, update.Date, workoutId).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			shouldFail: false,
		},
		{
			name:      "Internal error",
			userId:    1,
			workoutId: 1,
			update:    &entity.UpdateWorkout{},
			mockBehaviour: func(workoutId, userId int64, update *entity.UpdateWorkout) {
				rowsSelect := sqlmock.NewRows([]string{"user_id", "trainer_id"}).AddRow(int64(1), int64(0))
				mock.ExpectQuery("SELECT (.+) FROM workouts").
					WithArgs(workoutId).WillReturnRows(rowsSelect)
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE workouts SET (.+) WHERE (.+)").
					WithArgs(update.Title, update.Description, update.Date, workoutId).WillReturnError(errors.New("internal error"))
				mock.ExpectRollback()
			},
			shouldFail: true,
		},
		{
			name:      "No workout was found",
			userId:    1,
			workoutId: 1,
			update:    &entity.UpdateWorkout{},
			mockBehaviour: func(workoutId, userId int64, update *entity.UpdateWorkout) {
				rowsSelect := sqlmock.NewRows([]string{"user_id", "trainer_id"})
				mock.ExpectQuery("SELECT (.+) FROM workouts").
					WithArgs(workoutId).WillReturnRows(rowsSelect)
			},
			shouldFail: true,
		},
		{
			name:      "No access to workout",
			userId:    1,
			workoutId: 1,
			update:    &entity.UpdateWorkout{},
			mockBehaviour: func(workoutId, userId int64, update *entity.UpdateWorkout) {
				rowsSelect := sqlmock.NewRows([]string{"user_id", "trainer_id"}).AddRow(int64(100), int64(101))
				mock.ExpectQuery("SELECT (.+) FROM workouts").
					WithArgs(workoutId).WillReturnRows(rowsSelect)
			},
			shouldFail: true,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour(test.workoutId, test.userId, test.update)

			r := NewUserRepository(db)
			err = r.UpdateWorkout(test.workoutId, test.userId, test.update)

			if test.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}

func TestUserRepository_DeleteWorkout(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type mockBehaviour func(workoutId, userId int64)

	table := []struct {
		name          string
		userId        int64
		workoutId     int64
		mockBehaviour mockBehaviour
		shouldFail    bool
	}{
		{
			name:   "Ok",
			userId: 1,
			mockBehaviour: func(workoutId, userId int64) {
				rows := sqlmock.NewRows([]string{"trainer_id", "user_id"}).AddRow(int64(0), int64(1))
				mock.ExpectQuery("SELECT user_id, trainer_id FROM workouts").
					WithArgs(workoutId).WillReturnRows(rows)
				mock.ExpectExec("DELETE FROM workouts").
					WithArgs(workoutId).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			shouldFail: false,
		},
		{
			name:   "No workout",
			userId: 1,
			mockBehaviour: func(workoutId, userId int64) {
				mock.ExpectQuery("SELECT user_id, trainer_id FROM workouts").
					WithArgs(workoutId).WillReturnError(errors.New("sql: no rows in result set"))
			},
			shouldFail: true,
		},
		{
			name:   "Internal error",
			userId: 1,
			mockBehaviour: func(workoutId, userId int64) {
				rows := sqlmock.NewRows([]string{"trainer_id", "user_id"}).AddRow(int64(0), int64(1))
				mock.ExpectQuery("SELECT user_id, trainer_id FROM workouts").
					WithArgs(workoutId).WillReturnRows(rows)
				mock.ExpectExec("DELETE FROM workouts").
					WithArgs(workoutId).WillReturnError(errors.New("internal error"))
			},
			shouldFail: true,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour(test.workoutId, test.userId)

			r := NewUserRepository(db)

			err = r.DeleteWorkout(test.workoutId, test.userId)
			if test.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserRepository_GetAllTrainers(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type mockBehaviour func()

	table := []struct {
		name          string
		mockBehaviour mockBehaviour
		shouldFail    bool
		shouldReturn  []*entity.User
	}{
		{
			name: "Ok",
			mockBehaviour: func() {
				rows := sqlmock.NewRows([]string{"id", "email", "name", "surname"}).
					AddRow(int64(1), "test1", "test1", "test1").
					AddRow(int64(2), "test2", "test2", "test2")
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(entity.TrainerRole).WillReturnRows(rows)
			},
			shouldFail: false,
			shouldReturn: []*entity.User{
				{Id: 1, Email: "test1", Name: "test1", Surname: "test1"},
				{Id: 2, Email: "test2", Name: "test2", Surname: "test2"},
			},
		},
		{
			name: "No trainers",
			mockBehaviour: func() {
				rows := sqlmock.NewRows([]string{"id", "email", "name", "surname"})
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(entity.TrainerRole).WillReturnRows(rows)
			},
			shouldFail:   false,
			shouldReturn: []*entity.User{},
		},
		{
			name: "Internal error",
			mockBehaviour: func() {
				mock.ExpectQuery("SELECT (.+) FROM users").WillReturnError(errors.New("internal error"))
			},
			shouldFail:   true,
			shouldReturn: nil,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour()

			r := NewUserRepository(db)
			got, err := r.GetTrainers()
			if test.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, got, test.shouldReturn)
		})
	}
}

func TestUserRepository_GetTrainerById(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type mockBehaviour func(trainerId int64)

	table := []struct {
		name          string
		trainerId     int64
		mockBehaviour mockBehaviour
		shouldFail    bool
		shouldReturn  *entity.User
	}{
		{
			name:      "Ok",
			trainerId: 1,
			mockBehaviour: func(trainerId int64) {
				rows := sqlmock.NewRows([]string{"id", "email", "name", "surname"}).
					AddRow(trainerId, "test1", "test1", "test1")
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(trainerId).WillReturnRows(rows)
			},
			shouldFail: false,
			shouldReturn: &entity.User{
				Id: 1, Email: "test1", Name: "test1", Surname: "test1"},
		},
		{
			name: "No trainer on id",
			mockBehaviour: func(trainerId int64) {
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(trainerId).WillReturnError(errors.New("sql: no rows in result set"))
			},
			shouldFail:   true,
			shouldReturn: nil,
		},
		{
			name: "Internal error",
			mockBehaviour: func(trainerId int64) {
				mock.ExpectQuery("SELECT (.+) FROM users").WillReturnError(errors.New("internal error"))
			},
			shouldFail:   true,
			shouldReturn: nil,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour(test.trainerId)

			r := NewUserRepository(db)
			got, err := r.GetTrainerById(test.trainerId)
			if test.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, got, test.shouldReturn)
		})
	}
}

func TestUserRepository_GetUserPartnerships(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type mockBehaviour func(userId int64)

	table := []struct {
		name          string
		userId        int64
		mockBehaviour mockBehaviour
		shouldFail    bool
		shouldReturn  []*entity.Partnership
	}{
		{
			name:   "Ok",
			userId: 1,
			mockBehaviour: func(userId int64) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "trainer_id", "status"}).
					AddRow(int64(1), int64(1), int64(2), entity.StatusApproved).
					AddRow(int64(2), int64(1), int64(3), entity.StatusRequest).
					AddRow(int64(3), int64(1), int64(4), entity.StatusEndedByUser)
				mock.ExpectQuery("SELECT (.+) FROM partnerships").
					WithArgs(userId).WillReturnRows(rows)
			},
			shouldFail: false,
			shouldReturn: []*entity.Partnership{
				{Id: 1, UserId: 1, TrainerId: 2, Status: entity.StatusApproved},
				{Id: 2, UserId: 1, TrainerId: 3, Status: entity.StatusRequest},
				{Id: 3, UserId: 1, TrainerId: 4, Status: entity.StatusEndedByUser},
			},
		},
		{
			name:   "No partnerships",
			userId: 1,
			mockBehaviour: func(userId int64) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "trainer_id", "status"})
				mock.ExpectQuery("SELECT (.+) FROM partnerships").
					WithArgs(userId).WillReturnRows(rows)
			},
			shouldFail:   false,
			shouldReturn: []*entity.Partnership{},
		},
		{
			name:   "Internal error",
			userId: 1,
			mockBehaviour: func(userId int64) {
				mock.ExpectQuery("SELECT (.+) FROM partnerships").
					WithArgs(userId).WillReturnError(errors.New("internal error"))
			},
			shouldFail:   true,
			shouldReturn: nil,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour(test.userId)

			r := NewUserRepository(db)
			got, err := r.GetUserPartnerships(test.userId)
			if test.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, got, test.shouldReturn)
		})
	}
}

func TestUserRepository_GetPartnership(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type mockBehaviour func(trainerId, userId int64)

	table := []struct {
		name          string
		userId        int64
		trainerId     int64
		mockBehaviour mockBehaviour
		shouldFail    bool
		shouldReturn  *entity.Partnership
	}{
		{
			name:      "Ok",
			userId:    1,
			trainerId: 2,
			mockBehaviour: func(trainerId, userId int64) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "trainer_id", "status"}).AddRow(int64(1), int64(1), int64(2), entity.StatusApproved)
				mock.ExpectQuery("SELECT (.+) FROM partnerships").WithArgs(trainerId, userId).WillReturnRows(rows)
			},
			shouldFail:   false,
			shouldReturn: &entity.Partnership{Id: 1, UserId: 1, TrainerId: 2, Status: entity.StatusApproved},
		},
		{
			name:      "No partnership",
			userId:    1,
			trainerId: 2,
			mockBehaviour: func(trainerId, userId int64) {
				mock.ExpectQuery("SELECT (.+) FROM partnerships").WithArgs(trainerId, userId).WillReturnError(errors.New("sql: no rows in result set"))
			},
			shouldFail:   true,
			shouldReturn: &entity.Partnership{},
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			r := NewUserRepository(db)
			test.mockBehaviour(test.trainerId, test.userId)

			got, err := r.GetPartnership(test.trainerId, test.userId)
			if test.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, got, test.shouldReturn)
		})
	}
}

func TestUserRepository_SendRequestToTrainer(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type mockBehaviour func(trainerId, userId int64)

	table := []struct {
		name          string
		userId        int64
		trainerId     int64
		mockBehaviour mockBehaviour
		shouldFail    bool
		shouldReturn  int64
	}{
		{
			name:      "Ok Insert",
			userId:    1,
			trainerId: 2,
			mockBehaviour: func(trainerId, userId int64) {
				trainerRow := sqlmock.NewRows([]string{"id", "role"}).
					AddRow(trainerId, entity.TrainerRole)
				idRow := sqlmock.NewRows([]string{"id"}).
					AddRow(int64(1))
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(trainerId).WillReturnRows(trainerRow)
				mock.ExpectQuery("SELECT (.+) FROM partnerships").
					WithArgs(trainerId, userId).WillReturnError(errors.New("sql: no rows in result set"))
				mock.ExpectQuery("INSERT INTO partnerships").
					WithArgs(trainerId, userId).WillReturnRows(idRow)
			},
			shouldFail:   false,
			shouldReturn: 1,
		},
		{
			name:      "Ok Update Ended",
			userId:    1,
			trainerId: 2,
			mockBehaviour: func(trainerId, userId int64) {
				trainerRow := sqlmock.NewRows([]string{"id", "role"}).
					AddRow(trainerId, entity.TrainerRole)
				partnershipRow := sqlmock.NewRows([]string{"id", "trainer_id", "user_id", "status"}).
					AddRow(int64(1), int64(2), int64(1), entity.StatusEndedByUser)
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(trainerId).WillReturnRows(trainerRow)
				mock.ExpectQuery("SELECT (.+) FROM partnerships").
					WithArgs(trainerId, userId).WillReturnRows(partnershipRow)
				mock.ExpectExec("UPDATE partnerships SET (.+)").
					WithArgs(userId).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			shouldFail:   false,
			shouldReturn: 1,
		},
		{
			name:      "Ok Update Request",
			userId:    1,
			trainerId: 2,
			mockBehaviour: func(trainerId, userId int64) {
				trainerRow := sqlmock.NewRows([]string{"id", "role"}).
					AddRow(trainerId, entity.TrainerRole)
				partnershipRow := sqlmock.NewRows([]string{"id", "trainer_id", "user_id", "status"}).
					AddRow(int64(1), int64(2), int64(1), entity.StatusRequest)
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(trainerId).WillReturnRows(trainerRow)
				mock.ExpectQuery("SELECT (.+) FROM partnerships").
					WithArgs(trainerId, userId).WillReturnRows(partnershipRow)
			},
			shouldFail:   false,
			shouldReturn: 1,
		},
		{
			name:      "Not a trainer",
			userId:    1,
			trainerId: 2,
			mockBehaviour: func(trainerId, userId int64) {
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(trainerId).WillReturnError(errors.New("sql: no rows in result set"))
			},
			shouldFail:   true,
			shouldReturn: -1,
		},
		{
			name:      "Already approved partnership",
			userId:    1,
			trainerId: 2,
			mockBehaviour: func(trainerId, userId int64) {
				trainerRow := sqlmock.NewRows([]string{"id", "role"}).
					AddRow(trainerId, entity.TrainerRole)
				partnershipRow := sqlmock.NewRows([]string{"id", "trainer_id", "user_id", "status"}).
					AddRow(int64(1), int64(2), int64(1), entity.StatusApproved)
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(trainerId).WillReturnRows(trainerRow)
				mock.ExpectQuery("SELECT (.+) FROM partnerships").
					WithArgs(trainerId, userId).WillReturnRows(partnershipRow)
			},
			shouldFail:   true,
			shouldReturn: -1,
		},
		{
			name:      "Undefined status",
			userId:    1,
			trainerId: 2,
			mockBehaviour: func(trainerId, userId int64) {
				trainerRow := sqlmock.NewRows([]string{"id", "role"}).
					AddRow(trainerId, entity.TrainerRole)
				partnershipRow := sqlmock.NewRows([]string{"id", "trainer_id", "user_id", "status"}).
					AddRow(int64(1), int64(2), int64(1), "?????")
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(trainerId).WillReturnRows(trainerRow)
				mock.ExpectQuery("SELECT (.+) FROM partnerships").
					WithArgs(trainerId, userId).WillReturnRows(partnershipRow)
			},
			shouldFail:   true,
			shouldReturn: 0,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour(test.trainerId, test.userId)

			r := NewUserRepository(db)
			got, err := r.SendRequestToTrainer(test.trainerId, test.userId)
			if test.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, got, test.shouldReturn)
		})
	}
}

func TestUserRepository_EndPartnershipWithTrainer(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type mockBehaviour func(trainerId, userId int64)

	table := []struct {
		name          string
		userId        int64
		trainerId     int64
		mockBehaviour mockBehaviour
		shouldFail    bool
		shouldReturn  int64
	}{
		{
			name:      "Ok",
			userId:    1,
			trainerId: 2,
			mockBehaviour: func(trainerId, userId int64) {
				partnershipRow := sqlmock.NewRows([]string{"id", "user_id", "trainer_id", "status"}).
					AddRow(int64(1), int64(1), int64(2), entity.StatusApproved)
				mock.ExpectQuery("SELECT (.+) FROM partnerships").
					WithArgs(trainerId, userId).WillReturnRows(partnershipRow)
				mock.ExpectExec("UPDATE partnerships SET (.+)").
					WithArgs(int64(1)).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			shouldFail:   false,
			shouldReturn: 1,
		},
		{
			name:      "No approved partnership",
			userId:    1,
			trainerId: 2,
			mockBehaviour: func(trainerId, userId int64) {
				partnershipRow := sqlmock.NewRows([]string{"id", "user_id", "trainer_id", "status"}).
					AddRow(int64(1), int64(1), int64(2), entity.StatusEndedByUser)
				mock.ExpectQuery("SELECT (.+) FROM partnerships").
					WithArgs(trainerId, userId).WillReturnRows(partnershipRow)
			},
			shouldFail:   true,
			shouldReturn: -1,
		},
		{
			name:      "No partnership",
			userId:    1,
			trainerId: 2,
			mockBehaviour: func(trainerId, userId int64) {
				mock.ExpectQuery("SELECT (.+) FROM partnerships").
					WithArgs(trainerId, userId).WillReturnError(errors.New("sql: no rows in result set"))
			},
			shouldFail:   true,
			shouldReturn: -1,
		},
		{
			name:      "Internal error",
			userId:    1,
			trainerId: 2,
			mockBehaviour: func(trainerId, userId int64) {
				partnershipRow := sqlmock.NewRows([]string{"id", "user_id", "trainer_id", "status"}).
					AddRow(int64(1), int64(1), int64(2), entity.StatusApproved)
				mock.ExpectQuery("SELECT (.+) FROM partnerships").
					WithArgs(trainerId, userId).WillReturnRows(partnershipRow)
				mock.ExpectExec("UPDATE partnerships SET (.+)").
					WithArgs(int64(1)).WillReturnError(errors.New("internal error"))
			},
			shouldFail:   true,
			shouldReturn: 0,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour(test.trainerId, test.userId)

			r := NewUserRepository(db)
			got, err := r.EndPartnershipWithTrainer(test.trainerId, test.userId)
			if test.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, got, test.shouldReturn)
		})
	}

}

func TestUserRepository_GetTrainerUsers(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type mockBehaviour func(trainerId int64)

	table := []struct {
		name          string
		trainerId     int64
		mockBehaviour mockBehaviour
		shouldFail    bool
		shouldReturn  []*entity.User
	}{
		{
			name:      "Ok",
			trainerId: 1,
			mockBehaviour: func(trainerId int64) {
				rowsTrainer := sqlmock.NewRows([]string{"id", "role"}).AddRow(int64(1), entity.TrainerRole)
				rowsSelect := sqlmock.NewRows([]string{"id", "email", "name", "surname"}).
					AddRow(int64(2), "test1", "test1", "test1").
					AddRow(int64(3), "test2", "test2", "test2")
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(trainerId).WillReturnRows(rowsTrainer)
				mock.ExpectQuery("SELECT (.+) FROM users").
					WillReturnRows(rowsSelect)
			},
			shouldFail: false,
			shouldReturn: []*entity.User{
				{Id: 2, Email: "test1", Name: "test1", Surname: "test1"},
				{Id: 3, Email: "test2", Name: "test2", Surname: "test2"},
			},
		},
		{
			name:      "Not a trainer",
			trainerId: 1,
			mockBehaviour: func(trainerId int64) {
				rowsTrainer := sqlmock.NewRows([]string{"id", "role"}).AddRow(int64(1), entity.UserRole)
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(trainerId).WillReturnRows(rowsTrainer)
			},
			shouldFail:   true,
			shouldReturn: nil,
		},
		{
			name:      "Empty output",
			trainerId: 1,
			mockBehaviour: func(trainerId int64) {
				rowsTrainer := sqlmock.NewRows([]string{"id", "role"}).AddRow(int64(1), entity.TrainerRole)
				rowsSelect := sqlmock.NewRows([]string{"id", "email", "name", "surname"})
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(trainerId).WillReturnRows(rowsTrainer)
				mock.ExpectQuery("SELECT (.+) FROM users").
					WillReturnRows(rowsSelect)
			},
			shouldFail:   false,
			shouldReturn: []*entity.User{},
		},
		{
			name:      "Internal error",
			trainerId: 1,
			mockBehaviour: func(trainerId int64) {
				rowsTrainer := sqlmock.NewRows([]string{"id", "role"}).AddRow(int64(1), entity.TrainerRole)
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(trainerId).WillReturnRows(rowsTrainer)
				mock.ExpectQuery("SELECT (.+) FROM users").
					WillReturnError(errors.New("internal error"))
			},
			shouldFail:   true,
			shouldReturn: nil,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			test.mockBehaviour(test.trainerId)

			r := NewUserRepository(db)
			got, err := r.GetTrainerUsers(test.trainerId)
			if test.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, got, test.shouldReturn)
		})
	}
}

func TestUserRepository_GetTrainerRequests(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type mockBehaviour func(trainerId int64)

	table := []struct {
		name          string
		trainerId     int64
		mockBehaviour mockBehaviour
		shouldFail    bool
		shouldReturn  []*entity.Request
	}{
		{
			name:      "Ok",
			trainerId: 1,
			mockBehaviour: func(trainerId int64) {
				rowsTrainer := sqlmock.NewRows([]string{"id", "role"}).AddRow(int64(1), entity.TrainerRole)
				rowsSelect := sqlmock.NewRows([]string{"request_id", "user_id", "email", "name", "surname"}).
					AddRow(int64(1), int64(2), "test1", "test1", "test1").
					AddRow(int64(2), int64(3), "test2", "test2", "test2")
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(trainerId).WillReturnRows(rowsTrainer)
				mock.ExpectQuery("SELECT (.+) FROM users").
					WillReturnRows(rowsSelect)
			},
			shouldFail: false,
			shouldReturn: []*entity.Request{
				{RequestId: 1, UserId: 2, Email: "test1", Name: "test1", Surname: "test1"},
				{RequestId: 2, UserId: 3, Email: "test2", Name: "test2", Surname: "test2"},
			},
		},
		{
			name:      "Not a trainer",
			trainerId: 1,
			mockBehaviour: func(trainerId int64) {
				rowsTrainer := sqlmock.NewRows([]string{"id", "role"}).AddRow(int64(1), entity.UserRole)
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(trainerId).WillReturnRows(rowsTrainer)
			},
			shouldFail:   true,
			shouldReturn: nil,
		},
		{
			name:      "Empty output",
			trainerId: 1,
			mockBehaviour: func(trainerId int64) {
				rowsTrainer := sqlmock.NewRows([]string{"id", "role"}).AddRow(int64(1), entity.TrainerRole)
				rowsSelect := sqlmock.NewRows([]string{"request_id", "user_id", "email", "name", "surname"})
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(trainerId).WillReturnRows(rowsTrainer)
				mock.ExpectQuery("SELECT (.+) FROM users").
					WillReturnRows(rowsSelect)
			},
			shouldFail:   false,
			shouldReturn: []*entity.Request{},
		},
		{
			name:      "Internal error",
			trainerId: 1,
			mockBehaviour: func(trainerId int64) {
				rowsTrainer := sqlmock.NewRows([]string{"id", "role"}).AddRow(int64(1), entity.TrainerRole)
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(trainerId).WillReturnRows(rowsTrainer)
				mock.ExpectQuery("SELECT (.+) FROM users").
					WillReturnError(errors.New("internal error"))
			},
			shouldFail:   true,
			shouldReturn: nil,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			test.mockBehaviour(test.trainerId)

			r := NewUserRepository(db)
			got, err := r.GetTrainerRequests(test.trainerId)
			if test.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, got, test.shouldReturn)
		})
	}
}

func TestUserRepository_GetTrainerUserById(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type mockBehaviour func(trainerId, userId int64)

	table := []struct {
		name          string
		trainerId     int64
		userId        int64
		mockBehaviour mockBehaviour
		shouldFail    bool
		shouldReturn  *entity.User
	}{
		{
			name:      "Ok",
			trainerId: 1,
			userId:    2,
			mockBehaviour: func(trainerId, userId int64) {
				rowsTrainer := sqlmock.NewRows([]string{"id", "role"}).AddRow(int64(1), entity.TrainerRole)
				rowsPartnership := sqlmock.NewRows([]string{"id", "user_id", "trainer_id", "status"}).
					AddRow(int64(1), int64(2), int64(1), entity.StatusApproved)
				rowsSelect := sqlmock.NewRows([]string{"id", "email", "name", "surname"}).
					AddRow(int64(2), "test1", "test1", "test1")
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(trainerId).WillReturnRows(rowsTrainer)
				mock.ExpectQuery("SELECT (.+) FROM partnerships").
					WithArgs(trainerId, userId).WillReturnRows(rowsPartnership)
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(userId).WillReturnRows(rowsSelect)
			},
			shouldFail:   false,
			shouldReturn: &entity.User{Id: 2, Email: "test1", Name: "test1", Surname: "test1"},
		},
		{
			name:      "No approved partnership",
			trainerId: 1,
			userId:    2,
			mockBehaviour: func(trainerId, userId int64) {
				rowsTrainer := sqlmock.NewRows([]string{"id", "role"}).AddRow(int64(1), entity.TrainerRole)
				rowsPartnership := sqlmock.NewRows([]string{"id", "user_id", "trainer_id", "status"}).
					AddRow(int64(1), int64(2), int64(1), entity.StatusRequest)
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(trainerId).WillReturnRows(rowsTrainer)
				mock.ExpectQuery("SELECT (.+) FROM partnerships").
					WithArgs(trainerId, userId).WillReturnRows(rowsPartnership)
			},
			shouldFail:   true,
			shouldReturn: nil,
		},
		{
			name:      "No partnership",
			trainerId: 1,
			userId:    2,
			mockBehaviour: func(trainerId, userId int64) {
				rowsTrainer := sqlmock.NewRows([]string{"id", "role"}).AddRow(int64(1), entity.TrainerRole)
				rowsPartnership := sqlmock.NewRows([]string{"id", "user_id", "trainer_id", "status"})
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(trainerId).WillReturnRows(rowsTrainer)
				mock.ExpectQuery("SELECT (.+) FROM partnerships").
					WithArgs(trainerId, userId).WillReturnRows(rowsPartnership)
			},
			shouldFail:   true,
			shouldReturn: nil,
		},
		{
			name:      "Internal error",
			trainerId: 1,
			userId:    2,
			mockBehaviour: func(trainerId, userId int64) {
				rowsTrainer := sqlmock.NewRows([]string{"id", "role"}).AddRow(int64(1), entity.TrainerRole)
				rowsPartnership := sqlmock.NewRows([]string{"id", "user_id", "trainer_id", "status"}).
					AddRow(int64(1), int64(2), int64(1), entity.StatusApproved)
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(trainerId).WillReturnRows(rowsTrainer)
				mock.ExpectQuery("SELECT (.+) FROM partnerships").
					WithArgs(trainerId, userId).WillReturnRows(rowsPartnership)
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(userId).WillReturnError(errors.New("internal error"))
			},
			shouldFail:   true,
			shouldReturn: nil,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			test.mockBehaviour(test.trainerId, test.userId)

			r := NewUserRepository(db)
			got, err := r.GetTrainerUserById(test.trainerId, test.userId)
			if test.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, got, test.shouldReturn)
		})
	}
}

func TestUserRepository_GetTrainerRequestById(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type mockBehaviour func(trainerId, requestId int64)

	table := []struct {
		name          string
		trainerId     int64
		requestId     int64
		mockBehaviour mockBehaviour
		shouldFail    bool
		shouldReturn  *entity.Request
	}{
		{
			name:      "Ok",
			trainerId: 1,
			requestId: 2,
			mockBehaviour: func(trainerId, requestId int64) {
				rowsPartnership := sqlmock.NewRows([]string{"id", "user_id", "trainer_id", "status"}).
					AddRow(int64(1), int64(2), int64(1), entity.StatusRequest)
				rowsSelect := sqlmock.NewRows([]string{"request_id", "user_id", "email", "name", "surname"}).
					AddRow(int64(1), int64(2), "test1", "test1", "test1")
				mock.ExpectQuery("SELECT (.+) FROM partnerships").
					WithArgs(requestId).WillReturnRows(rowsPartnership)
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(requestId).WillReturnRows(rowsSelect)
			},
			shouldFail:   false,
			shouldReturn: &entity.Request{RequestId: 1, UserId: 2, Email: "test1", Name: "test1", Surname: "test1"},
		},
		{
			name:      "No request",
			trainerId: 1,
			requestId: 2,
			mockBehaviour: func(trainerId, requestId int64) {
				rowsPartnership := sqlmock.NewRows([]string{"id", "user_id", "trainer_id", "status"}).
					AddRow(int64(1), int64(2), int64(1), entity.StatusEndedByUser)
				mock.ExpectQuery("SELECT (.+) FROM partnerships").
					WithArgs(requestId).WillReturnRows(rowsPartnership)
			},
			shouldFail:   true,
			shouldReturn: nil,
		},
		{
			name:      "No access to request",
			trainerId: 1,
			requestId: 2,
			mockBehaviour: func(trainerId, requestId int64) {
				rowsPartnership := sqlmock.NewRows([]string{"id", "user_id", "trainer_id", "status"}).
					AddRow(int64(1), int64(4), int64(1), entity.StatusRequest)
				mock.ExpectQuery("SELECT (.+) FROM partnerships").
					WithArgs(requestId).WillReturnRows(rowsPartnership)
			},
			shouldFail:   true,
			shouldReturn: nil,
		},
		{
			name:      "Internal error",
			trainerId: 1,
			requestId: 2,
			mockBehaviour: func(trainerId, requestId int64) {
				rowsPartnership := sqlmock.NewRows([]string{"id", "user_id", "trainer_id", "status"}).
					AddRow(int64(1), int64(2), int64(1), entity.StatusRequest)
				mock.ExpectQuery("SELECT (.+) FROM partnerships").
					WithArgs(requestId).WillReturnRows(rowsPartnership)
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(requestId).WillReturnError(errors.New("internal error"))
			},
			shouldFail:   true,
			shouldReturn: nil,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			test.mockBehaviour(test.trainerId, test.requestId)

			r := NewUserRepository(db)
			got, err := r.GetTrainerRequestById(test.trainerId, test.requestId)
			if test.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, got, test.shouldReturn)
		})
	}
}

func TestUserRepository_InitPartnershipWithUser(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type mockBehaviour func(trainerId, userId int64)

	table := []struct {
		name          string
		trainerId     int64
		userId        int64
		mockBehaviour mockBehaviour
		shouldFail    bool
		shouldReturn  int64
	}{
		{
			name:      "Ok",
			trainerId: 1,
			userId:    2,
			mockBehaviour: func(trainerId, userId int64) {
				rowsUser := sqlmock.NewRows([]string{"id"}).AddRow(int64(2))
				rowsPartnership := sqlmock.NewRows([]string{"id", "user_id", "trainer_id", "status"})
				rowsResult := sqlmock.NewRows([]string{"id"}).AddRow(int64(1))
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(userId).WillReturnRows(rowsUser)
				mock.ExpectQuery("SELECT (.+) FROM partnerships").
					WithArgs(trainerId, userId).WillReturnRows(rowsPartnership)
				mock.ExpectQuery("INSERT INTO partnerships").
					WithArgs(trainerId, userId).WillReturnRows(rowsResult)
			},
			shouldFail:   false,
			shouldReturn: 1,
		},
		{
			name:      "Ok (Approved)",
			trainerId: 1,
			userId:    2,
			mockBehaviour: func(trainerId, userId int64) {
				rowsUser := sqlmock.NewRows([]string{"id"}).AddRow(int64(2))
				rowsPartnership := sqlmock.NewRows([]string{"id", "user_id", "trainer_id", "status"}).
					AddRow(int64(1), int64(2), int64(1), entity.StatusApproved)
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(userId).WillReturnRows(rowsUser)
				mock.ExpectQuery("SELECT (.+) FROM partnerships").
					WithArgs(trainerId, userId).WillReturnRows(rowsPartnership)
			},
			shouldFail:   false,
			shouldReturn: 1,
		},
		{
			name:      "Ok (Ended by trainer)",
			trainerId: 1,
			userId:    2,
			mockBehaviour: func(trainerId, userId int64) {
				rowsUser := sqlmock.NewRows([]string{"id"}).AddRow(int64(2))
				rowsPartnership := sqlmock.NewRows([]string{"id", "user_id", "trainer_id", "status"}).
					AddRow(int64(1), int64(2), int64(1), entity.StatusEndedByTrainer)
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(userId).WillReturnRows(rowsUser)
				mock.ExpectQuery("SELECT (.+) FROM partnerships").
					WithArgs(trainerId, userId).WillReturnRows(rowsPartnership)
				mock.ExpectExec("UPDATE partnerships SET status").
					WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			shouldFail:   false,
			shouldReturn: 1,
		},
		{
			name:      "Ok (request)",
			trainerId: 1,
			userId:    2,
			mockBehaviour: func(trainerId, userId int64) {
				rowsUser := sqlmock.NewRows([]string{"id"}).AddRow(int64(2))
				rowsPartnership := sqlmock.NewRows([]string{"id", "user_id", "trainer_id", "status"}).
					AddRow(int64(1), int64(2), int64(1), entity.StatusRequest)
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(userId).WillReturnRows(rowsUser)
				mock.ExpectQuery("SELECT (.+) FROM partnerships").
					WithArgs(trainerId, userId).WillReturnRows(rowsPartnership)
				mock.ExpectExec("UPDATE partnerships SET status").
					WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			shouldFail:   false,
			shouldReturn: 1,
		},
		{
			name:      "Bad request (Ended by User)",
			trainerId: 1,
			userId:    2,
			mockBehaviour: func(trainerId, userId int64) {
				rowsUser := sqlmock.NewRows([]string{"id"}).AddRow(int64(2))
				rowsPartnership := sqlmock.NewRows([]string{"id", "user_id", "trainer_id", "status"}).
					AddRow(int64(1), int64(2), int64(1), entity.StatusEndedByUser)
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(userId).WillReturnRows(rowsUser)
				mock.ExpectQuery("SELECT (.+) FROM partnerships").
					WithArgs(trainerId, userId).WillReturnRows(rowsPartnership)
			},
			shouldFail:   true,
			shouldReturn: 1,
		},
		{
			name:      "Bad request (Undefined status)",
			trainerId: 1,
			userId:    2,
			mockBehaviour: func(trainerId, userId int64) {
				rowsUser := sqlmock.NewRows([]string{"id"}).AddRow(int64(2))
				rowsPartnership := sqlmock.NewRows([]string{"id", "user_id", "trainer_id", "status"}).
					AddRow(int64(1), int64(2), int64(1), "")
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(userId).WillReturnRows(rowsUser)
				mock.ExpectQuery("SELECT (.+) FROM partnerships").
					WithArgs(trainerId, userId).WillReturnRows(rowsPartnership)
			},
			shouldFail:   true,
			shouldReturn: 0,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			test.mockBehaviour(test.trainerId, test.userId)

			r := NewUserRepository(db)
			got, err := r.InitPartnershipWithUser(test.trainerId, test.userId)
			if test.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, got, test.shouldReturn)
		})
	}
}

func TestUserRepository_EndPartnershipWithUser(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type mockBehaviour func(trainerId, userId int64)

	table := []struct {
		name          string
		trainerId     int64
		userId        int64
		mockBehaviour mockBehaviour
		shouldFail    bool
		shouldReturn  int64
	}{
		{
			name:      "Ok",
			trainerId: 1,
			userId:    2,
			mockBehaviour: func(trainerId, userId int64) {
				rowsPartnership := sqlmock.NewRows([]string{"id", "user_id", "trainer_id", "status"}).
					AddRow(int64(1), int64(2), int64(1), entity.StatusApproved)
				mock.ExpectQuery("SELECT (.+) FROM partnerships").
					WithArgs(trainerId, userId).WillReturnRows(rowsPartnership)
				mock.ExpectExec("UPDATE partnerships SET").
					WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			shouldFail:   false,
			shouldReturn: 1,
		},
		{
			name:      "No partnership",
			trainerId: 1,
			userId:    2,
			mockBehaviour: func(trainerId, userId int64) {
				rowsPartnership := sqlmock.NewRows([]string{"id", "user_id", "trainer_id", "status"})
				mock.ExpectQuery("SELECT (.+) FROM partnerships").
					WithArgs(trainerId, userId).WillReturnRows(rowsPartnership)
			},
			shouldFail:   true,
			shouldReturn: -1,
		},
		{
			name:      "No approved partnership",
			trainerId: 1,
			userId:    2,
			mockBehaviour: func(trainerId, userId int64) {
				rowsPartnership := sqlmock.NewRows([]string{"id", "user_id", "trainer_id", "status"}).
					AddRow(int64(1), int64(2), int64(1), entity.StatusRequest)
				mock.ExpectQuery("SELECT (.+) FROM partnerships").
					WithArgs(trainerId, userId).WillReturnRows(rowsPartnership)
			},
			shouldFail:   true,
			shouldReturn: -1,
		},
		{
			name:      "Internal error",
			trainerId: 1,
			userId:    2,
			mockBehaviour: func(trainerId, userId int64) {
				rowsPartnership := sqlmock.NewRows([]string{"id", "user_id", "trainer_id", "status"}).
					AddRow(int64(1), int64(2), int64(1), entity.StatusApproved)
				mock.ExpectQuery("SELECT (.+) FROM partnerships").
					WithArgs(trainerId, userId).WillReturnRows(rowsPartnership)
				mock.ExpectExec("UPDATE partnerships SET").
					WithArgs(1).WillReturnError(errors.New("internal error"))
			},
			shouldFail:   true,
			shouldReturn: 0,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			test.mockBehaviour(test.trainerId, test.userId)

			r := NewUserRepository(db)
			got, err := r.EndPartnershipWithUser(test.trainerId, test.userId)
			if test.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, got, test.shouldReturn)
		})
	}
}

func TestUserRepository_AcceptRequest(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type mockBehaviour func(trainerId, requestId int64)

	table := []struct {
		name          string
		trainerId     int64
		requestId     int64
		mockBehaviour mockBehaviour
		shouldFail    bool
		shouldReturn  int64
	}{
		{
			name:      "Ok",
			trainerId: 1,
			requestId: 1,
			mockBehaviour: func(trainerId, requestId int64) {
				rowsPartnership := sqlmock.NewRows([]string{"id"}).
					AddRow(int64(1))
				mock.ExpectQuery("UPDATE partnerships SET").
					WithArgs(trainerId, requestId).WillReturnRows(rowsPartnership)
			},
			shouldFail:   false,
			shouldReturn: 1,
		},
		{
			name:      "No request",
			trainerId: 1,
			requestId: 2,
			mockBehaviour: func(trainerId, requestId int64) {
				mock.ExpectQuery("UPDATE partnerships SET").
					WithArgs(trainerId, requestId).WillReturnError(errors.New("no rows"))
			},
			shouldFail:   true,
			shouldReturn: -1,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			test.mockBehaviour(test.trainerId, test.requestId)

			r := NewUserRepository(db)
			got, err := r.AcceptRequest(test.trainerId, test.requestId)
			if test.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, got, test.shouldReturn)
		})
	}
}

func TestUserRepository_DenyRequest(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type mockBehaviour func(trainerId, requestId int64)

	table := []struct {
		name          string
		trainerId     int64
		requestId     int64
		mockBehaviour mockBehaviour
		shouldFail    bool
	}{
		{
			name:      "Ok",
			trainerId: 1,
			requestId: 1,
			mockBehaviour: func(trainerId, requestId int64) {
				mock.ExpectExec("DELETE FROM partnerships").
					WithArgs(trainerId, requestId).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			shouldFail: false,
		},
		{
			name:      "No request",
			trainerId: 1,
			requestId: 2,
			mockBehaviour: func(trainerId, requestId int64) {
				mock.ExpectExec("DELETE FROM partnerships").
					WithArgs(trainerId, requestId).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			shouldFail: true,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			test.mockBehaviour(test.trainerId, test.requestId)

			r := NewUserRepository(db)
			err := r.DenyRequest(test.trainerId, test.requestId)
			if test.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserRepository_CreateWorkoutAsTrainer(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type mockBehaviour func(workout *entity.Workout)

	table := []struct {
		name          string
		workout       entity.Workout
		mockBehaviour mockBehaviour
		shouldFail    bool
		shouldReturn  int64
	}{
		{
			name: "Ok",
			workout: entity.Workout{
				Id:          1,
				Title:       "test",
				Description: "test",
				UserId:      1,
				TrainerId:   sql.NullInt64{Int64: 2, Valid: true},
				Date:        time.Now(),
			},
			mockBehaviour: func(workout *entity.Workout) {
				rowResult := sqlmock.NewRows([]string{"id"}).AddRow(int64(1))
				rowPartnership := sqlmock.NewRows([]string{"id", "status"}).AddRow(int64(1), entity.StatusApproved)
				mock.ExpectQuery("SELECT (.+) FROM partnerships").WillReturnRows(rowPartnership)
				mock.ExpectQuery("INSERT INTO workouts").
					WithArgs(workout.Title, workout.TrainerId, workout.UserId, workout.Description, workout.Date).
					WillReturnRows(rowResult)
			},
			shouldFail:   false,
			shouldReturn: int64(1),
		},
		{
			name: "No partnership",
			workout: entity.Workout{
				Id:          1,
				Title:       "test",
				Description: "test",
				UserId:      1,
				TrainerId:   sql.NullInt64{Int64: 2, Valid: true},
				Date:        time.Now(),
			},
			mockBehaviour: func(workout *entity.Workout) {
				rowPartnership := sqlmock.NewRows([]string{"id", "status"})
				mock.ExpectQuery("SELECT (.+) FROM partnerships").WillReturnRows(rowPartnership)
			},
			shouldFail:   true,
			shouldReturn: int64(-1),
		},
		{
			name: "No approved partnership",
			workout: entity.Workout{
				Id:          1,
				Title:       "test",
				Description: "test",
				UserId:      1,
				TrainerId:   sql.NullInt64{Int64: 2, Valid: true},
				Date:        time.Now(),
			},
			mockBehaviour: func(workout *entity.Workout) {
				rowPartnership := sqlmock.NewRows([]string{"id", "status"}).AddRow(int64(1), entity.StatusEndedByUser)
				mock.ExpectQuery("SELECT (.+) FROM partnerships").WillReturnRows(rowPartnership)
			},
			shouldFail:   true,
			shouldReturn: int64(-1),
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			r := NewUserRepository(db)
			test.mockBehaviour(&test.workout)

			got, err := r.CreateWorkoutAsTrainer(&test.workout)
			if test.shouldFail {
				assert.Error(t, err)
				assert.Equal(t, got, test.shouldReturn)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, got, test.shouldReturn)
			}
		})
	}
}

func TestUserRepository_GetTrainerWorkouts(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type mockBehaviour func(trainerId int64)

	table := []struct {
		name          string
		trainerId     int64
		mockBehaviour mockBehaviour
		shouldFail    bool
		shouldReturn  []*entity.Workout
	}{
		{
			name:      "Ok",
			trainerId: 1,
			mockBehaviour: func(trainerId int64) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "trainer_id", "title"}).
					AddRow(int64(1), int64(2), int64(1), "test1").
					AddRow(int64(2), int64(4), int64(1), "test2").
					AddRow(int64(3), int64(3), int64(1), "test3")
				mock.ExpectQuery("SELECT (.+) FROM workouts").
					WithArgs(trainerId).WillReturnRows(rows)
			},
			shouldFail: false,
			shouldReturn: []*entity.Workout{
				{Id: 1, UserId: 2, TrainerId: sql.NullInt64{1, true}, Title: "test1"},
				{Id: 2, UserId: 4, TrainerId: sql.NullInt64{1, true}, Title: "test2"},
				{Id: 3, UserId: 3, TrainerId: sql.NullInt64{1, true}, Title: "test3"},
			},
		},
		{
			name:      "Empty output",
			trainerId: 1,
			mockBehaviour: func(trainerId int64) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "trainer_id", "title"})
				mock.ExpectQuery("SELECT (.+) FROM workouts").
					WithArgs(trainerId).WillReturnRows(rows)
			},
			shouldFail:   false,
			shouldReturn: []*entity.Workout{},
		},
		{
			name:      "Internal error",
			trainerId: 1,
			mockBehaviour: func(trainerId int64) {
				mock.ExpectQuery("SELECT (.+) FROM workouts").
					WithArgs(trainerId).WillReturnError(errors.New("internal error"))
			},
			shouldFail:   true,
			shouldReturn: nil,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			r := NewUserRepository(db)
			test.mockBehaviour(test.trainerId)

			got, err := r.GetTrainerWorkouts(test.trainerId)
			if test.shouldFail {
				assert.Error(t, err)
				assert.Equal(t, got, test.shouldReturn)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, got, test.shouldReturn)
			}
		})
	}
}

func TestUserRepository_GetTrainerWorkoutsWithUser(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type mockBehaviour func(trainerId, userId int64)

	table := []struct {
		name          string
		trainerId     int64
		userId        int64
		mockBehaviour mockBehaviour
		shouldFail    bool
		shouldReturn  []*entity.Workout
	}{
		{
			name:      "Ok",
			trainerId: 1,
			userId:    2,
			mockBehaviour: func(trainerId, userId int64) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "trainer_id", "title"}).
					AddRow(int64(1), int64(2), int64(1), "test1").
					AddRow(int64(2), int64(2), int64(1), "test2").
					AddRow(int64(3), int64(2), int64(1), "test3")
				mock.ExpectQuery("SELECT (.+) FROM workouts").
					WithArgs(trainerId, userId).WillReturnRows(rows)
			},
			shouldFail: false,
			shouldReturn: []*entity.Workout{
				{Id: 1, UserId: 2, TrainerId: sql.NullInt64{1, true}, Title: "test1"},
				{Id: 2, UserId: 2, TrainerId: sql.NullInt64{1, true}, Title: "test2"},
				{Id: 3, UserId: 2, TrainerId: sql.NullInt64{1, true}, Title: "test3"},
			},
		},
		{
			name:      "Empty output",
			trainerId: 1,
			mockBehaviour: func(trainerId, userId int64) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "trainer_id", "title"})
				mock.ExpectQuery("SELECT (.+) FROM workouts").
					WithArgs(trainerId, userId).WillReturnRows(rows)
			},
			shouldFail:   false,
			shouldReturn: []*entity.Workout{},
		},
		{
			name:      "Internal error",
			trainerId: 1,
			mockBehaviour: func(trainerId, userId int64) {
				mock.ExpectQuery("SELECT (.+) FROM workouts").
					WithArgs(trainerId, userId).WillReturnError(errors.New("internal error"))
			},
			shouldFail:   true,
			shouldReturn: nil,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			r := NewUserRepository(db)
			test.mockBehaviour(test.trainerId, test.userId)

			got, err := r.GetTrainerWorkoutsWithUser(test.trainerId, test.userId)
			if test.shouldFail {
				assert.Error(t, err)
				assert.Equal(t, got, test.shouldReturn)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, got, test.shouldReturn)
			}
		})
	}
}
