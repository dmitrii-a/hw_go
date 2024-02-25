package tests

import (
	"time"

	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/common"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/domain"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/presentation/grpc/api/v1"
	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func GenerateTestEvent() *domain.Event {
	e := &domain.Event{}
	err := faker.FakeData(e)
	if common.IsErr(err) {
		panic(err)
	}
	e.StartTime = time.Now()
	endTime := e.StartTime.Add(time.Hour)
	e.EndTime = &endTime
	e.ID = faker.UUIDHyphenated(options.WithGenerateUniqueValues(true))
	e.NormalizeTime()
	return e
}

func GetEventStartEndTime(e1, e2 *domain.Event) (time.Time, time.Time) {
	var startTime, endTime time.Time
	if e1.StartTime.Before(e2.StartTime) {
		startTime = e1.StartTime
	} else {
		startTime = e2.StartTime
	}
	if e1.EndTime.After(*e2.EndTime) {
		endTime = *e1.EndTime
	} else {
		endTime = *e2.EndTime
	}
	return startTime, endTime
}

func CreateTestEventRequest(event *domain.Event) *pb.EventRequest {
	var (
		endTime    *timestamppb.Timestamp
		notifyTime *timestamppb.Timestamp
	)
	if event.EndTime != nil {
		endTime = timestamppb.New(*event.EndTime)
	}
	if event.NotifyTime != nil {
		notifyTime = timestamppb.New(*event.NotifyTime)
	}
	return &pb.EventRequest{
		Event: &pb.Event{
			Id:          event.ID,
			Title:       event.Title,
			StartTime:   timestamppb.New(event.StartTime),
			EndTime:     endTime,
			NotifyTime:  notifyTime,
			Description: event.Description,
			UserId:      event.UserID,
		},
		RequestId: faker.UUIDDigit(options.WithGenerateUniqueValues(true)),
	}
}
