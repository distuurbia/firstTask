// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/distuurbia/firstTask/internal/model"

	service "github.com/distuurbia/firstTask/internal/service"
)

// UserService is an autogenerated mock type for the UserService type
type UserService struct {
	mock.Mock
}

// Login provides a mock function with given fields: ctx, user
func (_m *UserService) Login(ctx context.Context, user *model.User) (service.TokenPair, error) {
	ret := _m.Called(ctx, user)

	var r0 service.TokenPair
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.User) (service.TokenPair, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.User) service.TokenPair); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(service.TokenPair)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Refresh provides a mock function with given fields: ctx, tokenPair
func (_m *UserService) Refresh(ctx context.Context, tokenPair service.TokenPair) (service.TokenPair, error) {
	ret := _m.Called(ctx, tokenPair)

	var r0 service.TokenPair
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, service.TokenPair) (service.TokenPair, error)); ok {
		return rf(ctx, tokenPair)
	}
	if rf, ok := ret.Get(0).(func(context.Context, service.TokenPair) service.TokenPair); ok {
		r0 = rf(ctx, tokenPair)
	} else {
		r0 = ret.Get(0).(service.TokenPair)
	}

	if rf, ok := ret.Get(1).(func(context.Context, service.TokenPair) error); ok {
		r1 = rf(ctx, tokenPair)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SignUp provides a mock function with given fields: ctx, user
func (_m *UserService) SignUp(ctx context.Context, user *model.User) error {
	ret := _m.Called(ctx, user)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.User) error); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewUserService creates a new instance of UserService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserService(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserService {
	mock := &UserService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
