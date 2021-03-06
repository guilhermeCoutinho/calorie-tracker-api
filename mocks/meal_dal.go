// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/guilhermeCoutinho/api-studies/dal (interfaces: MealDAL)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	dal "github.com/guilhermeCoutinho/api-studies/dal"
	models "github.com/guilhermeCoutinho/api-studies/models"
	reflect "reflect"
)

// MockMealDAL is a mock of MealDAL interface
type MockMealDAL struct {
	ctrl     *gomock.Controller
	recorder *MockMealDALMockRecorder
}

// MockMealDALMockRecorder is the mock recorder for MockMealDAL
type MockMealDALMockRecorder struct {
	mock *MockMealDAL
}

// NewMockMealDAL creates a new mock instance
func NewMockMealDAL(ctrl *gomock.Controller) *MockMealDAL {
	mock := &MockMealDAL{ctrl: ctrl}
	mock.recorder = &MockMealDALMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMealDAL) EXPECT() *MockMealDALMockRecorder {
	return m.recorder
}

// DeleteMeal mocks base method
func (m *MockMealDAL) DeleteMeal(arg0 context.Context, arg1 uuid.UUID, arg2 *uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMeal", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteMeal indicates an expected call of DeleteMeal
func (mr *MockMealDALMockRecorder) DeleteMeal(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMeal", reflect.TypeOf((*MockMealDAL)(nil).DeleteMeal), arg0, arg1, arg2)
}

// GetMeals mocks base method
func (m *MockMealDAL) GetMeals(arg0 context.Context, arg1, arg2 *uuid.UUID, arg3 *dal.QueryOptions) ([]*models.MealWithLimit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMeals", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]*models.MealWithLimit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMeals indicates an expected call of GetMeals
func (mr *MockMealDALMockRecorder) GetMeals(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMeals", reflect.TypeOf((*MockMealDAL)(nil).GetMeals), arg0, arg1, arg2, arg3)
}

// InsertMeal mocks base method
func (m *MockMealDAL) InsertMeal(arg0 context.Context, arg1 *models.Meal) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertMeal", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertMeal indicates an expected call of InsertMeal
func (mr *MockMealDALMockRecorder) InsertMeal(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertMeal", reflect.TypeOf((*MockMealDAL)(nil).InsertMeal), arg0, arg1)
}

// UpsertMeal mocks base method
func (m *MockMealDAL) UpsertMeal(arg0 context.Context, arg1 *models.Meal) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertMeal", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertMeal indicates an expected call of UpsertMeal
func (mr *MockMealDALMockRecorder) UpsertMeal(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertMeal", reflect.TypeOf((*MockMealDAL)(nil).UpsertMeal), arg0, arg1)
}
