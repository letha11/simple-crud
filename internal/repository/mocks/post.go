// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/post.go
//
// Generated by this command:
//
//	mockgen -source=internal/repository/post.go -destination=internal/repository/mocks/post.go
//

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	reflect "reflect"

	models "github.com/simple-crud-go/internal/models"
	gomock "go.uber.org/mock/gomock"
)

// MockPostRepo is a mock of PostRepo interface.
type MockPostRepo struct {
	ctrl     *gomock.Controller
	recorder *MockPostRepoMockRecorder
}

// MockPostRepoMockRecorder is the mock recorder for MockPostRepo.
type MockPostRepoMockRecorder struct {
	mock *MockPostRepo
}

// NewMockPostRepo creates a new mock instance.
func NewMockPostRepo(ctrl *gomock.Controller) *MockPostRepo {
	mock := &MockPostRepo{ctrl: ctrl}
	mock.recorder = &MockPostRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPostRepo) EXPECT() *MockPostRepoMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockPostRepo) Create(post *models.Post) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", post)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockPostRepoMockRecorder) Create(post any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockPostRepo)(nil).Create), post)
}

// GetAll mocks base method.
func (m *MockPostRepo) GetAll() ([]models.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]models.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockPostRepoMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockPostRepo)(nil).GetAll))
}

// GetById mocks base method.
func (m *MockPostRepo) GetById(id int) (*models.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", id)
	ret0, _ := ret[0].(*models.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockPostRepoMockRecorder) GetById(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockPostRepo)(nil).GetById), id)
}

// Update mocks base method.
func (m *MockPostRepo) Update(post *models.Post) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", post)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockPostRepoMockRecorder) Update(post any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockPostRepo)(nil).Update), post)
}
