// Code generated by mockery v2.42.1. DO NOT EDIT.

package clients

import (
	mock "github.com/stretchr/testify/mock"
	gorm "gorm.io/gorm"
)

// MockRepository is an autogenerated mock type for the Repository type
type MockRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: client
func (_m *MockRepository) Create(client Client) (*Client, error) {
	ret := _m.Called(client)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *Client
	var r1 error
	if rf, ok := ret.Get(0).(func(Client) (*Client, error)); ok {
		return rf(client)
	}
	if rf, ok := ret.Get(0).(func(Client) *Client); ok {
		r0 = rf(client)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Client)
		}
	}

	if rf, ok := ret.Get(1).(func(Client) error); ok {
		r1 = rf(client)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteOne provides a mock function with given fields: id
func (_m *MockRepository) DeleteOne(id int) (*Client, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteOne")
	}

	var r0 *Client
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (*Client, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int) *Client); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Client)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindMany provides a mock function with given fields:
func (_m *MockRepository) FindMany() ([]Client, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for FindMany")
	}

	var r0 []Client
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]Client, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []Client); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]Client)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOne provides a mock function with given fields: id
func (_m *MockRepository) FindOne(id int) (*Client, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for FindOne")
	}

	var r0 *Client
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (*Client, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int) *Client); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Client)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WithTx provides a mock function with given fields: tx
func (_m *MockRepository) WithTx(tx *gorm.DB) Repository {
	ret := _m.Called(tx)

	if len(ret) == 0 {
		panic("no return value specified for WithTx")
	}

	var r0 Repository
	if rf, ok := ret.Get(0).(func(*gorm.DB) Repository); ok {
		r0 = rf(tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Repository)
		}
	}

	return r0
}

// NewMockRepository creates a new instance of MockRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRepository {
	mock := &MockRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
