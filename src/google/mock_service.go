// Code generated by mockery v2.42.1. DO NOT EDIT.

package google

import (
	dto "aesir/src/google/dto"

	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"

	sheets "google.golang.org/api/sheets/v4"
)

// MockService is an autogenerated mock type for the Service type
type MockService struct {
	mock.Mock
}

// AppendRow provides a mock function with given fields: createVoCDto
func (_m *MockService) AppendRow(createVoCDto *dto.CreateVoCDto) error {
	ret := _m.Called(createVoCDto)

	if len(ret) == 0 {
		panic("no return value specified for AppendRow")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*dto.CreateVoCDto) error); ok {
		r0 = rf(createVoCDto)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindSheet provides a mock function with given fields:
func (_m *MockService) FindSheet() (*sheets.ValueRange, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for FindSheet")
	}

	var r0 *sheets.ValueRange
	var r1 error
	if rf, ok := ret.Get(0).(func() (*sheets.ValueRange, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *sheets.ValueRange); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sheets.ValueRange)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WithTx provides a mock function with given fields: tx
func (_m *MockService) WithTx(tx *gorm.DB) Service {
	ret := _m.Called(tx)

	if len(ret) == 0 {
		panic("no return value specified for WithTx")
	}

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