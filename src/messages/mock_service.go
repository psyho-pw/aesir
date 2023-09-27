// Code generated by mockery v2.33.3. DO NOT EDIT.

package messages

import (
	mock "github.com/stretchr/testify/mock"
	gorm "gorm.io/gorm"
)

// MockService is an autogenerated mock type for the Service type
type MockService struct {
	mock.Mock
}

// FindMany provides a mock function with given fields:
func (_m *MockService) FindMany() ([]Message, error) {
	ret := _m.Called()

	var r0 []Message
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]Message, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []Message); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]Message)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateTimestampsByChannelIds provides a mock function with given fields: channelIds, threshold
func (_m *MockService) UpdateTimestampsByChannelIds(channelIds []int, threshold int) error {
	ret := _m.Called(channelIds, threshold)

	var r0 error
	if rf, ok := ret.Get(0).(func([]int, int) error); ok {
		r0 = rf(channelIds, threshold)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WithTx provides a mock function with given fields: tx
func (_m *MockService) WithTx(tx *gorm.DB) Service {
	ret := _m.Called(tx)

	var r0 Service
	if rf, ok := ret.Get(0).(func(*gorm.DB) Service); ok {
		r0 = rf(tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Service)
		}
	}

	return r0
}

// NewMockService creates a new instance of MockService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockService {
	mock := &MockService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
