package tests

import (
	"math/rand"
	"time"

	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/common"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/domain"
	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
)

func GenerateTestEvent() *domain.Event {
	e := &domain.Event{}
	err := faker.FakeData(e)
	if common.IsErr(err) {
		panic(err)
	}
	e.EndTime = e.StartTime.Add(time.Hour * time.Duration(rand.Intn(48))) //nolint:gosec
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
	if e1.EndTime.After(e2.EndTime) {
		endTime = e1.EndTime
	} else {
		endTime = e2.EndTime
	}
	return startTime, endTime
}
