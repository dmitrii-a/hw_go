package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/application"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/domain"
	pb "github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/presentation/grpc/api/v1"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/tests"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/tests/mocks"
	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGrpcEventService_ConvertEvent(t *testing.T) {
	s := &grpcEventService{}
	event := tests.GenerateTestEvent()
	result := s.eventResponse(event)

	require.NotNil(t, result)
	require.NotNil(t, result.Event)
	require.Equal(t, event.ID, result.Event.Id)
	require.Equal(t, event.Title, result.Event.Title)
	require.Equal(t, timestamppb.New(event.StartTime), result.Event.StartTime)
	require.Equal(t, timestamppb.New(*event.EndTime), result.Event.EndTime)
	require.Equal(t, timestamppb.New(*event.NotifyTime), result.Event.NotifyTime)
	require.Equal(t, event.UserID, result.Event.UserId)
	require.Equal(t, timestamppb.New(*event.CreatedTime), result.Event.CreatedTime)
}

func TestGrpcEventService_ConvertEvents(t *testing.T) {
	s := &grpcEventService{}
	events := []*domain.Event{
		tests.GenerateTestEvent(),
		tests.GenerateTestEvent(),
	}
	result := s.eventsResponse(events)

	require.NotNil(t, result)
	require.NotNil(t, result.Events)
	require.Len(t, result.Events, 2)

	require.Equal(t, events[0].ID, result.Events[0].Id)
	require.Equal(t, events[0].Title, result.Events[0].Title)
	require.Equal(t, timestamppb.New(events[0].StartTime), result.Events[0].StartTime)
	require.Equal(t, timestamppb.New(*events[0].EndTime), result.Events[0].EndTime)
	require.Equal(t, timestamppb.New(*events[0].NotifyTime), result.Events[0].NotifyTime)
	require.Equal(t, events[0].UserID, result.Events[0].UserId)
	require.Equal(t, timestamppb.New(*events[0].CreatedTime), result.Events[0].CreatedTime)

	require.Equal(t, events[1].ID, result.Events[1].Id)
	require.Equal(t, events[1].Title, result.Events[1].Title)
	require.Equal(t, timestamppb.New(events[1].StartTime), result.Events[1].StartTime)
	require.Equal(t, timestamppb.New(*events[1].EndTime), result.Events[1].EndTime)
	require.Equal(t, timestamppb.New(*events[1].NotifyTime), result.Events[1].NotifyTime)
	require.Equal(t, events[1].UserID, result.Events[1].UserId)
	require.Equal(t, timestamppb.New(*events[1].CreatedTime), result.Events[1].CreatedTime)
}

func TestGrpcEventService_GetEvent(t *testing.T) {
	mockRepo := new(mocks.EventRepository)
	event := tests.GenerateTestEvent()
	mockRepo.On("Get", event.ID).Return(event, nil)
	s := grpcEventService{service: application.NewEventService(mockRepo)}
	result, err := s.GetEvent(context.Background(), &pb.EventIDRequest{Id: event.ID})

	mockRepo.AssertExpectations(t)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, event.ID, result.Event.Id)
}

func TestGrpcEventService_AddEvent(t *testing.T) {
	mockRepo := new(mocks.EventRepository)
	event := tests.GenerateTestEvent()
	event.CreatedTime = nil
	mockRepo.On("Add", mock.Anything).Run(func(args mock.Arguments) {
		e := args[0].(*domain.Event)
		createTime := time.Now().UTC()
		e.CreatedTime = &createTime
		e.NormalizeTime()
	}).Return(nil)

	s := grpcEventService{service: application.NewEventService(mockRepo)}
	result, err := s.CreateEvent(context.Background(), tests.CreateTestEventRequest(event))

	mockRepo.AssertExpectations(t)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotEmpty(t, result.Event.Id)
}

func TestGrpcEventService_AddEventError(t *testing.T) {
	mockRepo := new(mocks.EventRepository)
	event := tests.GenerateTestEvent()
	event.CreatedTime = nil
	mockRepo.On("Add", mock.Anything).Return(domain.ErrEventCreate)

	s := grpcEventService{service: application.NewEventService(mockRepo)}
	result, err := s.CreateEvent(context.Background(), tests.CreateTestEventRequest(event))

	mockRepo.AssertExpectations(t)
	require.Error(t, err)
	require.Nil(t, result)
}

