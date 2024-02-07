package service

import (
	"context"
	"errors"

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

func (s *grpcEventService) convertEvent(e *domain.Event) *pb.EventResponse {
	return &pb.EventResponse{
		Event: &pb.Event{
			Id:          e.ID,
			Title:       e.Title,
			StartTime:   timestamppb.New(e.StartTime),
			EndTime:     timestamppb.New(e.EndTime),
			NotifyTime:  timestamppb.New(e.NotifyTime),
			UserId:      e.UserID,
			CreatedTime: timestamppb.New(*e.CreatedTime),
		},
	}
}

func (s *grpcEventService) convertEvents(e []*domain.Event) *pb.EventsResponse {
	events := make([]*pb.Event, len(e))
	for i, event := range e {
		events[i] = &pb.Event{
			Id:          event.ID,
			Title:       event.Title,
			StartTime:   timestamppb.New(event.StartTime),
			EndTime:     timestamppb.New(event.EndTime),
			NotifyTime:  timestamppb.New(event.NotifyTime),
			UserId:      event.UserID,
			CreatedTime: timestamppb.New(*event.CreatedTime),
		}
	}
	return &pb.EventsResponse{Events: events}
}

func (s *grpcEventService) convertToEvent(e *pb.Event) *domain.Event {
	return &domain.Event{
		ID:          e.Id,
		Title:       e.Title,
		StartTime:   e.StartTime.AsTime(),
		EndTime:     e.EndTime.AsTime(),
		NotifyTime:  e.NotifyTime.AsTime(),
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
	return s.convertEvent(event), nil
}

// AddEvent adds a new event.
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
		return nil, status.Errorf(codes.Unknown, err.Error())
	}
	return s.convertEvent(event), nil
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
		return nil, status.Errorf(codes.Unknown, "error updating event: %v", err)
	}
	return s.convertEvent(event), nil
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
		return nil, status.Errorf(codes.Unknown, "error deleting event: %v", err)
	}
	return nil, nil
}

// GetEventsForPeriod returns a list of events for the specified period.
func (s *grpcEventService) GetEventsForPeriod(
	_ context.Context,
	timePeriodRequest *pb.TimePeriodRequest,
) (*pb.EventsResponse, error) {
	err := timePeriodRequest.ValidateAll()
	if common.IsErr(err) {
		return nil, err
	}
	events, err := s.service.ListForPeriod(
		timePeriodRequest.StartTime.AsTime(), timePeriodRequest.EndTime.AsTime(),
	)
	if common.IsErr(err) {
		return nil, status.Errorf(codes.Unknown, "error getting events for period: %v", err)
	}
	return s.convertEvents(events), nil
}
