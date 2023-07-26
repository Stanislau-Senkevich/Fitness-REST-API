package handler

import (
	"Fitness_REST_API/internal/entity"
	"Fitness_REST_API/internal/service"
	mockService "Fitness_REST_API/internal/service/mocks"
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler_getAllUsersFullInfo(t *testing.T) {
	type mockBehaviour func(r *mockService.MockAdmin)

	table := []struct {
		name                 string
		mockBehaviour        mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			mockBehaviour: func(r *mockService.MockAdmin) {
				r.EXPECT().GetUsersId(entity.UserRole).Return([]int64{1, 2}, nil)
				r.EXPECT().GetUserFullInfoById(int64(1)).
					Return(&entity.UserInfo{
						Id: 1, Email: "test", Role: entity.UserRole, Name: "test", Surname: "test",
						Partnerships: []*entity.Partnership{
							{Id: 1, UserId: 1, TrainerId: 3, Status: entity.StatusApproved},
							{Id: 2, UserId: 1, TrainerId: 5, Status: entity.StatusRequest},
						},
						Workouts: []*entity.Workout{
							{Id: 1, UserId: 1, TrainerId: sql.NullInt64{Int64: 2, Valid: true}, Title: "test"},
							{Id: 2, UserId: 1, TrainerId: sql.NullInt64{Int64: 0, Valid: false}, Title: "test"},
						},
					}, nil)
				r.EXPECT().GetUserFullInfoById(int64(2)).
					Return(&entity.UserInfo{
						Id: 2, Email: "test", Role: entity.UserRole, Name: "test", Surname: "test",
						Partnerships: []*entity.Partnership{
							{Id: 3, UserId: 2, TrainerId: 10, Status: entity.StatusApproved},
							{Id: 4, UserId: 2, TrainerId: 8, Status: entity.StatusRequest},
						},
						Workouts: []*entity.Workout{
							{Id: 3, UserId: 2, TrainerId: sql.NullInt64{Int64: 8, Valid: true}, Title: "test"},
							{Id: 4, UserId: 2, TrainerId: sql.NullInt64{Int64: 0, Valid: false}, Title: "test"},
						},
					}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"id":1,"email":"test","role":"user","name":"test","surname":"test","created_at":"0001-01-01T00:00:00Z","partnerships":[{"id":1,"user_id":1,"trainer_id":3,"status":"approved","created_at":"0001-01-01T00:00:00Z","ended_at":{"Time":"0001-01-01T00:00:00Z","Valid":false}},{"id":2,"user_id":1,"trainer_id":5,"status":"request","created_at":"0001-01-01T00:00:00Z","ended_at":{"Time":"0001-01-01T00:00:00Z","Valid":false}}],"workouts":[{"id":1,"title":"test","user_id":1,"trainer_id":{"Int64":2,"Valid":true},"date":"0001-01-01T00:00:00Z"},{"id":2,"title":"test","user_id":1,"trainer_id":{"Int64":0,"Valid":false},"date":"0001-01-01T00:00:00Z"}]},{"id":2,"email":"test","role":"user","name":"test","surname":"test","created_at":"0001-01-01T00:00:00Z","partnerships":[{"id":3,"user_id":2,"trainer_id":10,"status":"approved","created_at":"0001-01-01T00:00:00Z","ended_at":{"Time":"0001-01-01T00:00:00Z","Valid":false}},{"id":4,"user_id":2,"trainer_id":8,"status":"request","created_at":"0001-01-01T00:00:00Z","ended_at":{"Time":"0001-01-01T00:00:00Z","Valid":false}}],"workouts":[{"id":3,"title":"test","user_id":2,"trainer_id":{"Int64":8,"Valid":true},"date":"0001-01-01T00:00:00Z"},{"id":4,"title":"test","user_id":2,"trainer_id":{"Int64":0,"Valid":false},"date":"0001-01-01T00:00:00Z"}]}]`,
		},
		{
			name: "No users",
			mockBehaviour: func(r *mockService.MockAdmin) {
				r.EXPECT().GetUsersId(entity.UserRole).Return([]int64{}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[]`,
		},
		{
			name: "Internal error",
			mockBehaviour: func(r *mockService.MockAdmin) {
				r.EXPECT().GetUsersId(entity.UserRole).Return(nil, errors.New("internal error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"internal error"}`,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mockService.NewMockAdmin(c)
			test.mockBehaviour(repo)

			services := &service.Services{Admin: repo}
			handler := &Handler{services: services}

			r := gin.New()
			r.GET("/user", handler.getAllUsersFullInfo)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/user", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_getTrainersInfo(t *testing.T) {
	type mockBehaviour func(r *mockService.MockAdmin)

	table := []struct {
		name                 string
		mockBehaviour        mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			mockBehaviour: func(r *mockService.MockAdmin) {
				r.EXPECT().GetUsersId(entity.TrainerRole).Return([]int64{1, 2}, nil)
				r.EXPECT().GetUserFullInfoById(int64(1)).
					Return(&entity.UserInfo{
						Id: 1, Email: "test", Role: entity.TrainerRole, Name: "test", Surname: "test",
						Partnerships: []*entity.Partnership{
							{Id: 1, UserId: 10, TrainerId: 1, Status: entity.StatusApproved},
							{Id: 2, UserId: 11, TrainerId: 1, Status: entity.StatusRequest},
						},
						Workouts: []*entity.Workout{
							{Id: 1, UserId: 10, TrainerId: sql.NullInt64{Int64: 1, Valid: true}, Title: "test"},
							{Id: 2, UserId: 11, TrainerId: sql.NullInt64{Int64: 1, Valid: true}, Title: "test"},
						},
					}, nil)
				r.EXPECT().GetUserFullInfoById(int64(2)).
					Return(&entity.UserInfo{
						Id: 2, Email: "test", Role: entity.TrainerRole, Name: "test", Surname: "test",
						Partnerships: []*entity.Partnership{
							{Id: 3, UserId: 22, TrainerId: 2, Status: entity.StatusApproved},
							{Id: 4, UserId: 20, TrainerId: 2, Status: entity.StatusRequest},
						},
						Workouts: []*entity.Workout{
							{Id: 3, UserId: 22, TrainerId: sql.NullInt64{Int64: 2, Valid: true}, Title: "test"},
							{Id: 4, UserId: 20, TrainerId: sql.NullInt64{Int64: 2, Valid: true}, Title: "test"},
						},
					}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"id":1,"email":"test","role":"trainer","name":"test","surname":"test","created_at":"0001-01-01T00:00:00Z","partnerships":[{"id":1,"user_id":10,"trainer_id":1,"status":"approved","created_at":"0001-01-01T00:00:00Z","ended_at":{"Time":"0001-01-01T00:00:00Z","Valid":false}},{"id":2,"user_id":11,"trainer_id":1,"status":"request","created_at":"0001-01-01T00:00:00Z","ended_at":{"Time":"0001-01-01T00:00:00Z","Valid":false}}],"workouts":[{"id":1,"title":"test","user_id":10,"trainer_id":{"Int64":1,"Valid":true},"date":"0001-01-01T00:00:00Z"},{"id":2,"title":"test","user_id":11,"trainer_id":{"Int64":1,"Valid":true},"date":"0001-01-01T00:00:00Z"}]},{"id":2,"email":"test","role":"trainer","name":"test","surname":"test","created_at":"0001-01-01T00:00:00Z","partnerships":[{"id":3,"user_id":22,"trainer_id":2,"status":"approved","created_at":"0001-01-01T00:00:00Z","ended_at":{"Time":"0001-01-01T00:00:00Z","Valid":false}},{"id":4,"user_id":20,"trainer_id":2,"status":"request","created_at":"0001-01-01T00:00:00Z","ended_at":{"Time":"0001-01-01T00:00:00Z","Valid":false}}],"workouts":[{"id":3,"title":"test","user_id":22,"trainer_id":{"Int64":2,"Valid":true},"date":"0001-01-01T00:00:00Z"},{"id":4,"title":"test","user_id":20,"trainer_id":{"Int64":2,"Valid":true},"date":"0001-01-01T00:00:00Z"}]}]`,
		},
		{
			name: "No users",
			mockBehaviour: func(r *mockService.MockAdmin) {
				r.EXPECT().GetUsersId(entity.TrainerRole).Return([]int64{}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[]`,
		},
		{
			name: "Internal error",
			mockBehaviour: func(r *mockService.MockAdmin) {
				r.EXPECT().GetUsersId(entity.TrainerRole).Return(nil, errors.New("internal error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"internal error"}`,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mockService.NewMockAdmin(c)
			test.mockBehaviour(repo)

			services := &service.Services{Admin: repo}
			handler := &Handler{services: services}

			r := gin.New()
			r.GET("/user", handler.getTrainersInfo)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/user", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_getUserFullInfoById(t *testing.T) {
	type mockBehaviour func(r *mockService.MockAdmin, userId int64)

	table := []struct {
		name                 string
		userId               int64
		mockBehaviour        mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:   "Ok",
			userId: 1,
			mockBehaviour: func(r *mockService.MockAdmin, userId int64) {
				r.EXPECT().GetUserFullInfoById(userId).
					Return(&entity.UserInfo{
						Id: 1, Email: "test", Role: entity.UserRole, Name: "test", Surname: "test",
						Partnerships: []*entity.Partnership{
							{Id: 1, UserId: 1, TrainerId: 3, Status: entity.StatusApproved},
							{Id: 2, UserId: 1, TrainerId: 5, Status: entity.StatusRequest},
						},
						Workouts: []*entity.Workout{
							{Id: 1, UserId: 1, TrainerId: sql.NullInt64{Int64: 2, Valid: true}, Title: "test"},
							{Id: 2, UserId: 1, TrainerId: sql.NullInt64{Int64: 0, Valid: false}, Title: "test"},
						},
					}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1,"email":"test","role":"user","name":"test","surname":"test","created_at":"0001-01-01T00:00:00Z","partnerships":[{"id":1,"user_id":1,"trainer_id":3,"status":"approved","created_at":"0001-01-01T00:00:00Z","ended_at":{"Time":"0001-01-01T00:00:00Z","Valid":false}},{"id":2,"user_id":1,"trainer_id":5,"status":"request","created_at":"0001-01-01T00:00:00Z","ended_at":{"Time":"0001-01-01T00:00:00Z","Valid":false}}],"workouts":[{"id":1,"title":"test","user_id":1,"trainer_id":{"Int64":2,"Valid":true},"date":"0001-01-01T00:00:00Z"},{"id":2,"title":"test","user_id":1,"trainer_id":{"Int64":0,"Valid":false},"date":"0001-01-01T00:00:00Z"}]}`,
		},
		{
			name:                 "Invalid id param",
			userId:               -1,
			mockBehaviour:        func(r *mockService.MockAdmin, userId int64) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid id param"}`,
		},
		{
			name:   "Invalid id",
			userId: 1,
			mockBehaviour: func(r *mockService.MockAdmin, userId int64) {
				r.EXPECT().GetUserFullInfoById(userId).
					Return(nil, errors.New("bad request"))
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"bad request"}`,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mockService.NewMockAdmin(c)
			test.mockBehaviour(repo, test.userId)

			services := &service.Services{Admin: repo}
			handler := &Handler{services: services}

			r := gin.New()
			r.GET("/user/:id", handler.getUserFullInfoByID)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/user/%d", test.userId), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_createUser(t *testing.T) {
	type mockBehavior func(r *mockService.MockAdmin, inputUser entity.User)

	table := []struct {
		name                 string
		inputBody            string
		inputUser            entity.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok User",
			inputBody: `{"email":"testEmail", "password":"testPassword", "name":"testName", "surname":"testSurname"}`,
			inputUser: entity.User{Email: "testEmail", PasswordHash: "testPassword", Role: entity.UserRole, Name: "testName", Surname: "testSurname"},
			mockBehavior: func(r *mockService.MockAdmin, inputUser entity.User) {
				r.EXPECT().CreateUser(&inputUser).Return(int64(1), nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:      "Ok Trainer",
			inputBody: `{"email":"testEmail", "password":"testPassword", "role":"trainer", "name":"testName", "surname":"testSurname"}`,
			inputUser: entity.User{Email: "testEmail", PasswordHash: "testPassword", Role: entity.TrainerRole, Name: "testName", Surname: "testSurname"},
			mockBehavior: func(r *mockService.MockAdmin, inputUser entity.User) {
				r.EXPECT().CreateUser(&inputUser).Return(int64(1), nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:                 "Invalid email JSON",
			inputBody:            `{"emailA":"testEmail", "password":"testPassword", "name":"testName", "surname":"testSurname"}`,
			inputUser:            entity.User{Email: "testEmail", PasswordHash: "testPassword", Role: entity.UserRole, Name: "testName", Surname: "testSurname"},
			mockBehavior:         func(r *mockService.MockAdmin, inputUser entity.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"Key: 'User.Email' Error:Field validation for 'Email' failed on the 'required' tag"}`,
		},
		{
			name:                 "Invalid password JSON",
			inputBody:            `{"email":"testEmail", "passwordA":"testPassword", "name":"testName", "surname":"testSurname"}`,
			inputUser:            entity.User{Email: "testEmail", PasswordHash: "testPassword", Role: entity.UserRole, Name: "testName", Surname: "testSurname"},
			mockBehavior:         func(r *mockService.MockAdmin, inputUser entity.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"Key: 'User.PasswordHash' Error:Field validation for 'PasswordHash' failed on the 'required' tag"}`,
		},
		{
			name:                 "Invalid name JSON",
			inputBody:            `{"email":"testEmail", "password":"testPassword", "nameA":"testName", "surname":"testSurname"}`,
			inputUser:            entity.User{Email: "testEmail", PasswordHash: "testPassword", Name: "testName", Surname: "testSurname"},
			mockBehavior:         func(r *mockService.MockAdmin, inputUser entity.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"Key: 'User.Name' Error:Field validation for 'Name' failed on the 'required' tag"}`,
		},
		{
			name:                 "Invalid surname JSON",
			inputBody:            `{"email":"testEmail", "password":"testPassword", "name":"testName", "surnameA":"testSurname"}`,
			inputUser:            entity.User{Email: "testEmail", PasswordHash: "testPassword", Name: "testName", Surname: "testSurname"},
			mockBehavior:         func(r *mockService.MockAdmin, inputUser entity.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"Key: 'User.Surname' Error:Field validation for 'Surname' failed on the 'required' tag"}`,
		},
		{
			name:      "Internal Server Error",
			inputBody: `{"email":"testEmail", "password":"testPassword", "name":"testName", "surname":"testSurname"}`,
			inputUser: entity.User{Email: "testEmail", PasswordHash: "testPassword", Role: entity.UserRole, Name: "testName", Surname: "testSurname"},
			mockBehavior: func(r *mockService.MockAdmin, inputUser entity.User) {
				r.EXPECT().CreateUser(&inputUser).Return(int64(0), errors.New("internal error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"internal error"}`,
		},
		{
			name:      "Email has already reserved",
			inputBody: `{"email":"testEmail", "password":"testPassword", "name":"testName", "surname":"testSurname"}`,
			inputUser: entity.User{Email: "testEmail", PasswordHash: "testPassword", Role: entity.UserRole, Name: "testName", Surname: "testSurname"},
			mockBehavior: func(r *mockService.MockAdmin, inputUser entity.User) {
				r.EXPECT().CreateUser(&inputUser).Return(int64(-1), errors.New("reserved email"))
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"reserved email"}`,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mockService.NewMockAdmin(c)
			test.mockBehavior(repo, test.inputUser)

			services := &service.Services{Admin: repo}
			handler := &Handler{services: services}

			r := gin.New()
			r.POST("/user", handler.createUser)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/user",
				bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)

		})
	}
}

func TestHandler_updateUser(t *testing.T) {
	type mockBehaviour func(r *mockService.MockAdmin, u *mockService.MockUser, userId int64)

	table := []struct {
		name               string
		userId             int64
		inputBody          string
		update             *entity.UserUpdate
		mockBehaviour      mockBehaviour
		expectedStatusCode int
	}{
		{
			name:      "Ok",
			userId:    1,
			inputBody: `{"email":"testNew", "password":"qwerty123", "name":"testNew", "surname":"testNew"}`,
			mockBehaviour: func(r *mockService.MockAdmin, u *mockService.MockUser, userId int64) {
				user := &entity.User{
					Id:           1,
					Email:        "testOld",
					PasswordHash: "afcsdv",
					Role:         entity.UserRole,
					Name:         "testOld",
					Surname:      "testOld",
					CreatedAt:    time.Unix(0, 1),
				}
				update := &entity.UserUpdate{
					Email:    "testNew",
					Password: "5cec175b165e3d5e62c9e13ce848ef6feac81bff",
					Role:     entity.UserRole,
					Name:     "testNew",
					Surname:  "testNew",
				}
				u.EXPECT().GetUserInfoById(userId).Return(user, nil)
				u.EXPECT().GetPasswordHash("qwerty123").Return("5cec175b165e3d5e62c9e13ce848ef6feac81bff")
				r.EXPECT().UpdateUser(userId, update).Return(nil)
			},
			expectedStatusCode: 200,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mockService.NewMockAdmin(c)
			us := mockService.NewMockUser(c)
			test.mockBehaviour(repo, us, test.userId)

			services := &service.Services{Admin: repo, User: us}
			handler := &Handler{services: services}

			r := gin.New()
			r.PUT("/user/:id", handler.updateUser)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", fmt.Sprintf("/user/%d", test.userId),
				bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
		})
	}
}

func TestHandler_deleteUser(t *testing.T) {
	type mockBehaviour func(r *mockService.MockAdmin, userId int64)

	table := []struct {
		name               string
		userId             int64
		mockBehaviour      mockBehaviour
		expectedStatusCode int
	}{
		{
			name:   "Ok",
			userId: 1,
			mockBehaviour: func(r *mockService.MockAdmin, userId int64) {
				r.EXPECT().DeleteUser(userId).Return(nil)
			},
			expectedStatusCode: 200,
		},
		{
			name:               "Invalid id param",
			userId:             -1,
			mockBehaviour:      func(r *mockService.MockAdmin, userId int64) {},
			expectedStatusCode: 400,
		},
		{
			name:   "No user to delete",
			userId: 1,
			mockBehaviour: func(r *mockService.MockAdmin, userId int64) {
				r.EXPECT().DeleteUser(userId).Return(errors.New("error due deleting"))
			},
			expectedStatusCode: 400,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mockService.NewMockAdmin(c)
			test.mockBehaviour(repo, test.userId)

			services := &service.Services{Admin: repo}
			handler := &Handler{services: services}

			r := gin.New()
			r.DELETE("/user/:id", handler.deleteUser)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/user/%d", test.userId), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
		})
	}
}
