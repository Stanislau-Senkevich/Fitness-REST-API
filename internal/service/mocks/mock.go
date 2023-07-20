// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	entity "Fitness_REST_API/internal/entity"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAdmin is a mock of Admin interface.
type MockAdmin struct {
	ctrl     *gomock.Controller
	recorder *MockAdminMockRecorder
}

// MockAdminMockRecorder is the mock recorder for MockAdmin.
type MockAdminMockRecorder struct {
	mock *MockAdmin
}

// NewMockAdmin creates a new mock instance.
func NewMockAdmin(ctrl *gomock.Controller) *MockAdmin {
	mock := &MockAdmin{ctrl: ctrl}
	mock.recorder = &MockAdminMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAdmin) EXPECT() *MockAdminMockRecorder {
	return m.recorder
}

// ParseToken mocks base method.
func (m *MockAdmin) ParseToken(token string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseToken", token)
	ret0, _ := ret[0].(error)
	return ret0
}

// ParseToken indicates an expected call of ParseToken.
func (mr *MockAdminMockRecorder) ParseToken(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseToken", reflect.TypeOf((*MockAdmin)(nil).ParseToken), token)
}

// SignIn mocks base method.
func (m *MockAdmin) SignIn(login, passwordHash string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignIn", login, passwordHash)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignIn indicates an expected call of SignIn.
func (mr *MockAdminMockRecorder) SignIn(login, passwordHash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignIn", reflect.TypeOf((*MockAdmin)(nil).SignIn), login, passwordHash)
}

// MockUser is a mock of User interface.
type MockUser struct {
	ctrl     *gomock.Controller
	recorder *MockUserMockRecorder
}

// MockUserMockRecorder is the mock recorder for MockUser.
type MockUserMockRecorder struct {
	mock *MockUser
}

// NewMockUser creates a new mock instance.
func NewMockUser(ctrl *gomock.Controller) *MockUser {
	mock := &MockUser{ctrl: ctrl}
	mock.recorder = &MockUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUser) EXPECT() *MockUserMockRecorder {
	return m.recorder
}

// CreateWorkoutAsUser mocks base method.
func (m *MockUser) CreateWorkoutAsUser(workout *entity.Workout) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateWorkoutAsUser", workout)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateWorkoutAsUser indicates an expected call of CreateWorkoutAsUser.
func (mr *MockUserMockRecorder) CreateWorkoutAsUser(workout interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateWorkoutAsUser", reflect.TypeOf((*MockUser)(nil).CreateWorkoutAsUser), workout)
}

// DeleteWorkout mocks base method.
func (m *MockUser) DeleteWorkout(workoutId, userId int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteWorkout", workoutId, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteWorkout indicates an expected call of DeleteWorkout.
func (mr *MockUserMockRecorder) DeleteWorkout(workoutId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteWorkout", reflect.TypeOf((*MockUser)(nil).DeleteWorkout), workoutId, userId)
}

// EndPartnershipWithTrainer mocks base method.
func (m *MockUser) EndPartnershipWithTrainer(trainerId, userId int64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EndPartnershipWithTrainer", trainerId, userId)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EndPartnershipWithTrainer indicates an expected call of EndPartnershipWithTrainer.
func (mr *MockUserMockRecorder) EndPartnershipWithTrainer(trainerId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EndPartnershipWithTrainer", reflect.TypeOf((*MockUser)(nil).EndPartnershipWithTrainer), trainerId, userId)
}

// GetAllTrainers mocks base method.
func (m *MockUser) GetAllTrainers() ([]*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllTrainers")
	ret0, _ := ret[0].([]*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllTrainers indicates an expected call of GetAllTrainers.
func (mr *MockUserMockRecorder) GetAllTrainers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllTrainers", reflect.TypeOf((*MockUser)(nil).GetAllTrainers))
}

// GetAllUserWorkouts mocks base method.
func (m *MockUser) GetAllUserWorkouts(id int64) ([]*entity.Workout, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllUserWorkouts", id)
	ret0, _ := ret[0].([]*entity.Workout)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllUserWorkouts indicates an expected call of GetAllUserWorkouts.
func (mr *MockUserMockRecorder) GetAllUserWorkouts(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllUserWorkouts", reflect.TypeOf((*MockUser)(nil).GetAllUserWorkouts), id)
}

// GetTrainerById mocks base method.
func (m *MockUser) GetTrainerById(id int64) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTrainerById", id)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTrainerById indicates an expected call of GetTrainerById.
func (mr *MockUserMockRecorder) GetTrainerById(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTrainerById", reflect.TypeOf((*MockUser)(nil).GetTrainerById), id)
}

// GetUser mocks base method.
func (m *MockUser) GetUser(id int64) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", id)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockUserMockRecorder) GetUser(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockUser)(nil).GetUser), id)
}

// GetUserPartnerships mocks base method.
func (m *MockUser) GetUserPartnerships(userId int64) ([]*entity.Partnership, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserPartnerships", userId)
	ret0, _ := ret[0].([]*entity.Partnership)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserPartnerships indicates an expected call of GetUserPartnerships.
func (mr *MockUserMockRecorder) GetUserPartnerships(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserPartnerships", reflect.TypeOf((*MockUser)(nil).GetUserPartnerships), userId)
}

// GetWorkoutById mocks base method.
func (m *MockUser) GetWorkoutById(workoutId, userId int64) (*entity.Workout, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWorkoutById", workoutId, userId)
	ret0, _ := ret[0].(*entity.Workout)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWorkoutById indicates an expected call of GetWorkoutById.
func (mr *MockUserMockRecorder) GetWorkoutById(workoutId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWorkoutById", reflect.TypeOf((*MockUser)(nil).GetWorkoutById), workoutId, userId)
}

// ParseToken mocks base method.
func (m *MockUser) ParseToken(token string) (int64, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseToken", token)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ParseToken indicates an expected call of ParseToken.
func (mr *MockUserMockRecorder) ParseToken(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseToken", reflect.TypeOf((*MockUser)(nil).ParseToken), token)
}

// SendRequestToTrainer mocks base method.
func (m *MockUser) SendRequestToTrainer(trainerId, userId int64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendRequestToTrainer", trainerId, userId)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendRequestToTrainer indicates an expected call of SendRequestToTrainer.
func (mr *MockUserMockRecorder) SendRequestToTrainer(trainerId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendRequestToTrainer", reflect.TypeOf((*MockUser)(nil).SendRequestToTrainer), trainerId, userId)
}

// SignIn mocks base method.
func (m *MockUser) SignIn(email, passwordHash, role string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignIn", email, passwordHash, role)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignIn indicates an expected call of SignIn.
func (mr *MockUserMockRecorder) SignIn(email, passwordHash, role interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignIn", reflect.TypeOf((*MockUser)(nil).SignIn), email, passwordHash, role)
}

// SignUp mocks base method.
func (m *MockUser) SignUp(user *entity.User) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUp", user)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignUp indicates an expected call of SignUp.
func (mr *MockUserMockRecorder) SignUp(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUp", reflect.TypeOf((*MockUser)(nil).SignUp), user)
}

// UpdateWorkout mocks base method.
func (m *MockUser) UpdateWorkout(workoutId, userId int64, update *entity.UpdateWorkout) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateWorkout", workoutId, userId, update)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateWorkout indicates an expected call of UpdateWorkout.
func (mr *MockUserMockRecorder) UpdateWorkout(workoutId, userId, update interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateWorkout", reflect.TypeOf((*MockUser)(nil).UpdateWorkout), workoutId, userId, update)
}