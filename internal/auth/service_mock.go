// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package auth is a generated GoMock package.
package auth

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	model "github.com/vliubezny/gstore/internal/model"
	reflect "reflect"
)

// MockService is a mock of Service interface
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// Register mocks base method
func (m *MockService) Register(ctx context.Context, user model.User, password string) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", ctx, user, password)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register
func (mr *MockServiceMockRecorder) Register(ctx, user, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockService)(nil).Register), ctx, user, password)
}

// Login mocks base method
func (m *MockService) Login(ctx context.Context, email, password string) (TokenPair, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, email, password)
	ret0, _ := ret[0].(TokenPair)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login
func (mr *MockServiceMockRecorder) Login(ctx, email, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockService)(nil).Login), ctx, email, password)
}

// Refresh mocks base method
func (m *MockService) Refresh(ctx context.Context, refreshToken string) (TokenPair, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Refresh", ctx, refreshToken)
	ret0, _ := ret[0].(TokenPair)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Refresh indicates an expected call of Refresh
func (mr *MockServiceMockRecorder) Refresh(ctx, refreshToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refresh", reflect.TypeOf((*MockService)(nil).Refresh), ctx, refreshToken)
}

// Revoke mocks base method
func (m *MockService) Revoke(ctx context.Context, refreshToken string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Revoke", ctx, refreshToken)
	ret0, _ := ret[0].(error)
	return ret0
}

// Revoke indicates an expected call of Revoke
func (mr *MockServiceMockRecorder) Revoke(ctx, refreshToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Revoke", reflect.TypeOf((*MockService)(nil).Revoke), ctx, refreshToken)
}

// ValidateAccessToken mocks base method
func (m *MockService) ValidateAccessToken(token string) (AccessTokenClaims, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateAccessToken", token)
	ret0, _ := ret[0].(AccessTokenClaims)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateAccessToken indicates an expected call of ValidateAccessToken
func (mr *MockServiceMockRecorder) ValidateAccessToken(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateAccessToken", reflect.TypeOf((*MockService)(nil).ValidateAccessToken), token)
}

// UpdateUserPermissions mocks base method
func (m *MockService) UpdateUserPermissions(ctx context.Context, user model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserPermissions", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserPermissions indicates an expected call of UpdateUserPermissions
func (mr *MockServiceMockRecorder) UpdateUserPermissions(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserPermissions", reflect.TypeOf((*MockService)(nil).UpdateUserPermissions), ctx, user)
}