func TestGrpcEventService_UpdateEvent(t *testing.T) {
	mockRepo := new(mocks.EventRepository)
	event := tests.GenerateTestEvent()
	event.CreatedTime = nil
	mockRepo.On("Update", event).Run(func(args mock.Arguments) {
		e := args[0].(*domain.Event)
		createTime := time.Now().UTC()
		e.CreatedTime = &createTime
		e.NormalizeTime()
	}).Return(nil)

	s := grpcEventService{service: application.NewEventService(mockRepo)}
	result, err := s.UpdateEvent(context.Background(), tests.CreateTestEventRequest(event))

	mockRepo.AssertExpectations(t)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, event.ID, result.Event.Id)
}

func TestGrpcEventService_UpdateEventError(t *testing.T) {
	mockRepo := new(mocks.EventRepository)
	event := tests.GenerateTestEvent()
	event.CreatedTime = nil
	mockRepo.On("Update", event).Return(domain.ErrEventNotExist)

	s := grpcEventService{service: application.NewEventService(mockRepo)}
	result, err := s.UpdateEvent(context.Background(), tests.CreateTestEventRequest(event))

	mockRepo.AssertExpectations(t)
	require.Error(t, err)
	require.Nil(t, result)
}

func TestGrpcEventService_DeleteEvent(t *testing.T) {
	mockRepo := new(mocks.EventRepository)
	event := tests.GenerateTestEvent()
	event.CreatedTime = nil
	mockRepo.On("Delete", event.ID).Return(nil)

	s := grpcEventService{service: application.NewEventService(mockRepo)}
	result, err := s.DeleteEvent(context.Background(), &pb.EventIDRequest{Id: event.ID})

	mockRepo.AssertExpectations(t)
	require.NoError(t, err)
	require.Equal(t, new(emptypb.Empty), result)
}

func TestGrpcEventService_DeleteEventError(t *testing.T) {
	mockRepo := new(mocks.EventRepository)
	event := tests.GenerateTestEvent()
	event.CreatedTime = nil
	mockRepo.On("Delete", event.ID).Return(nil)

	s := grpcEventService{service: application.NewEventService(mockRepo)}
	result, err := s.DeleteEvent(
		context.Background(),
		&pb.EventIDRequest{
			Id:        event.ID,
			RequestId: faker.UUIDDigit(options.WithGenerateUniqueValues(true)),
		},
	)

	mockRepo.AssertExpectations(t)
	require.NoError(t, err)
	require.Equal(t, new(emptypb.Empty), result)
}

func TestGrpcEventService_GetEventsByPeriod(t *testing.T) {
	mockRepo := new(mocks.EventRepository)
	events := []*domain.Event{tests.GenerateTestEvent(), tests.GenerateTestEvent()}
	startTime, endTime := tests.GetEventStartEndTime(events[0], events[1])
	mockRepo.On("GetEventsByPeriod", startTime, endTime).Return(events, nil)

	s := grpcEventService{service: application.NewEventService(mockRepo)}
	result, err := s.GetEventsByPeriod(
		context.Background(),
		&pb.TimePeriodRequest{
			StartTime: timestamppb.New(startTime),
			EndTime:   timestamppb.New(endTime),
			RequestId: faker.UUIDDigit(options.WithGenerateUniqueValues(true)),
		},
	)

	mockRepo.AssertExpectations(t)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Events, 2)
}

func TestGrpcEventService_GetEventsByPeriodError(t *testing.T) {
	mockRepo := new(mocks.EventRepository)
	events := []*domain.Event{tests.GenerateTestEvent(), tests.GenerateTestEvent()}
	startTime, endTime := tests.GetEventStartEndTime(events[0], events[1])
	mockRepo.On("GetEventsByPeriod", startTime, endTime).Return(nil, errors.New("error"))

	s := grpcEventService{service: application.NewEventService(mockRepo)}
	result, err := s.GetEventsByPeriod(
		context.Background(),
		&pb.TimePeriodRequest{
			StartTime: timestamppb.New(startTime),
			EndTime:   timestamppb.New(endTime),
			RequestId: faker.UUIDDigit(options.WithGenerateUniqueValues(true)),
		},
	)

	mockRepo.AssertExpectations(t)
	require.Error(t, err)
	require.Nil(t, result)
}
