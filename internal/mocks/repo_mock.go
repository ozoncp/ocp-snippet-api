// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ozoncp/ocp-snippet-api/internal/repo (interfaces: Repo)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/ozoncp/ocp-snippet-api/internal/models"
)

// MockRepo is a mock of Repo interface.
type MockRepo struct {
	ctrl     *gomock.Controller
	recorder *MockRepoMockRecorder
}

// MockRepoMockRecorder is the mock recorder for MockRepo.
type MockRepoMockRecorder struct {
	mock *MockRepo
}

// NewMockRepo creates a new mock instance.
func NewMockRepo(ctrl *gomock.Controller) *MockRepo {
	mock := &MockRepo{ctrl: ctrl}
	mock.recorder = &MockRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepo) EXPECT() *MockRepoMockRecorder {
	return m.recorder
}

// AddSnippets mocks base method.
func (m *MockRepo) AddSnippets(arg0 context.Context, arg1 []models.Snippet) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddSnippets", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddSnippets indicates an expected call of AddSnippets.
func (mr *MockRepoMockRecorder) AddSnippets(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSnippets", reflect.TypeOf((*MockRepo)(nil).AddSnippets), arg0, arg1)
}

// DescribeSnippet mocks base method.
func (m *MockRepo) DescribeSnippet(arg0 context.Context, arg1 uint64) (*models.Snippet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DescribeSnippet", arg0, arg1)
	ret0, _ := ret[0].(*models.Snippet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeSnippet indicates an expected call of DescribeSnippet.
func (mr *MockRepoMockRecorder) DescribeSnippet(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeSnippet", reflect.TypeOf((*MockRepo)(nil).DescribeSnippet), arg0, arg1)
}

// ListSnippets mocks base method.
func (m *MockRepo) ListSnippets(arg0 context.Context, arg1, arg2 uint64) ([]models.Snippet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSnippets", arg0, arg1, arg2)
	ret0, _ := ret[0].([]models.Snippet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSnippets indicates an expected call of ListSnippets.
func (mr *MockRepoMockRecorder) ListSnippets(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSnippets", reflect.TypeOf((*MockRepo)(nil).ListSnippets), arg0, arg1, arg2)
}

// RemoveSnippet mocks base method.
func (m *MockRepo) RemoveSnippet(arg0 context.Context, arg1 uint64) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveSnippet", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RemoveSnippet indicates an expected call of RemoveSnippet.
func (mr *MockRepoMockRecorder) RemoveSnippet(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveSnippet", reflect.TypeOf((*MockRepo)(nil).RemoveSnippet), arg0, arg1)
}

// UpdateSnippet mocks base method.
func (m *MockRepo) UpdateSnippet(arg0 context.Context, arg1 models.Snippet) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSnippet", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateSnippet indicates an expected call of UpdateSnippet.
func (mr *MockRepoMockRecorder) UpdateSnippet(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSnippet", reflect.TypeOf((*MockRepo)(nil).UpdateSnippet), arg0, arg1)
}
