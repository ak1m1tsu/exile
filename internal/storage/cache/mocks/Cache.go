// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	context "context"
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// Cache is an autogenerated mock type for the Cache type
type Cache struct {
	mock.Mock
}

// Del provides a mock function with given fields: ctx, key
func (_m *Cache) Del(ctx context.Context, key string) error {
	ret := _m.Called(ctx, key)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, key
func (_m *Cache) Get(ctx context.Context, key string) ([]byte, bool, error) {
	ret := _m.Called(ctx, key)

	var r0 []byte
	var r1 bool
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]byte, bool, error)); ok {
		return rf(ctx, key)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []byte); ok {
		r0 = rf(ctx, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) bool); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Get(1).(bool)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string) error); ok {
		r2 = rf(ctx, key)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Set provides a mock function with given fields: ctx, key, value, ttl
func (_m *Cache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	ret := _m.Called(ctx, key, value, ttl)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []byte, time.Duration) error); ok {
		r0 = rf(ctx, key, value, ttl)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewCache interface {
	mock.TestingT
	Cleanup(func())
}

// NewCache creates a new instance of Cache. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCache(t mockConstructorTestingTNewCache) *Cache {
	mock := &Cache{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}