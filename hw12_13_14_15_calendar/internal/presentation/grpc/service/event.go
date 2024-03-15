package service

import (
	"context"
	"errors"
	"time"

	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/application"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/common"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/domain"
	pb "github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/presentation/grpc/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type grpcEventService struct {
	pb.EventServiceV1Server
	service *application.EventService
}

// NewGrpcEventService returns a new instance of the grpc event service.
func NewGrpcEventService() pb.EventServiceV1Server {
	return &grpcEventService{
		service: application.EventApplicationService,
	}
}

func (s *grpcEventService) convertEventTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

func (s *grpcEventService) convertEvent(e *domain.Event) *pb.Event {
	return &pb.Event{
		Id:          e.ID,
		Title:       e.Title,
		StartTime:   timestamppb.New(e.StartTime),
		EndTime:     s.convertEventTimestamp(e.EndTime),
		NotifyTime:  s.convertEventTimestamp(e.NotifyTime),
		UserId:      e.UserID,
		CreatedTime: timestamppb.New(*e.CreatedTime),
	}
}

func (s *grpcEventService) eventResponse(e *domain.Event) *pb.EventResponse {
	return &pb.EventResponse{
		Event: s.convertEvent(e),
	}
}

func (s *grpcEventService) eventsResponse(events []*domain.Event) *pb.EventsResponse {
	pbEvents := make([]*pb.Event, len(events))
	for i, e := range events {
		pbEvents[i] = s.convertEvent(e)
	}
	return &pb.EventsResponse{Events: pbEvents}
}

func (s *grpcEventService) convertToEvent(e *pb.Event) *domain.Event {
	var (
		endTime    *time.Time
		notifyTime *time.Time
	)
	if e.EndTime != nil {
		t := e.EndTime.AsTime()
		endTime = &t
	}
	if e.NotifyTime != nil {
		t := e.NotifyTime.AsTime()
		notifyTime = &t
	}
	return &domain.Event{
		ID:          e.Id,
		Title:       e.Title,
		StartTime:   e.StartTime.AsTime(),
		EndTime:     endTime,
		NotifyTime:  notifyTime,
		Description: e.Description,
		UserID:      e.UserId,
	}
}

// GetEvent returns an event by ID.
func (s *grpcEventService) GetEvent(
	_ context.Context,
	eventID *pb.EventIDRequest,
) (*pb.EventResponse, error) {
	err := eventID.ValidateAll()
	if common.IsErr(err) {
		return nil, err
	}
	event, err := s.service.Get(eventID.Id)
	if common.IsErr(err) {
		if errors.Is(err, domain.ErrEventNotExist) {
			return nil, status.Errorf(codes.NotFound, "event not found")
		}
		return nil, status.Errorf(codes.Unknown, "error getting event: %v", err)
	}
	return s.eventResponse(event), nil
}

// CreateEvent adds a new event.
func (s *grpcEventService) CreateEvent(
	_ context.Context,
	eventRequest *pb.EventRequest,
) (*pb.EventResponse, error) {
	err := eventRequest.ValidateAll()
	if common.IsErr(err) {
		return nil, err
	}
	event := s.convertToEvent(eventRequest.Event)
	err = s.service.Create(event)
	if common.IsErr(err) {
		for _, domainErr := range []error{domain.ErrEndTime, domain.ErrNotifyTime} {
			if errors.Is(err, domainErr) {
				return nil, status.Errorf(codes.InvalidArgument, err.Error())
			}
		}
		return nil, status.Errorf(codes.Unknown, err.Error())
	}
	return s.eventResponse(event), nil
}

// UpdateEvent updates an event.
func (s *grpcEventService) UpdateEvent(
	_ context.Context,
	eventRequest *pb.EventRequest,
) (*pb.EventResponse, error) {
	err := eventRequest.ValidateAll()
	if common.IsErr(err) {
		return nil, err
	}
	event := s.convertToEvent(eventRequest.Event)
	err = s.service.Update(event)
	if common.IsErr(err) {
		if errors.Is(err, domain.ErrEventNotExist) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		if errors.Is(err, domain.ErrUUID) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		return nil, status.Errorf(codes.Unknown, "error updating event: %v", err)
	}
	return s.eventResponse(event), nil
}

// DeleteEvent deletes an event by ID.
func (s *grpcEventService) DeleteEvent(
	_ context.Context,
	eventIDRequest *pb.EventIDRequest,
) (*emptypb.Empty, error) {
	err := eventIDRequest.ValidateAll()
	if common.IsErr(err) {
		return nil, err
	}
	err = s.service.Delete(eventIDRequest.Id)
	if common.IsErr(err) {
		if errors.Is(err, domain.ErrEventNotExist) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		if errors.Is(err, domain.ErrUUID) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		return nil, status.Errorf(codes.Unknown, "error deleting event: %v", err)
	}
	return new(emptypb.Empty), nil
}

// GetEventsByPeriod returns a list of events for the specified period.
func (s *grpcEventService) GetEventsByPeriod(
	_ context.Context,
	timePeriodRequest *pb.TimePeriodRequest,
) (*pb.EventsResponse, error) {
	err := timePeriodRequest.ValidateAll()
	if common.IsErr(err) {
		return nil, err
	}
	events, err := s.service.ListByPeriod(
		timePeriodRequest.StartTime.AsTime(), timePeriodRequest.EndTime.AsTime(),
	)
	if common.IsErr(err) {
		return nil, status.Errorf(codes.Unknown, "error getting events for period: %v", err)
	}
	return s.eventsResponse(events), nil
}
