package handler

import (
	"Fitness_REST_API/internal/entity"
	"Fitness_REST_API/internal/service"
	mock_service "Fitness_REST_API/internal/service/mocks"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_adminIdentity(t *testing.T) {
	type mockBehavior func(r *mock_service.MockAdmin, token string)

	table := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *mock_service.MockAdmin, token string) {
				r.EXPECT().ParseToken(token).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "ok",
		},
		{
			name:                 "Invalid header name",
			headerName:           "Authorizationn",
			headerValue:          "Bearer token",
			token:                "token",
			mockBehavior:         func(r *mock_service.MockAdmin, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"empty auth header"}`,
		},
		{
			name:                 "Empty header name",
			headerName:           "",
			headerValue:          "Bearer token",
			token:                "token",
			mockBehavior:         func(r *mock_service.MockAdmin, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"empty auth header"}`,
		},
		{
			name:                 "Empty token",
			headerName:           "Authorization",
			headerValue:          "Bearer ",
			token:                "",
			mockBehavior:         func(r *mock_service.MockAdmin, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"invalid auth header"}`,
		},
		{
			name:                 "Empty token",
			headerName:           "Authorization",
			headerValue:          "Bearer ",
			token:                "token",
			mockBehavior:         func(r *mock_service.MockAdmin, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"invalid auth header"}`,
		},
		{
			name:        "Parse error",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *mock_service.MockAdmin, token string) {
				r.EXPECT().ParseToken(token).Return(errors.New("some parsing error"))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"some parsing error"}`,
		},
	}

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockAdmin(c)
			test.mockBehavior(repo, test.token)

			services := &service.Services{Admin: repo}
			handler := &Handler{services: services}

			r := gin.New()
			r.GET("/identity", handler.adminIdentity)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/identity", nil)
			req.Header.Set(test.headerName, test.headerValue)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_userIdentity(t *testing.T) {
	type mockBehavior func(r *mock_service.MockUser, token string)

	table := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "User Ok",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *mock_service.MockUser, token string) {
				r.EXPECT().ParseToken(token).Return(int64(1), entity.UserRole, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "1",
		},
		{
			name:        "Trainer Ok",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *mock_service.MockUser, token string) {
				r.EXPECT().ParseToken(token).Return(int64(2), entity.TrainerRole, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "2",
		},
		{
			name:                 "Empty header name",
			headerName:           "",
			headerValue:          "Bearer token",
			token:                "token",
			mockBehavior:         func(r *mock_service.MockUser, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"empty auth header"}`,
		},
		{
			name:                 "Empty token",
			headerName:           "Authorization",
			headerValue:          "Bearer ",
			token:                "",
			mockBehavior:         func(r *mock_service.MockUser, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"invalid auth header"}`,
		},
		{
			name:                 "Empty token 2",
			headerName:           "Authorization",
			headerValue:          "Bearer ",
			token:                "token",
			mockBehavior:         func(r *mock_service.MockUser, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"invalid auth header"}`,
		},
		{
			name:        "Parse error",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *mock_service.MockUser, token string) {
				var empty entity.Role
				r.EXPECT().ParseToken(token).Return(int64(-1), empty, errors.New("some parsing error"))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"some parsing error"}`,
		},
	}

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockUser(c)
			test.mockBehavior(repo, test.token)

			services := &service.Services{User: repo}
			handler := &Handler{services: services}

			r := gin.New()
			r.GET("/identity", handler.userIdentity, func(c *gin.Context) {
				id, _ := c.Get(userIdCtx)
				c.String(200, "%d", id)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/identity", nil)
			req.Header.Set(test.headerName, test.headerValue)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_trainerIdentity(t *testing.T) {
	type mockBehavior func(r *mock_service.MockUser, token string)

	table := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "User Ok",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *mock_service.MockUser, token string) {
				r.EXPECT().ParseToken(token).Return(int64(1), entity.TrainerRole, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "1",
		},
		{
			name:        "Trainer Ok",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *mock_service.MockUser, token string) {
				r.EXPECT().ParseToken(token).Return(int64(2), entity.TrainerRole, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "2",
		},
		{
			name:                 "Empty header name",
			headerName:           "",
			headerValue:          "Bearer token",
			token:                "token",
			mockBehavior:         func(r *mock_service.MockUser, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"empty auth header"}`,
		},
		{
			name:                 "Empty token",
			headerName:           "Authorization",
			headerValue:          "Bearer ",
			token:                "",
			mockBehavior:         func(r *mock_service.MockUser, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"invalid auth header"}`,
		},
		{
			name:                 "Empty token 2",
			headerName:           "Authorization",
			headerValue:          "Bearer ",
			token:                "token",
			mockBehavior:         func(r *mock_service.MockUser, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"invalid auth header"}`,
		},
		{
			name:        "Parse error",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *mock_service.MockUser, token string) {
				var empty entity.Role
				r.EXPECT().ParseToken(token).Return(int64(-1), empty, errors.New("some parsing error"))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"some parsing error"}`,
		},
	}

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockUser(c)
			test.mockBehavior(repo, test.token)

			services := &service.Services{User: repo}
			handler := &Handler{services: services}

			r := gin.New()
			r.GET("/identity", handler.trainerIdentity, func(c *gin.Context) {
				id, _ := c.Get(userIdCtx)
				c.String(200, "%d", id)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/identity", nil)
			req.Header.Set(test.headerName, test.headerValue)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
