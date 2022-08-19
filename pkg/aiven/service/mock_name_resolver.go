// Code generated by mockery v2.12.1. DO NOT EDIT.

package service

import (
	aiven "github.com/aiven/aiven-go-client"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// MockNameResolver is an autogenerated mock type for the NameResolver type
type MockNameResolver struct {
	mock.Mock
}

// List provides a mock function with given fields: project
func (_m *MockNameResolver) List(project string) ([]*aiven.Service, error) {
	ret := _m.Called(project)

	var r0 []*aiven.Service
	if rf, ok := ret.Get(0).(func(string) []*aiven.Service); ok {
		r0 = rf(project)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*aiven.Service)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(project)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ResolveKafkaServiceName provides a mock function with given fields: project
func (_m *MockNameResolver) ResolveKafkaServiceName(project string) (string, error) {
	ret := _m.Called(project)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(project)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(project)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockNameResolver creates a new instance of MockNameResolver. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockNameResolver(t testing.TB) *MockNameResolver {
	mock := &MockNameResolver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
