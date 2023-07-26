package handler

import (
	"Fitness_REST_API/internal/entity"
	"Fitness_REST_API/internal/service"
	mockService "Fitness_REST_API/internal/service/mocks"
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	"net/http/httptest"
	"testing"
)

func TestHandler_adminSignIn(t *testing.T) {
	type mockBehavior func(r *mockService.MockAdmin, signInInput adminSignInInput)

	table := []struct {
		name                 string
		inputBody            string
		signInInput          adminSignInInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			inputBody:   `{"login":"testLogin", "password":"testPassword"}`,
			signInInput: adminSignInInput{Login: "testLogin", Password: "testPassword"},
			mockBehavior: func(r *mockService.MockAdmin, signInInput adminSignInInput) {
				r.EXPECT().SignIn(signInInput.Login, signInInput.Password).Return("token", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"token":"token"}`,
		},
		{
			name:                 "Not Bindable JSON Login",
			inputBody:            `{"login1":"testLogin", "password":"testPassword"}`,
			signInInput:          adminSignInInput{Login: "testLogin", Password: "testPassword"},
			mockBehavior:         func(r *mockService.MockAdmin, signInInput adminSignInInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"Key: 'adminSignInInput.Login' Error:Field validation for 'Login' failed on the 'required' tag"}`,
		},
		{
			name:                 "Not Bindable JSON Password",
			inputBody:            `{"login":"testLogin", "password1":"testPassword"}`,
			signInInput:          adminSignInInput{Login: "testLogin", Password: "testPassword"},
			mockBehavior:         func(r *mockService.MockAdmin, signInInput adminSignInInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"Key: 'adminSignInInput.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`,
		},
		{
			name:                 "Not Bindable JSON Login and Password",
			inputBody:            `{"login1":"testLogin", "password1":"testPassword"}`,
			signInInput:          adminSignInInput{Login: "testLogin", Password: "testPassword"},
			mockBehavior:         func(r *mockService.MockAdmin, signInInput adminSignInInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"Key: 'adminSignInInput.Login' Error:Field validation for 'Login' failed on the 'required' tag\nKey: 'adminSignInInput.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`,
		},
		{
			name:        "Invalid Login or Password",
			inputBody:   `{"login":"testLogin", "password":"testPassword"}`,
			signInInput: adminSignInInput{Login: "testLogin", Password: "testPassword"},
			mockBehavior: func(r *mockService.MockAdmin, signInInput adminSignInInput) {
				r.EXPECT().SignIn(signInInput.Login, signInInput.Password).Return("", errors.New("invalid login or password"))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"invalid login or password"}`,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mockService.NewMockAdmin(c)
			test.mockBehavior(repo, test.signInInput)

			services := &service.Services{Admin: repo}
			handler := &Handler{services: services}

			r := gin.New()
			r.POST("/admin/sign-in", handler.adminSignIn)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/admin/sign-in",
				bytes.NewBufferString(test.inputBody))
			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_signIn(t *testing.T) {
	type mockBehavior func(r *mockService.MockUser, signInInput userSignInInput)

	table := []struct {
		name                 string
		inputBody            string
		signInInput          userSignInInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			inputBody:   `{"email":"testEmail", "password":"testPassword"}`,
			signInInput: userSignInInput{Email: "testEmail", Password: "testPassword"},
			mockBehavior: func(r *mockService.MockUser, signInInput userSignInInput) {
				r.EXPECT().SignIn(signInInput.Email, signInInput.Password, entity.UserRole).Return("token", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"token":"token"}`,
		},
		{
			name:                 "Not Bindable JSON Email",
			inputBody:            `{"email1":"testEmail", "password":"testPassword"}`,
			signInInput:          userSignInInput{Email: "testEmail", Password: "testPassword"},
			mockBehavior:         func(r *mockService.MockUser, signInInput userSignInInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"Key: 'userSignInInput.Email' Error:Field validation for 'Email' failed on the 'required' tag"}`,
		},
		{
			name:                 "Not Bindable JSON Password",
			inputBody:            `{"email":"testEmail", "password1":"testPassword"}`,
			signInInput:          userSignInInput{Email: "testEmail", Password: "testPassword"},
			mockBehavior:         func(r *mockService.MockUser, signInInput userSignInInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"Key: 'userSignInInput.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`,
		},
		{
			name:                 "Not Bindable JSON Email and Password",
			inputBody:            `{"email1":"testEmail", "password1":"testPassword"}`,
			signInInput:          userSignInInput{Email: "testEmail", Password: "testPassword"},
			mockBehavior:         func(r *mockService.MockUser, signInInput userSignInInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"Key: 'userSignInInput.Email' Error:Field validation for 'Email' failed on the 'required' tag\nKey: 'userSignInInput.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`,
		},
		{
			name:        "Invalid Email or Password",
			inputBody:   `{"email":"testEmail", "password":"testPassword"}`,
			signInInput: userSignInInput{Email: "testEmail", Password: "testPassword"},
			mockBehavior: func(r *mockService.MockUser, signInInput userSignInInput) {
				r.EXPECT().SignIn(signInInput.Email, signInInput.Password, entity.UserRole).Return("", errors.New("invalid email or password"))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"invalid email or password"}`,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mockService.NewMockUser(c)
			test.mockBehavior(repo, test.signInInput)

			services := &service.Services{User: repo}
			handler := &Handler{services: services}

			r := gin.New()
			r.POST("/sign-in", handler.signIn)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-in",
				bytes.NewBufferString(test.inputBody))
			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_trainerSignIn(t *testing.T) {
	type mockBehavior func(r *mockService.MockUser, signInInput userSignInInput)

	table := []struct {
		name                 string
		inputBody            string
		signInInput          userSignInInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			inputBody:   `{"email":"testEmail", "password":"testPassword"}`,
			signInInput: userSignInInput{Email: "testEmail", Password: "testPassword"},
			mockBehavior: func(r *mockService.MockUser, signInInput userSignInInput) {
				r.EXPECT().SignIn(signInInput.Email, signInInput.Password, entity.TrainerRole).Return("token", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"token":"token"}`,
		},
		{
			name:                 "Not Bindable JSON Email",
			inputBody:            `{"email1":"testEmail", "password":"testPassword"}`,
			signInInput:          userSignInInput{Email: "testEmail", Password: "testPassword"},
			mockBehavior:         func(r *mockService.MockUser, signInInput userSignInInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"Key: 'userSignInInput.Email' Error:Field validation for 'Email' failed on the 'required' tag"}`,
		},
		{
			name:                 "Not Bindable JSON Password",
			inputBody:            `{"email":"testEmail", "password1":"testPassword"}`,
			signInInput:          userSignInInput{Email: "testEmail", Password: "testPassword"},
			mockBehavior:         func(r *mockService.MockUser, signInInput userSignInInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"Key: 'userSignInInput.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`,
		},
		{
			name:                 "Not Bindable JSON Email and Password",
			inputBody:            `{"email1":"testEmail", "password1":"testPassword"}`,
			signInInput:          userSignInInput{Email: "testEmail", Password: "testPassword"},
			mockBehavior:         func(r *mockService.MockUser, signInInput userSignInInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"Key: 'userSignInInput.Email' Error:Field validation for 'Email' failed on the 'required' tag\nKey: 'userSignInInput.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`,
		},
		{
			name:        "Invalid Email or Password",
			inputBody:   `{"email":"testEmail", "password":"testPassword"}`,
			signInInput: userSignInInput{Email: "testEmail", Password: "testPassword"},
			mockBehavior: func(r *mockService.MockUser, signInInput userSignInInput) {
				r.EXPECT().SignIn(signInInput.Email, signInInput.Password, entity.TrainerRole).Return("", errors.New("invalid email or password"))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"invalid email or password"}`,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mockService.NewMockUser(c)
			test.mockBehavior(repo, test.signInInput)

			services := &service.Services{User: repo}
			handler := &Handler{services: services}

			r := gin.New()
			r.POST("/sign-in", handler.trainerSignIn)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-in",
				bytes.NewBufferString(test.inputBody))
			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(r *mockService.MockUser, inputUser entity.User)

	table := []struct {
		name                 string
		inputBody            string
		inputUser            entity.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"email":"testEmail", "password":"testPassword", "name":"testName", "surname":"testSurname"}`,
			inputUser: entity.User{Email: "testEmail", PasswordHash: "testPassword", Name: "testName", Surname: "testSurname"},
			mockBehavior: func(r *mockService.MockUser, inputUser entity.User) {
				r.EXPECT().SignUp(&inputUser).Return(int64(1), nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:                 "Invalid email JSON",
			inputBody:            `{"emailA":"testEmail", "password":"testPassword", "name":"testName", "surname":"testSurname"}`,
			inputUser:            entity.User{Email: "testEmail", PasswordHash: "testPassword", Name: "testName", Surname: "testSurname"},
			mockBehavior:         func(r *mockService.MockUser, inputUser entity.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"Key: 'User.Email' Error:Field validation for 'Email' failed on the 'required' tag"}`,
		},
		{
			name:                 "Invalid password JSON",
			inputBody:            `{"email":"testEmail", "passwordA":"testPassword", "name":"testName", "surname":"testSurname"}`,
			inputUser:            entity.User{Email: "testEmail", PasswordHash: "testPassword", Name: "testName", Surname: "testSurname"},
			mockBehavior:         func(r *mockService.MockUser, inputUser entity.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"Key: 'User.PasswordHash' Error:Field validation for 'PasswordHash' failed on the 'required' tag"}`,
		},
		{
			name:                 "Invalid name JSON",
			inputBody:            `{"email":"testEmail", "password":"testPassword", "nameA":"testName", "surname":"testSurname"}`,
			inputUser:            entity.User{Email: "testEmail", PasswordHash: "testPassword", Name: "testName", Surname: "testSurname"},
			mockBehavior:         func(r *mockService.MockUser, inputUser entity.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"Key: 'User.Name' Error:Field validation for 'Name' failed on the 'required' tag"}`,
		},
		{
			name:                 "Invalid surname JSON",
			inputBody:            `{"email":"testEmail", "password":"testPassword", "name":"testName", "surnameA":"testSurname"}`,
			inputUser:            entity.User{Email: "testEmail", PasswordHash: "testPassword", Name: "testName", Surname: "testSurname"},
			mockBehavior:         func(r *mockService.MockUser, inputUser entity.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"Key: 'User.Surname' Error:Field validation for 'Surname' failed on the 'required' tag"}`,
		},
		{
			name:      "Internal Server Error",
			inputBody: `{"email":"testEmail", "password":"testPassword", "name":"testName", "surname":"testSurname"}`,
			inputUser: entity.User{Email: "testEmail", PasswordHash: "testPassword", Name: "testName", Surname: "testSurname"},
			mockBehavior: func(r *mockService.MockUser, inputUser entity.User) {
				r.EXPECT().SignUp(&inputUser).Return(int64(0), errors.New("internal error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"internal error"}`,
		},
		{
			name:      "Email has already reserved",
			inputBody: `{"email":"testEmail", "password":"testPassword", "name":"testName", "surname":"testSurname"}`,
			inputUser: entity.User{Email: "testEmail", PasswordHash: "testPassword", Name: "testName", Surname: "testSurname"},
			mockBehavior: func(r *mockService.MockUser, inputUser entity.User) {
				r.EXPECT().SignUp(&inputUser).Return(int64(-1), errors.New("reserved email"))
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"reserved email"}`,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mockService.NewMockUser(c)
			test.mockBehavior(repo, test.inputUser)

			services := &service.Services{User: repo}
			handler := &Handler{services: services}

			r := gin.New()
			r.POST("/sign-up", handler.signUp)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up",
				bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)

		})
	}
}
