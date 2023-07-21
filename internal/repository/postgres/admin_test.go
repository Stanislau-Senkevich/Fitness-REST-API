package postgres

import (
	"Fitness_REST_API/internal/entity"
	"errors"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
)

func TestAdminRepository_Authorize(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type mockBehavior func(admin entity.Admin)

	type args struct {
		login    string
		password string
	}

	table := []struct {
		name         string
		admin        entity.Admin
		args         args
		mockBehavior mockBehavior
		shouldFail   bool
	}{
		{
			name:  "Ok",
			admin: entity.Admin{Id: 1, Login: "test", PasswordHash: "test"},
			args:  args{login: "test", password: "test"},
			mockBehavior: func(admin entity.Admin) {
				rows := sqlmock.NewRows([]string{"id", "login", "password_hash"}).AddRow(admin.Id, admin.Login, admin.PasswordHash)
				mock.ExpectQuery("SELECT (.+) FROM admins").WillReturnRows(rows)
			},
		},
		{
			name:  "Empty fields",
			admin: entity.Admin{Id: 1, Login: "test", PasswordHash: "test"},
			args:  args{login: "", password: ""},
			mockBehavior: func(admin entity.Admin) {
				mock.ExpectQuery("SELECT (.+) FROM admins").WillReturnError(errors.New("no rows selected"))
			},
			shouldFail: true,
		},
	}

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior(test.admin)

			r := NewAdminRepository(db)

			err := r.Authorize(test.args.login, test.args.password)
			if test.shouldFail {
				assert.Error(t, err)
				t.Skip("OK")
			}

			assert.NoError(t, err)
		})
	}

}
