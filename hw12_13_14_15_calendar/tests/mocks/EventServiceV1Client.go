// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	context "context"

	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	mock "github.com/stretchr/testify/mock"

	pb "github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/presentation/grpc/api/v1"
)

// EventServiceV1Client is an autogenerated mock type for the EventServiceV1Client type
type EventServiceV1Client struct {
	mock.Mock
}

// CreateEvent provides a mock function with given fields: ctx, in, opts
func (_m *EventServiceV1Client) CreateEvent(ctx context.Context, in *pb.EventRequest, opts ...grpc.CallOption) (*pb.EventResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for CreateEvent")
	}

	var r0 *pb.EventResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *pb.EventRequest, ...grpc.CallOption) (*pb.EventResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *pb.EventRequest, ...grpc.CallOption) *pb.EventResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pb.EventResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *pb.EventRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteEvent provides a mock function with given fields: ctx, in, opts
func (_m *EventServiceV1Client) DeleteEvent(ctx context.Context, in *pb.EventIDRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for DeleteEvent")
	}

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *pb.EventIDRequest, ...grpc.CallOption) (*emptypb.Empty, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *pb.EventIDRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *pb.EventIDRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEvent provides a mock function with given fields: ctx, in, opts
func (_m *EventServiceV1Client) GetEvent(ctx context.Context, in *pb.EventIDRequest, opts ...grpc.CallOption) (*pb.EventResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetEvent")
	}

	var r0 *pb.EventResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *pb.EventIDRequest, ...grpc.CallOption) (*pb.EventResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *pb.EventIDRequest, ...grpc.CallOption) *pb.EventResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pb.EventResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *pb.EventIDRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEventsByPeriod provides a mock function with given fields: ctx, in, opts
func (_m *EventServiceV1Client) GetEventsByPeriod(ctx context.Context, in *pb.TimePeriodRequest, opts ...grpc.CallOption) (*pb.EventsResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetEventsByPeriod")
	}

	var r0 *pb.EventsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *pb.TimePeriodRequest, ...grpc.CallOption) (*pb.EventsResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *pb.TimePeriodRequest, ...grpc.CallOption) *pb.EventsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pb.EventsResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *pb.TimePeriodRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateEvent provides a mock function with given fields: ctx, in, opts
func (_m *EventServiceV1Client) UpdateEvent(ctx context.Context, in *pb.EventRequest, opts ...grpc.CallOption) (*pb.EventResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for UpdateEvent")
	}

	var r0 *pb.EventResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *pb.EventRequest, ...grpc.CallOption) (*pb.EventResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *pb.EventRequest, ...grpc.CallOption) *pb.EventResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pb.EventResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *pb.EventRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewEventServiceV1Client creates a new instance of EventServiceV1Client. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEventServiceV1Client(t interface {
	mock.TestingT
	Cleanup(func())
}) *EventServiceV1Client {
	mock := &EventServiceV1Client{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
