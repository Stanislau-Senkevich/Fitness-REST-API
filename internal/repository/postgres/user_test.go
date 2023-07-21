package postgres

import (
	"Fitness_REST_API/internal/entity"
	"database/sql"
	"errors"
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
				TrainerId:   sql.NullInt64{2, true},
				Date:        time.Now(),
			},
			mockBehaviour: func(workout *entity.Workout) {
				row := sqlmock.NewRows([]string{"id"}).AddRow(workout.Id)
				rowTrainer := sqlmock.NewRows([]string{"id", "role"}).AddRow(workout.Id, entity.TrainerRole)
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
				TrainerId:   sql.NullInt64{2, true},
				Date:        time.Now(),
			},
			mockBehaviour: func(workout *entity.Workout) {
				rowTrainer := sqlmock.NewRows([]string{"id", "role"}).AddRow(workout.Id, entity.TrainerRole)
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
				TrainerId:   sql.NullInt64{2, true},
				Date:        time.Now(),
			},
			mockBehaviour: func(workout *entity.Workout) {
				rowTrainer := sqlmock.NewRows([]string{"id", "role"}).AddRow(workout.Id, entity.UserRole)
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
				TrainerId:   sql.NullInt64{2, true},
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
	}
}
