package handler

import (
	"Fitness_REST_API/internal/entity"
	"Fitness_REST_API/internal/service"
	mock_service "Fitness_REST_API/internal/service/mocks"
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_getTrainerUsers(t *testing.T) {
	type mockBehaviour func(r *mock_service.MockUser, trainerId int64)

	table := []struct {
		name                 string
		trainerId            int64
		mockBehaviour        mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			trainerId: 1,
			mockBehaviour: func(r *mock_service.MockUser, trainerId int64) {
				r.EXPECT().GetTrainerUsers(trainerId).Return([]*entity.User{
					{Id: 100, Email: "test1", Name: "test1", Surname: "test1"},
					{Id: 101, Email: "test2", Name: "test2", Surname: "test2"},
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"users":[{"id":100,"email":"test1","name":"test1","surname":"test1","created_at":"0001-01-01T00:00:00Z"},{"id":101,"email":"test2","name":"test2","surname":"test2","created_at":"0001-01-01T00:00:00Z"}]}`, //nolint
		},
		{
			name:                 "Invalid id",
			trainerId:            -1,
			mockBehaviour:        func(r *mock_service.MockUser, trainerId int64) {},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"invalid id"}`,
		},
		{
			name:      "No users",
			trainerId: 1,
			mockBehaviour: func(r *mock_service.MockUser, trainerId int64) {
				r.EXPECT().GetTrainerUsers(trainerId).Return([]*entity.User{}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"users":[]}`,
		},
		{
			name:      "Internal error",
			trainerId: 1,
			mockBehaviour: func(r *mock_service.MockUser, trainerId int64) {
				r.EXPECT().GetTrainerUsers(trainerId).Return(nil, errors.New("internal error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"internal error"}`,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_service.NewMockUser(c)
			test.mockBehaviour(user, test.trainerId)

			services := &service.Services{User: user}
			handler := &Handler{services: services}

			r := gin.New()
			r.GET("/user", handler.getTrainerUsers)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/user", nil)
			ctx, _ := gin.CreateTestContext(w)
			ctx.Set(userIdCtx, test.trainerId)
			r.ServeHTTP(w, req.WithContext(ctx))

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_getTrainerRequests(t *testing.T) {
	type mockBehaviour func(r *mock_service.MockUser, trainerId int64)

	table := []struct {
		name                 string
		trainerId            int64
		mockBehaviour        mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			trainerId: 1,
			mockBehaviour: func(r *mock_service.MockUser, trainerId int64) {
				r.EXPECT().GetTrainerRequests(trainerId).Return([]*entity.Request{
					{RequestId: 1, UserId: 100, Email: "test", Name: "test", Surname: "test"},
					{RequestId: 2, UserId: 101, Email: "test", Name: "test", Surname: "test"},
					{RequestId: 3, UserId: 102, Email: "test", Name: "test", Surname: "test"},
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"request_id":1,"user_id":100,"email":"test","name":"test","surname":"test","send_at":"0001-01-01T00:00:00Z"},{"request_id":2,"user_id":101,"email":"test","name":"test","surname":"test","send_at":"0001-01-01T00:00:00Z"},{"request_id":3,"user_id":102,"email":"test","name":"test","surname":"test","send_at":"0001-01-01T00:00:00Z"}]`, //nolint
		},
		{
			name:                 "Invalid id",
			trainerId:            -1,
			mockBehaviour:        func(r *mock_service.MockUser, trainerId int64) {},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"invalid id"}`,
		},
		{
			name:      "No requests",
			trainerId: 1,
			mockBehaviour: func(r *mock_service.MockUser, trainerId int64) {
				r.EXPECT().GetTrainerRequests(trainerId).Return([]*entity.Request{}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[]`,
		},
		{
			name:      "Internal error",
			trainerId: 1,
			mockBehaviour: func(r *mock_service.MockUser, trainerId int64) {
				r.EXPECT().GetTrainerRequests(trainerId).Return(nil, errors.New("internal error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"internal error"}`,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_service.NewMockUser(c)
			test.mockBehaviour(user, test.trainerId)

			services := &service.Services{User: user}
			handler := &Handler{services: services}

			r := gin.New()
			r.GET("/request", handler.getTrainerRequests)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/request", nil)
			ctx, _ := gin.CreateTestContext(w)
			ctx.Set(userIdCtx, test.trainerId)
			r.ServeHTTP(w, req.WithContext(ctx))

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_getTrainerUserById(t *testing.T) {
	type mockBehaviour func(r *mock_service.MockUser, trainerId, userId int64)

	table := []struct {
		name                 string
		trainerId            int64
		userId               int64
		mockBehaviour        mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			trainerId: 1,
			userId:    2,
			mockBehaviour: func(r *mock_service.MockUser, trainerId, userId int64) {
				r.EXPECT().GetTrainerUserById(trainerId, userId).Return(&entity.User{
					Id: 2, Email: "test", Name: "test", Surname: "test",
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":2,"email":"test","name":"test","surname":"test","created_at":"0001-01-01T00:00:00Z"}`,
		},
		{
			name:                 "Invalid userId",
			trainerId:            1,
			userId:               -2,
			mockBehaviour:        func(r *mock_service.MockUser, trainerId, userId int64) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid id parameter"}`,
		},
		{
			name:                 "Invalid trainerId",
			trainerId:            -1,
			userId:               2,
			mockBehaviour:        func(r *mock_service.MockUser, trainerId, userId int64) {},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"invalid id"}`,
		},
		{
			name:      "No access to user",
			trainerId: 1,
			userId:    2,
			mockBehaviour: func(r *mock_service.MockUser, trainerId, userId int64) {
				r.EXPECT().GetTrainerUserById(trainerId, userId).Return(nil, errors.New("bad request"))
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"bad request"}`,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_service.NewMockUser(c)
			test.mockBehaviour(user, test.trainerId, test.userId)

			services := &service.Services{User: user}
			handler := &Handler{services: services}

			r := gin.New()
			r.GET("/user/:id", handler.getTrainerUserById)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/user/%d", test.userId), nil)
			ctx, _ := gin.CreateTestContext(w)
			ctx.Set(userIdCtx, test.trainerId)
			r.ServeHTTP(w, req.WithContext(ctx))

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_getTrainerRequestById(t *testing.T) {
	type mockBehaviour func(r *mock_service.MockUser, trainerId, requestId int64)

	table := []struct {
		name                 string
		trainerId            int64
		requestId            int64
		mockBehaviour        mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			trainerId: 1,
			requestId: 2,
			mockBehaviour: func(r *mock_service.MockUser, trainerId, requestId int64) {
				r.EXPECT().GetTrainerRequestById(trainerId, requestId).Return(&entity.Request{
					RequestId: 1, UserId: 100, Email: "test", Name: "test", Surname: "test",
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"request_id":1,"user_id":100,"email":"test","name":"test","surname":"test","send_at":"0001-01-01T00:00:00Z"}`, //nolint
		},
		{
			name:                 "Invalid requestId",
			trainerId:            1,
			requestId:            -2,
			mockBehaviour:        func(r *mock_service.MockUser, trainerId, requestId int64) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid id parameter"}`,
		},
		{
			name:                 "Invalid trainerId",
			trainerId:            -1,
			requestId:            2,
			mockBehaviour:        func(r *mock_service.MockUser, trainerId, requestId int64) {},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"invalid id"}`,
		},
		{
			name:      "No request was returned",
			trainerId: 1,
			requestId: 2,
			mockBehaviour: func(r *mock_service.MockUser, trainerId, requestId int64) {
				r.EXPECT().GetTrainerRequestById(trainerId, requestId).Return(nil, errors.New("no request"))
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"no request"}`,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_service.NewMockUser(c)
			test.mockBehaviour(user, test.trainerId, test.requestId)

			services := &service.Services{User: user}
			handler := &Handler{services: services}
			r := gin.New()
			r.GET("/request/:id", handler.getTrainerRequestById)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/request/%d", test.requestId), nil)
			ctx, _ := gin.CreateTestContext(w)
			ctx.Set(userIdCtx, test.trainerId)
			r.ServeHTTP(w, req.WithContext(ctx))

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_getTrainerWorkouts(t *testing.T) {
	type mockBehaviour func(r *mock_service.MockUser, trainerId int64)

	table := []struct {
		name                 string
		trainerId            int64
		mockBehaviour        mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			trainerId: 1,
			mockBehaviour: func(r *mock_service.MockUser, trainerId int64) {
				r.EXPECT().GetTrainerWorkouts(trainerId).Return([]*entity.Workout{
					{Id: 1, TrainerId: sql.NullInt64{Int64: 1, Valid: true}, UserId: 2, Title: "test1"},
					{Id: 2, TrainerId: sql.NullInt64{Int64: 1, Valid: true}, UserId: 3, Title: "test2"},
					{Id: 3, TrainerId: sql.NullInt64{Int64: 1, Valid: true}, UserId: 4, Title: "test3"},
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"workouts":[{"id":1,"title":"test1","user_id":2,"trainer_id":{"Int64":1,"Valid":true},"date":"0001-01-01T00:00:00Z"},{"id":2,"title":"test2","user_id":3,"trainer_id":{"Int64":1,"Valid":true},"date":"0001-01-01T00:00:00Z"},{"id":3,"title":"test3","user_id":4,"trainer_id":{"Int64":1,"Valid":true},"date":"0001-01-01T00:00:00Z"}]}`, //nolint
		},
		{
			name:                 "Invalid id",
			trainerId:            -1,
			mockBehaviour:        func(r *mock_service.MockUser, trainerId int64) {},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"invalid id"}`,
		},
		{
			name:      "No workouts",
			trainerId: 1,
			mockBehaviour: func(r *mock_service.MockUser, trainerId int64) {
				r.EXPECT().GetTrainerWorkouts(trainerId).Return([]*entity.Workout{}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"workouts":[]}`,
		},
		{
			name:      "Internal error",
			trainerId: 1,
			mockBehaviour: func(r *mock_service.MockUser, trainerId int64) {
				r.EXPECT().GetTrainerWorkouts(trainerId).Return(nil, errors.New("internal error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"internal error"}`,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_service.NewMockUser(c)
			test.mockBehaviour(user, test.trainerId)

			services := &service.Services{User: user}
			handler := &Handler{services: services}

			r := gin.New()
			r.GET("/workout", handler.getTrainerWorkouts)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/workout", nil)
			ctx, _ := gin.CreateTestContext(w)
			ctx.Set(userIdCtx, test.trainerId)
			r.ServeHTTP(w, req.WithContext(ctx))

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_getTrainerWorkoutsWithUser(t *testing.T) {
	type mockBehaviour func(r *mock_service.MockUser, trainerId, userId int64)

	table := []struct {
		name                 string
		trainerId            int64
		userId               int64
		mockBehaviour        mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			trainerId: 1,
			userId:    2,
			mockBehaviour: func(r *mock_service.MockUser, trainerId, userId int64) {
				r.EXPECT().GetTrainerWorkoutsWithUser(trainerId, userId).Return([]*entity.Workout{
					{Id: 1, TrainerId: sql.NullInt64{Int64: 1, Valid: true}, UserId: 2, Title: "test1"},
					{Id: 2, TrainerId: sql.NullInt64{Int64: 1, Valid: true}, UserId: 2, Title: "test2"},
					{Id: 3, TrainerId: sql.NullInt64{Int64: 1, Valid: true}, UserId: 2, Title: "test3"},
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"workouts":[{"id":1,"title":"test1","user_id":2,"trainer_id":{"Int64":1,"Valid":true},"date":"0001-01-01T00:00:00Z"},{"id":2,"title":"test2","user_id":2,"trainer_id":{"Int64":1,"Valid":true},"date":"0001-01-01T00:00:00Z"},{"id":3,"title":"test3","user_id":2,"trainer_id":{"Int64":1,"Valid":true},"date":"0001-01-01T00:00:00Z"}]}`, //nolint
		},
		{
			name:                 "Invalid trainerId",
			trainerId:            -1,
			userId:               2,
			mockBehaviour:        func(r *mock_service.MockUser, trainerId, userId int64) {},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"invalid id"}`,
		},
		{
			name:                 "Invalid userid",
			trainerId:            1,
			userId:               -2,
			mockBehaviour:        func(r *mock_service.MockUser, trainerId, userId int64) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid id parameter"}`,
		},
		{
			name:      "No workouts",
			trainerId: 1,
			userId:    500,
			mockBehaviour: func(r *mock_service.MockUser, trainerId, userId int64) {
				r.EXPECT().GetTrainerWorkoutsWithUser(trainerId, userId).Return([]*entity.Workout{}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"workouts":[]}`,
		},
		{
			name:      "Internal error",
			trainerId: 1,
			userId:    2,
			mockBehaviour: func(r *mock_service.MockUser, trainerId, userId int64) {
				r.EXPECT().GetTrainerWorkoutsWithUser(trainerId, userId).Return(nil, errors.New("internal error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"internal error"}`,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_service.NewMockUser(c)
			test.mockBehaviour(user, test.trainerId, test.userId)

			services := &service.Services{User: user}
			handler := &Handler{services: services}

			r := gin.New()
			r.GET("/workout/user/:id", handler.getTrainerWorkoutsWithUser)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/workout/user/%d", test.userId), nil)
			ctx, _ := gin.CreateTestContext(w)
			ctx.Set(userIdCtx, test.trainerId)
			r.ServeHTTP(w, req.WithContext(ctx))

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_createTrainerWorkout(t *testing.T) {
	type mockBehaviour func(r *mock_service.MockUser, inputWorkout entity.Workout, trainerId int64)

	table := []struct {
		name                 string
		trainerId            int64
		inputBody            string
		inputWorkout         entity.Workout
		mockBehaviour        mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:         "Ok",
			trainerId:    1,
			inputBody:    `{"title":"test", "description":"test", "user_id":2}`,
			inputWorkout: entity.Workout{Title: "test", Description: "test", UserId: 2},
			mockBehaviour: func(r *mock_service.MockUser, inputWorkout entity.Workout, trainerId int64) {
				_ = formatTrainerWorkout(&inputWorkout, trainerId)
				r.EXPECT().CreateWorkoutAsTrainer(&inputWorkout).Return(int64(1), nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"workout_id":1}`,
		},
		{
			name:         "Invalid trainer Id in workout",
			trainerId:    1,
			inputBody:    `{"title":"test", "description":"test", "user_id":2, "trainer_id":{"int64":4, "valid":true}}`,
			inputWorkout: entity.Workout{Title: "test", Description: "test", UserId: 2},
			mockBehaviour: func(r *mock_service.MockUser, inputWorkout entity.Workout, trainerId int64) {
				_ = formatTrainerWorkout(&inputWorkout, trainerId)
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"trainer_id from token and trainer_id from workout must match"}`,
		},
		{
			name:                 "Workout without a title (including empty workout)",
			trainerId:            1,
			inputBody:            `{}`,
			inputWorkout:         entity.Workout{},
			mockBehaviour:        func(r *mock_service.MockUser, inputWorkout entity.Workout, trainerId int64) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"Key: 'Workout.Title' Error:Field validation for 'Title' failed on the 'required' tag"}`, //nolint
		},
		{
			name:                 "Empty title",
			trainerId:            1,
			inputBody:            `{"title":"", "description":"test"}`,
			inputWorkout:         entity.Workout{},
			mockBehaviour:        func(r *mock_service.MockUser, inputWorkout entity.Workout, trainerId int64) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"Key: 'Workout.Title' Error:Field validation for 'Title' failed on the 'required' tag"}`, //nolint
		},
		{
			name:         "Internal error",
			trainerId:    1,
			inputBody:    `{"title":"test", "description":"test", "user_id":2}`,
			inputWorkout: entity.Workout{Title: "test", Description: "test", UserId: 2},
			mockBehaviour: func(r *mock_service.MockUser, inputWorkout entity.Workout, trainerId int64) {
				_ = formatTrainerWorkout(&inputWorkout, trainerId)
				r.EXPECT().CreateWorkoutAsTrainer(&inputWorkout).Return(int64(-1), errors.New("internal error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"internal error"}`,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockUser(c)
			test.mockBehaviour(repo, test.inputWorkout, test.trainerId)

			services := &service.Services{User: repo}
			handler := &Handler{services: services}

			r := gin.New()
			r.POST("/workout", handler.createTrainerWorkout)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/workout",
				bytes.NewBufferString(test.inputBody))
			ctx, _ := gin.CreateTestContext(w)
			ctx.Set(userIdCtx, test.trainerId)

			req = req.WithContext(ctx)
			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_initPartnershipWithUser(t *testing.T) {
	type mockBehaviour func(r *mock_service.MockUser, trainerId, userId int64)

	table := []struct {
		name                 string
		trainerId            int64
		userId               int64
		mockBehaviour        mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			trainerId: 1,
			userId:    2,
			mockBehaviour: func(r *mock_service.MockUser, trainerId, userId int64) {
				r.EXPECT().InitPartnershipWithUser(trainerId, userId).Return(int64(1), nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"partnership_id":1}`,
		},
		{
			name:                 "Invalid userId",
			trainerId:            1,
			userId:               -2,
			mockBehaviour:        func(r *mock_service.MockUser, trainerId, userId int64) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid id parameter"}`,
		},
		{
			name:                 "Invalid id",
			trainerId:            -1,
			userId:               2,
			mockBehaviour:        func(r *mock_service.MockUser, trainerId, userId int64) {},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"invalid id"}`,
		},
		{
			name:      "Bad userId",
			trainerId: 1,
			userId:    2,
			mockBehaviour: func(r *mock_service.MockUser, trainerId, userId int64) {
				r.EXPECT().InitPartnershipWithUser(trainerId, userId).Return(int64(-1), errors.New("bad userId"))
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"bad userId"}`,
		},
		{
			name:      "Internal error",
			trainerId: 1,
			userId:    2,
			mockBehaviour: func(r *mock_service.MockUser, trainerId, userId int64) {
				r.EXPECT().InitPartnershipWithUser(trainerId, userId).Return(int64(0), errors.New("internal error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"internal error"}`,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_service.NewMockUser(c)
			test.mockBehaviour(user, test.trainerId, test.userId)

			services := &service.Services{User: user}
			handler := &Handler{services: services}

			r := gin.New()
			r.POST("/user/:id", handler.initPartnershipWithUser)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/user/%d", test.userId), nil)
			ctx, _ := gin.CreateTestContext(w)
			ctx.Set(userIdCtx, test.trainerId)
			r.ServeHTTP(w, req.WithContext(ctx))

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_endPartnershipWithUser(t *testing.T) {
	type mockBehaviour func(r *mock_service.MockUser, trainerId, userId int64)

	table := []struct {
		name                 string
		trainerId            int64
		userId               int64
		mockBehaviour        mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			trainerId: 1,
			userId:    2,
			mockBehaviour: func(r *mock_service.MockUser, trainerId, userId int64) {
				r.EXPECT().EndPartnershipWithUser(trainerId, userId).Return(int64(1), nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"partnership_id":1}`,
		},
		{
			name:                 "Invalid userId",
			trainerId:            1,
			userId:               -2,
			mockBehaviour:        func(r *mock_service.MockUser, trainerId, userId int64) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid id parameter"}`,
		},
		{
			name:                 "Invalid id",
			trainerId:            -1,
			userId:               2,
			mockBehaviour:        func(r *mock_service.MockUser, trainerId, userId int64) {},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"invalid id"}`,
		},
		{
			name:      "Bad userId",
			trainerId: 1,
			userId:    2,
			mockBehaviour: func(r *mock_service.MockUser, trainerId, userId int64) {
				r.EXPECT().EndPartnershipWithUser(trainerId, userId).Return(int64(-1), errors.New("bad userId"))
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"bad userId"}`,
		},
		{
			name:      "Internal error",
			trainerId: 1,
			userId:    2,
			mockBehaviour: func(r *mock_service.MockUser, trainerId, userId int64) {
				r.EXPECT().EndPartnershipWithUser(trainerId, userId).Return(int64(0), errors.New("internal error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"internal error"}`,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_service.NewMockUser(c)
			test.mockBehaviour(user, test.trainerId, test.userId)

			services := &service.Services{User: user}
			handler := &Handler{services: services}

			r := gin.New()
			r.PUT("/user/:id", handler.endPartnershipWithUser)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/user/%d", test.userId), nil)
			ctx, _ := gin.CreateTestContext(w)
			ctx.Set(userIdCtx, test.trainerId)
			r.ServeHTTP(w, req.WithContext(ctx))

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_acceptRequest(t *testing.T) {
	type mockBehaviour func(r *mock_service.MockUser, trainerId, requestId int64)

	table := []struct {
		name                 string
		trainerId            int64
		requestId            int64
		mockBehaviour        mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			trainerId: 1,
			requestId: 2,
			mockBehaviour: func(r *mock_service.MockUser, trainerId, requestId int64) {
				r.EXPECT().AcceptRequest(trainerId, requestId).Return(int64(1), nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"partnership_id":1}`,
		},
		{
			name:                 "Invalid requestId",
			trainerId:            1,
			requestId:            -2,
			mockBehaviour:        func(r *mock_service.MockUser, trainerId, requestId int64) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid id parameter"}`,
		},
		{
			name:                 "Invalid trainerId",
			trainerId:            -1,
			requestId:            2,
			mockBehaviour:        func(r *mock_service.MockUser, trainerId, requestId int64) {},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"invalid id"}`,
		},
		{
			name:      "No request to accept",
			trainerId: 1,
			requestId: 2,
			mockBehaviour: func(r *mock_service.MockUser, trainerId, requestId int64) {
				r.EXPECT().AcceptRequest(trainerId, requestId).Return(int64(0), errors.New("no request"))
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"no request"}`,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_service.NewMockUser(c)
			test.mockBehaviour(user, test.trainerId, test.requestId)

			services := &service.Services{User: user}
			handler := &Handler{services: services}
			r := gin.New()
			r.PUT("/request/:id", handler.acceptRequest)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/request/%d", test.requestId), nil)
			ctx, _ := gin.CreateTestContext(w)
			ctx.Set(userIdCtx, test.trainerId)
			r.ServeHTTP(w, req.WithContext(ctx))

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_denyRequest(t *testing.T) {
	type mockBehaviour func(r *mock_service.MockUser, trainerId, requestId int64)

	table := []struct {
		name               string
		trainerId          int64
		requestId          int64
		mockBehaviour      mockBehaviour
		expectedStatusCode int
	}{
		{
			name:      "Ok",
			trainerId: 1,
			requestId: 2,
			mockBehaviour: func(r *mock_service.MockUser, trainerId, requestId int64) {
				r.EXPECT().DenyRequest(trainerId, requestId).Return(nil)
			},
			expectedStatusCode: 200,
		},
		{
			name:               "Invalid requestId",
			trainerId:          1,
			requestId:          -2,
			mockBehaviour:      func(r *mock_service.MockUser, trainerId, requestId int64) {},
			expectedStatusCode: 400,
		},
		{
			name:               "Invalid trainerId",
			trainerId:          -1,
			requestId:          2,
			mockBehaviour:      func(r *mock_service.MockUser, trainerId, requestId int64) {},
			expectedStatusCode: 500,
		},
		{
			name:      "No request to deny",
			trainerId: 1,
			requestId: 2,
			mockBehaviour: func(r *mock_service.MockUser, trainerId, requestId int64) {
				r.EXPECT().DenyRequest(trainerId, requestId).Return(errors.New("no request"))
			},
			expectedStatusCode: 400,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_service.NewMockUser(c)
			test.mockBehaviour(user, test.trainerId, test.requestId)

			services := &service.Services{User: user}
			handler := &Handler{services: services}
			r := gin.New()
			r.PUT("/request/:id", handler.denyRequest)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/request/%d", test.requestId), nil)
			ctx, _ := gin.CreateTestContext(w)
			ctx.Set(userIdCtx, test.trainerId)
			r.ServeHTTP(w, req.WithContext(ctx))

			assert.Equal(t, w.Code, test.expectedStatusCode)
		})
	}
}
