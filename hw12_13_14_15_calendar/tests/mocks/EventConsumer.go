// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// EventConsumer is an autogenerated mock type for the EventConsumer type
type EventConsumer struct {
	mock.Mock
}

// Consume provides a mock function with given fields: name
func (_m *EventConsumer) Consume(name string) (<-chan []byte, error) {
	ret := _m.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for Consume")
	}

	var r0 <-chan []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (<-chan []byte, error)); ok {
		return rf(name)
	}
	if rf, ok := ret.Get(0).(func(string) <-chan []byte); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan []byte)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewEventConsumer creates a new instance of EventConsumer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEventConsumer(t interface {
	mock.TestingT
	Cleanup(func())
}) *EventConsumer {
	mock := &EventConsumer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
