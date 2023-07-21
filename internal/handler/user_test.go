package handler

import (
	"Fitness_REST_API/internal/entity"
	"Fitness_REST_API/internal/service"
	mock_service "Fitness_REST_API/internal/service/mocks"
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	"net/http/httptest"
	"testing"
)

func TestHandler_getUserInfo(t *testing.T) {
	type mockBehaviour func(r *mock_service.MockUser, userId int64)

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
			mockBehaviour: func(r *mock_service.MockUser, userId int64) {
				r.EXPECT().GetUser(userId).Return(&entity.User{
					Id:           1,
					Email:        "test",
					PasswordHash: "test",
					Role:         "test",
					Name:         "test",
					Surname:      "test",
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1,"email":"test","password":"hidden","role":"test","name":"test","surname":"test","created_at":"0001-01-01T00:00:00Z"}`,
		},
		{
			name:                 "Invalid id",
			userId:               -1,
			mockBehaviour:        func(r *mock_service.MockUser, userId int64) {},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"invalid id"}`,
		},
		{
			name:   "Internal error",
			userId: 1,
			mockBehaviour: func(r *mock_service.MockUser, userId int64) {
				r.EXPECT().GetUser(userId).Return(nil, errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"something went wrong"}`,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockUser(c)
			test.mockBehaviour(repo, test.userId)

			services := &service.Services{User: repo}
			handler := &Handler{services: services}

			r := gin.New()
			r.GET("/user", handler.getUserInfo)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/user", nil)
			ctx, _ := gin.CreateTestContext(w)
			ctx.Set(userIdCtx, test.userId)

			req = req.WithContext(ctx)
			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_createWorkout(t *testing.T) {
	type mockBehaviour func(r *mock_service.MockUser, inputWorkout entity.Workout)

	table := []struct {
		name                 string
		userId               int64
		inputBody            string
		inputWorkout         entity.Workout
		mockBehaviour        mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:         "Ok",
			userId:       1,
			inputBody:    `{"title":"test", "description":"test"}`,
			inputWorkout: entity.Workout{Title: "test", Description: "test", UserId: 1},
			mockBehaviour: func(r *mock_service.MockUser, inputWorkout entity.Workout) {
				r.EXPECT().CreateWorkoutAsUser(&inputWorkout).Return(int64(1), nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:                 "Invalid User Id in workout",
			userId:               1,
			inputBody:            `{"title":"test", "description":"test", "user_id":2}`,
			inputWorkout:         entity.Workout{Title: "test", Description: "test", UserId: 2},
			mockBehaviour:        func(r *mock_service.MockUser, inputWorkout entity.Workout) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"user_id from token and user_id from workout must match"}`,
		},
		{
			name:                 "Workout without a title (including empty workout)",
			userId:               1,
			inputBody:            `{}`,
			inputWorkout:         entity.Workout{},
			mockBehaviour:        func(r *mock_service.MockUser, inputWorkout entity.Workout) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"Key: 'Workout.Title' Error:Field validation for 'Title' failed on the 'required' tag"}`,
		},
		{
			name:                 "Empty title",
			userId:               1,
			inputBody:            `{"title":"", "description":"test"}`,
			inputWorkout:         entity.Workout{},
			mockBehaviour:        func(r *mock_service.MockUser, inputWorkout entity.Workout) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"Key: 'Workout.Title' Error:Field validation for 'Title' failed on the 'required' tag"}`,
		},
		{
			name:         "Internal error",
			userId:       1,
			inputBody:    `{"title":"test", "description":"test"}`,
			inputWorkout: entity.Workout{Title: "test", Description: "test", UserId: 1},
			mockBehaviour: func(r *mock_service.MockUser, inputWorkout entity.Workout) {
				r.EXPECT().CreateWorkoutAsUser(&inputWorkout).Return(int64(-1), errors.New("internal error"))
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
			test.mockBehaviour(repo, test.inputWorkout)

			services := &service.Services{User: repo}
			handler := &Handler{services: services}

			r := gin.New()
			r.POST("/workout", handler.createWorkout)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/workout",
				bytes.NewBufferString(test.inputBody))
			ctx, _ := gin.CreateTestContext(w)
			ctx.Set(userIdCtx, test.userId)

			req = req.WithContext(ctx)
			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_getUserWorkouts(t *testing.T) {
	type mockBehaviour func(r *mock_service.MockUser, userId int64)

	table := []struct {
		name                 string
		userId               int64
		workouts             []*entity.Workout
		mockBehaviour        mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:   "Ok",
			userId: 1,
			mockBehaviour: func(r *mock_service.MockUser, userId int64) {
				r.EXPECT().GetAllUserWorkouts(userId).Return(
					[]*entity.Workout{
						{Title: "test1", Description: "test1"},
						{Title: "test2"},
					}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"id":0,"title":"test1","user_id":0,"trainer_id":{"Int64":0,"Valid":false},"description":"test1","date":"0001-01-01T00:00:00Z"},{"id":0,"title":"test2","user_id":0,"trainer_id":{"Int64":0,"Valid":false},"description":"","date":"0001-01-01T00:00:00Z"}]`,
		},
		{
			name:   "Empty workout",
			userId: 1,
			mockBehaviour: func(r *mock_service.MockUser, userId int64) {
				r.EXPECT().GetAllUserWorkouts(userId).Return(
					[]*entity.Workout{}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[]`,
		},
		{
			name:                 "Invalid id",
			userId:               -1,
			mockBehaviour:        func(r *mock_service.MockUser, userId int64) {},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"invalid id"}`,
		},
		{
			name:   "Internal error",
			userId: 1,
			mockBehaviour: func(r *mock_service.MockUser, userId int64) {
				r.EXPECT().GetAllUserWorkouts(userId).Return(nil, errors.New("internal error"))
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
			test.mockBehaviour(repo, test.userId)

			services := &service.Services{User: repo}
			handler := &Handler{services: services}

			router := gin.New()
			router.GET("/workout", handler.getUserWorkouts)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/workout", nil)
			ctx, _ := gin.CreateTestContext(w)
			ctx.Set(userIdCtx, test.userId)
			req = req.WithContext(ctx)

			router.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_getWorkoutById(t *testing.T) {
	type mockBehaviour func(r *mock_service.MockUser, workoutId, userId int64)

	table := []struct {
		name                 string
		userId               int64
		workoutId            int64
		mockBehaviour        mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			userId:    1,
			workoutId: 1,
			mockBehaviour: func(r *mock_service.MockUser, workoutId, userId int64) {
				r.EXPECT().GetWorkoutById(workoutId, userId).Return(&entity.Workout{Title: "test", Description: "test"}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":0,"title":"test","user_id":0,"trainer_id":{"Int64":0,"Valid":false},"description":"test","date":"0001-01-01T00:00:00Z"}`,
		},
		{
			name:                 "Invalid WorkoutId",
			userId:               1,
			workoutId:            -1,
			mockBehaviour:        func(r *mock_service.MockUser, workoutId, userId int64) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid id param"}`,
		},
		{
			name:      "No workout was found or no access to workout",
			userId:    1,
			workoutId: 100,
			mockBehaviour: func(r *mock_service.MockUser, workoutId, userId int64) {
				r.EXPECT().GetWorkoutById(workoutId, userId).Return(nil, errors.New("invalid id or no access to workout"))
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid id or no access to workout"}`,
		},
		{
			name:                 "Invalid UserId",
			userId:               -1,
			workoutId:            1,
			mockBehaviour:        func(r *mock_service.MockUser, workoutId, userId int64) {},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"invalid id"}`,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockUser(c)
			test.mockBehaviour(repo, test.workoutId, test.userId)

			services := &service.Services{User: repo}
			handler := &Handler{services: services}

			router := gin.New()
			router.GET("/workout/:id", handler.getWorkoutByID)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/workout/%d", test.workoutId), nil)
			ctx, _ := gin.CreateTestContext(w)
			ctx.Set(userIdCtx, test.userId)
			req = req.WithContext(ctx)

			router.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_updateWorkout(t *testing.T) {
	type mockBehaviour func(r *mock_service.MockUser, workoutId, userId int64, input entity.UpdateWorkout)

	table := []struct {
		name                 string
		userId               int64
		workoutId            int64
		inputBody            string
		updateWorkout        entity.UpdateWorkout
		mockBehaviour        mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:          "Ok",
			userId:        1,
			workoutId:     1,
			inputBody:     `{"title":"newTitle", "description":"newDesc"}`,
			updateWorkout: entity.UpdateWorkout{Title: "newTitle", Description: "newDesc"},
			mockBehaviour: func(r *mock_service.MockUser, workoutId, userId int64, input entity.UpdateWorkout) {
				r.EXPECT().UpdateWorkout(workoutId, userId, &input).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:                 "Invalid userId",
			userId:               -1,
			workoutId:            1,
			inputBody:            `{"title":"newTitle", "description":"newDesc"}`,
			updateWorkout:        entity.UpdateWorkout{Title: "newTitle", Description: "newDesc"},
			mockBehaviour:        func(r *mock_service.MockUser, workoutId, userId int64, input entity.UpdateWorkout) {},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"invalid id"}`,
		},
		{
			name:          "Invalid workoutId",
			userId:        1,
			workoutId:     -1,
			inputBody:     `{"title":"newTitle", "description":"newDesc"}`,
			updateWorkout: entity.UpdateWorkout{Title: "newTitle", Description: "newDesc"},
			mockBehaviour: func(r *mock_service.MockUser, workoutId, userId int64, input entity.UpdateWorkout) {
				r.EXPECT().UpdateWorkout(workoutId, userId, &input).Return(errors.New("invalid workoutId"))
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid workoutId"}`,
		},
		{
			name:          "Empty workout",
			userId:        1,
			workoutId:     1,
			inputBody:     `{}`,
			updateWorkout: entity.UpdateWorkout{},
			mockBehaviour: func(r *mock_service.MockUser, workoutId, userId int64, input entity.UpdateWorkout) {
				r.EXPECT().UpdateWorkout(workoutId, userId, &input).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}`,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockUser(c)
			test.mockBehaviour(repo, test.workoutId, test.userId, test.updateWorkout)

			services := &service.Services{User: repo}
			handler := &Handler{services: services}

			router := gin.New()
			router.PUT("/workout/:id", handler.updateWorkout)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", fmt.Sprintf("/workout/%d", test.workoutId),
				bytes.NewBufferString(test.inputBody))
			ctx, _ := gin.CreateTestContext(w)
			ctx.Set(userIdCtx, test.userId)
			req = req.WithContext(ctx)

			router.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_deleteWorkout(t *testing.T) {
	type mockBehaviour func(r *mock_service.MockUser, workoutId, userId int64)

	table := []struct {
		name               string
		userId             int64
		workoutId          int64
		mockBehaviour      mockBehaviour
		expectedStatusCode int
	}{
		{
			name:      "Ok",
			userId:    1,
			workoutId: 1,
			mockBehaviour: func(r *mock_service.MockUser, workoutId, userId int64) {
				r.EXPECT().DeleteWorkout(workoutId, userId).Return(nil)
			},
			expectedStatusCode: 200,
		},
		{
			name:               "Invalid userId",
			userId:             -1,
			workoutId:          1,
			mockBehaviour:      func(r *mock_service.MockUser, workoutId, userId int64) {},
			expectedStatusCode: 500,
		},
		{
			name:      "Invalid workoutId",
			userId:    1,
			workoutId: -1,
			mockBehaviour: func(r *mock_service.MockUser, workoutId, userId int64) {
				r.EXPECT().DeleteWorkout(workoutId, userId).Return(errors.New("invalid workoutId"))
			},
			expectedStatusCode: 400,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockUser(c)
			test.mockBehaviour(repo, test.workoutId, test.userId)

			services := &service.Services{User: repo}
			handler := &Handler{services: services}

			router := gin.New()
			router.DELETE("/workout/:id", handler.deleteWorkout)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/workout/%d", test.workoutId), nil)
			ctx, _ := gin.CreateTestContext(w)
			ctx.Set(userIdCtx, test.userId)
			req = req.WithContext(ctx)

			router.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
		})
	}
}

func TestHandler_getAllTrainers(t *testing.T) {
	type mockBehaviour func(r *mock_service.MockUser)
	table := []struct {
		name                 string
		mockBehaviour        mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			mockBehaviour: func(r *mock_service.MockUser) {
				r.EXPECT().GetAllTrainers().Return([]*entity.User{{Email: "test"}, {Email: "test2"}}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"id":0,"email":"test","password":"","role":"","name":"","surname":"","created_at":"0001-01-01T00:00:00Z"},{"id":0,"email":"test2","password":"","role":"","name":"","surname":"","created_at":"0001-01-01T00:00:00Z"}]`,
		},
		{
			name: "Internal error",
			mockBehaviour: func(r *mock_service.MockUser) {
				r.EXPECT().GetAllTrainers().Return(nil, errors.New("internal error"))
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
			test.mockBehaviour(repo)

			services := &service.Services{User: repo}
			handler := &Handler{services: services}

			router := gin.New()
			router.GET("/trainer", handler.getAllTrainers)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/trainer", nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_getTrainerById(t *testing.T) {
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
				r.EXPECT().GetTrainerById(trainerId).Return(&entity.User{Email: "test"}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":0,"email":"test","password":"","role":"","name":"","surname":"","created_at":"0001-01-01T00:00:00Z"}`,
		},
		{
			name:                 "Invalid trainerId",
			trainerId:            -1,
			mockBehaviour:        func(r *mock_service.MockUser, trainerId int64) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid id param"}`,
		},
		{
			name:      "No trainer was found",
			trainerId: 100,
			mockBehaviour: func(r *mock_service.MockUser, trainerId int64) {
				r.EXPECT().GetTrainerById(trainerId).Return(nil, errors.New("no trainer"))
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"no trainer was found on provided id"}`,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockUser(c)
			test.mockBehaviour(repo, test.trainerId)

			services := &service.Services{User: repo}
			handler := &Handler{services: services}

			router := gin.New()
			router.GET("/trainer/:id", handler.getTrainerByID)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/trainer/%d", test.trainerId), nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_getPartnerships(t *testing.T) {
	type mockBehaviour func(r *mock_service.MockUser, userId int64)
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
			mockBehaviour: func(r *mock_service.MockUser, userId int64) {
				r.EXPECT().GetUserPartnerships(userId).Return([]*entity.Partnership{{Id: 1, UserId: 1, TrainerId: 1}, {Id: 2, UserId: 1, TrainerId: 2}}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"Id":1,"UserId":1,"TrainerId":1,"Status":"","CreatedAt":"0001-01-01T00:00:00Z","EndedAt":{"Time":"0001-01-01T00:00:00Z","Valid":false}},{"Id":2,"UserId":1,"TrainerId":2,"Status":"","CreatedAt":"0001-01-01T00:00:00Z","EndedAt":{"Time":"0001-01-01T00:00:00Z","Valid":false}}]`,
		},
		{
			name:   "No partnerships",
			userId: 1,
			mockBehaviour: func(r *mock_service.MockUser, userId int64) {
				r.EXPECT().GetUserPartnerships(userId).Return([]*entity.Partnership{}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[]`,
		},
		{
			name:                 "Invalid userId",
			userId:               -1,
			mockBehaviour:        func(r *mock_service.MockUser, userId int64) {},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"invalid id"}`,
		},
		{
			name:   "Internal error",
			userId: 1,
			mockBehaviour: func(r *mock_service.MockUser, userId int64) {
				r.EXPECT().GetUserPartnerships(userId).Return(nil, errors.New("internal error"))
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
			test.mockBehaviour(repo, test.userId)

			services := &service.Services{User: repo}
			handler := &Handler{services: services}

			router := gin.New()
			router.GET("/partnership", handler.getPartnerships)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/partnership", nil)
			ctx, _ := gin.CreateTestContext(w)
			ctx.Set(userIdCtx, test.userId)

			router.ServeHTTP(w, req.WithContext(ctx))

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_sendRequestToTrainer(t *testing.T) {
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
			userId:    1,
			mockBehaviour: func(r *mock_service.MockUser, trainerId, userId int64) {
				r.EXPECT().SendRequestToTrainer(trainerId, userId).Return(int64(1), nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:                 "Invalid userId",
			trainerId:            1,
			userId:               -1,
			mockBehaviour:        func(r *mock_service.MockUser, trainerId, userId int64) {},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"invalid id"}`,
		},
		{
			name:                 "Invalid trainerId",
			trainerId:            -1,
			userId:               1,
			mockBehaviour:        func(r *mock_service.MockUser, trainerId, userId int64) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid id param"}`,
		},
		{
			name:      "Approved partnership already exists",
			trainerId: 1,
			userId:    1,
			mockBehaviour: func(r *mock_service.MockUser, trainerId, userId int64) {
				r.EXPECT().SendRequestToTrainer(trainerId, userId).Return(int64(-1), errors.New("already exists"))
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"already exists"}`,
		},
		{
			name:      "Internal error",
			trainerId: 1,
			userId:    1,
			mockBehaviour: func(r *mock_service.MockUser, trainerId, userId int64) {
				r.EXPECT().SendRequestToTrainer(trainerId, userId).Return(int64(0), errors.New("internal error"))
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
			test.mockBehaviour(repo, test.trainerId, test.userId)

			services := &service.Services{User: repo}
			handler := &Handler{services: services}

			router := gin.New()
			router.POST("/partnership/trainer/:id", handler.sendRequestToTrainer)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", fmt.Sprintf("/partnership/trainer/%d", test.trainerId), nil)
			ctx, _ := gin.CreateTestContext(w)
			ctx.Set(userIdCtx, test.userId)

			router.ServeHTTP(w, req.WithContext(ctx))

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_endPartnershipWithTrainer(t *testing.T) {
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
			userId:    1,
			mockBehaviour: func(r *mock_service.MockUser, trainerId, userId int64) {
				r.EXPECT().EndPartnershipWithTrainer(trainerId, userId).Return(int64(1), nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:                 "Invalid userId",
			trainerId:            1,
			userId:               -1,
			mockBehaviour:        func(r *mock_service.MockUser, trainerId, userId int64) {},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"invalid id"}`,
		},
		{
			name:                 "Invalid trainerId",
			trainerId:            -1,
			userId:               1,
			mockBehaviour:        func(r *mock_service.MockUser, trainerId, userId int64) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid id param"}`,
		},
		{
			name:      "No partnership to end",
			trainerId: 1,
			userId:    1,
			mockBehaviour: func(r *mock_service.MockUser, trainerId, userId int64) {
				r.EXPECT().EndPartnershipWithTrainer(trainerId, userId).Return(int64(-1), errors.New("no partnership to end"))
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"no partnership to end"}`,
		},
		{
			name:      "Internal error",
			trainerId: 1,
			userId:    1,
			mockBehaviour: func(r *mock_service.MockUser, trainerId, userId int64) {
				r.EXPECT().EndPartnershipWithTrainer(trainerId, userId).Return(int64(0), errors.New("internal error"))
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
			test.mockBehaviour(repo, test.trainerId, test.userId)

			services := &service.Services{User: repo}
			handler := &Handler{services: services}

			router := gin.New()
			router.PUT("/partnership/trainer/:id", handler.endPartnershipWithTrainer)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", fmt.Sprintf("/partnership/trainer/%d", test.trainerId), nil)
			ctx, _ := gin.CreateTestContext(w)
			ctx.Set(userIdCtx, test.userId)

			router.ServeHTTP(w, req.WithContext(ctx))

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
