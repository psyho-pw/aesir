// Code generated by mockery v2.32.2. DO NOT EDIT.

package channels

import (
	mock "github.com/stretchr/testify/mock"
	gorm "gorm.io/gorm"
)

// MockService is an autogenerated mock type for the Service type
type MockService struct {
	mock.Mock
}

// Create provides a mock function with given fields: channel
func (_m *MockService) Create(channel Channel) (*Channel, error) {
	ret := _m.Called(channel)

	var r0 *Channel
	var r1 error
	if rf, ok := ret.Get(0).(func(Channel) (*Channel, error)); ok {
		return rf(channel)
	}
	if rf, ok := ret.Get(0).(func(Channel) *Channel); ok {
		r0 = rf(channel)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Channel)
		}
	}

	if rf, ok := ret.Get(1).(func(Channel) error); ok {
		r1 = rf(channel)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateMany provides a mock function with given fields: channels
func (_m *MockService) CreateMany(channels []Channel) ([]Channel, error) {
	ret := _m.Called(channels)

	var r0 []Channel
	var r1 error
	if rf, ok := ret.Get(0).(func([]Channel) ([]Channel, error)); ok {
		return rf(channels)
	}
	if rf, ok := ret.Get(0).(func([]Channel) []Channel); ok {
		r0 = rf(channels)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]Channel)
		}
	}

	if rf, ok := ret.Get(1).(func([]Channel) error); ok {
		r1 = rf(channels)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteOneBySlackId provides a mock function with given fields: slackId
func (_m *MockService) DeleteOneBySlackId(slackId string) (*Channel, error) {
	ret := _m.Called(slackId)

	var r0 *Channel
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*Channel, error)); ok {
		return rf(slackId)
	}
	if rf, ok := ret.Get(0).(func(string) *Channel); ok {
		r0 = rf(slackId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Channel)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(slackId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindMany provides a mock function with given fields:
func (_m *MockService) FindMany() ([]Channel, error) {
	ret := _m.Called()

	var r0 []Channel
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]Channel, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []Channel); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]Channel)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindManyWithMessage provides a mock function with given fields:
func (_m *MockService) FindManyWithMessage() ([]Channel, error) {
	ret := _m.Called()

	var r0 []Channel
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]Channel, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []Channel); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]Channel)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOneBySlackId provides a mock function with given fields: slackId
func (_m *MockService) FindOneBySlackId(slackId string) (*Channel, error) {
	ret := _m.Called(slackId)

	var r0 *Channel
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*Channel, error)); ok {
		return rf(slackId)
	}
	if rf, ok := ret.Get(0).(func(string) *Channel); ok {
		r0 = rf(slackId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Channel)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(slackId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateOneBySlackId provides a mock function with given fields: slackId, channel
func (_m *MockService) UpdateOneBySlackId(slackId string, channel Channel) (*Channel, error) {
	ret := _m.Called(slackId, channel)

	var r0 *Channel
	var r1 error
	if rf, ok := ret.Get(0).(func(string, Channel) (*Channel, error)); ok {
		return rf(slackId, channel)
	}
	if rf, ok := ret.Get(0).(func(string, Channel) *Channel); ok {
		r0 = rf(slackId, channel)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Channel)
		}
	}

	if rf, ok := ret.Get(1).(func(string, Channel) error); ok {
		r1 = rf(slackId, channel)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
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
