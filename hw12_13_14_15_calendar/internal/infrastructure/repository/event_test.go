package repository

import (
	"database/sql/driver"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/common"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/domain"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/tests"
	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
)

func generateTestEvent() *domain.Event {
	e := &domain.Event{}
	err := faker.FakeData(e)
	if common.IsErr(err) {
		panic(err)
	}
	e.EndTime = e.StartTime.Add(time.Hour * time.Duration(rand.Intn(48)))
	e.ID = faker.UUIDHyphenated(options.WithGenerateUniqueValues(true))
	normalizeTime(e)
	return e
}

func getStartEndTime(e1, e2 *domain.Event) (time.Time, time.Time) {
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

type eventDBTestSuite struct {
	base tests.BaseDBTestSuite
	suite.Suite
	repo domain.EventRepository
}

func (s *eventDBTestSuite) SetupSuite() {
	s.base.SetupSuite()
	db = s.base.DB
	s.repo = NewEventDBRepository()
}

func (s *eventDBTestSuite) TearDownTest() {
	_, err := db.Exec("TRUNCATE TABLE event CASCADE")
	if common.IsErr(err) {
		panic(err)
	}
}

func (s *eventDBTestSuite) setEventInDB() *domain.Event {
	e := generateTestEvent()
	err := s.repo.AddEvent(e)
	s.NoError(err)
	return e
}

func (s *eventDBTestSuite) TestNonExistedEvent() {
	_ = s.setEventInDB()
	eventID := faker.UUIDHyphenated(options.WithGenerateUniqueValues(true))
	event, err := s.repo.GetEvent(eventID)
	s.Error(err)
	s.Nil(event)
}

func (s *eventDBTestSuite) TestAddEvent() {
	e := generateTestEvent()
	err := s.repo.AddEvent(e)
	s.NoError(err)
}

func (s *eventDBTestSuite) TestAddEventWithExistingID() {
	e := generateTestEvent()
	err := s.repo.AddEvent(e)
	s.NoError(err)
	err = s.repo.AddEvent(e)
	s.Error(err)
}

func (s *eventDBTestSuite) TestUpdateEvent() {
	e := generateTestEvent()
	err := s.repo.AddEvent(e)
	s.NoError(err)
	err = s.repo.UpdateEvent(e)
	s.NoError(err)
}

func (s *eventDBTestSuite) TestGetEvent() {
	e := s.setEventInDB()
	result, err := s.repo.GetEvent(e.ID)
	s.NoError(err)
	s.NotNil(result)
	s.Equal(e, result)
}

func (s *eventDBTestSuite) TestDeleteEvent() {
	e := s.setEventInDB()
	result, err := s.repo.GetEvent(e.ID)
	s.NoError(err)
	s.NotNil(result)
	err = s.repo.DeleteEvent(e.ID)
	s.NoError(err)
}

func (s *eventDBTestSuite) TestListEventsForPeriodWithNoEvents() {
	startTime := time.Now()
	endTime := time.Now().Add(time.Hour)
	events, err := s.repo.ListEventsForPeriod(startTime, endTime)
	s.NoError(err)
	s.Empty(events)
}

func (s *eventDBTestSuite) TestListEventsForPeriodWithSingleEvent() {
	e := s.setEventInDB()
	events, err := s.repo.ListEventsForPeriod(e.StartTime, e.EndTime)
	s.NoError(err)
	s.Len(events, 1)
	s.Equal(e, events[0])
}

func (s *eventDBTestSuite) TestListEventsForPeriodWithMultipleEvents() {
	e1 := s.setEventInDB()
	e2 := s.setEventInDB()
	startTime, endTime := getStartEndTime(e1, e2)
	events, err := s.repo.ListEventsForPeriod(startTime, endTime)
	s.NoError(err)
	s.Len(events, 2)
	s.Equal(events[0], e1)
	s.Equal(events[1], e2)
}

func (s *eventDBTestSuite) TestListEventsForPeriodWithEventOutsidePeriod() {
	e := s.setEventInDB()
	endTime := e.EndTime.Add(time.Minute)
	events, err := s.repo.ListEventsForPeriod(endTime, endTime)
	s.NoError(err)
	s.Empty(events)
}

func TestRunDBEventSuite(t *testing.T) {
	suite.Run(t, new(eventDBTestSuite))
}

type eventMockSQLTestSuite struct {
	suite.Suite
	repo domain.EventRepository
	mock sqlmock.Sqlmock
}

func (s *eventMockSQLTestSuite) SetupSuite() {
	s.repo = NewEventDBRepository()
}

func (s *eventMockSQLTestSuite) SetupTest() {
	mockDB, mock, err := sqlmock.New()
	if common.IsErr(err) {
		panic("An error was not expected when opening a stub database connection")
	}
	s.mock = mock
	db = sqlx.NewDb(mockDB, "sqlmock")
}

func (s *eventMockSQLTestSuite) setEventInDB() *domain.Event {
	e := generateTestEvent()
	rows := sqlmock.NewRows(
		[]string{
			"id", "title", "start_time", "end_time", "notify_time", "description", "user_id", "created_time",
		},
	).AddRow(e.ID, e.Title, e.StartTime, e.EndTime, e.NotifyTime, e.Description, e.UserID, e.CreatedTime)
	s.mock.ExpectQuery("^SELECT (.+) FROM event WHERE id = \\$1$").
		WithArgs(e.ID).
		WillReturnRows(rows)
	return e
}

func (s *eventMockSQLTestSuite) TestNonExistedEvent() {
	_ = s.setEventInDB()
	eventID := faker.UUIDHyphenated()
	event, err := s.repo.GetEvent(eventID)
	s.Error(err)
	s.Nil(event)
}

func (s *eventMockSQLTestSuite) TestAddEvent() {
	e := generateTestEvent()
	s.mock.ExpectExec("^INSERT INTO event (.+) VALUES (.+)$").
		WithArgs(
			e.ID,
			e.Title,
			e.StartTime,
			e.EndTime,
			e.NotifyTime,
			e.Description,
			e.UserID,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).WillReturnResult(sqlmock.NewResult(1, 1))
	err := s.repo.AddEvent(e)
	s.NoError(err)
}

func (s *eventMockSQLTestSuite) TestAddEventWithExistingID() {
	e := generateTestEvent()
	duplicateErr := fmt.Errorf("pq: duplicate key value violates unique constraint \"event_pkey\"")
	s.mock.ExpectExec("^INSERT INTO event (.+) VALUES (.+)$").
		WithArgs(
			e.ID,
			e.Title,
			e.StartTime,
			e.EndTime,
			e.NotifyTime,
			e.Description,
			e.UserID,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).WillReturnError(duplicateErr)
	err := s.repo.AddEvent(e)
	s.ErrorIs(err, duplicateErr)
}

func (s *eventMockSQLTestSuite) TestUpdateEvent() {
	e := generateTestEvent()
	s.mock.ExpectExec("^UPDATE event SET (.+) WHERE id = \\$9$").
		WithArgs(
			e.Title,
			e.StartTime,
			e.EndTime,
			e.NotifyTime,
			e.Description,
			e.UserID,
			e.CreatedTime,
			sqlmock.AnyArg(),
			e.ID,
		).WillReturnResult(sqlmock.NewResult(1, 1))
	err := s.repo.UpdateEvent(e)
	s.NoError(err)
}

func (s *eventMockSQLTestSuite) TestGetEvent() {
	e := s.setEventInDB()
	result, err := s.repo.GetEvent(e.ID)
	s.NoError(err)
	s.NotNil(result)
	s.Equal(e, result)
}

func (s *eventMockSQLTestSuite) TestDeleteEvent() {
	e := s.setEventInDB()
	result, err := s.repo.GetEvent(e.ID)
	s.NoError(err)
	s.NotNil(result)
	s.mock.ExpectExec("^DELETE FROM event WHERE id = \\$1$").
		WithArgs(e.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = s.repo.DeleteEvent(e.ID)
	s.NoError(err)
}

func (s *eventMockSQLTestSuite) mockPeriodSelect(
	startTime, endTime time.Time,
	valueRows ...[]driver.Value,
) {
	rows := sqlmock.NewRows(
		[]string{
			"id",
			"title",
			"start_time",
			"end_time",
			"notify_time",
			"description",
			"user_id",
			"created_time",
		},
	)
	for _, row := range valueRows {
		rows.AddRow(row...)
	}
	s.mock.ExpectQuery("^SELECT (.+) FROM event WHERE start_time >= \\$1 AND end_time <= \\$2$").
		WithArgs(startTime, endTime).
		WillReturnRows(rows)
}

func (s *eventMockSQLTestSuite) TestListEventsForPeriodWithNoEvents() {
	startTime := time.Now()
	endTime := time.Now().Add(time.Hour)
	s.mockPeriodSelect(startTime, endTime)
	events, err := s.repo.ListEventsForPeriod(startTime, endTime)
	s.NoError(err)
	s.Empty(events)
}

func (s *eventMockSQLTestSuite) TestListEventsForPeriodWithSingleEvent() {
	e := s.setEventInDB()
	result, err := s.repo.GetEvent(e.ID)
	s.NoError(err)
	s.Equal(e, result)
	s.mockPeriodSelect(
		result.StartTime,
		result.EndTime,
		[]driver.Value{
			e.ID,
			e.Title,
			e.StartTime,
			e.EndTime,
			e.NotifyTime,
			e.Description,
			e.UserID,
			e.CreatedTime,
		},
	)
	events, err := s.repo.ListEventsForPeriod(e.StartTime, e.EndTime)
	s.NoError(err)
	s.Len(events, 1)
	s.Equal(e, events[0])
}

func (s *eventMockSQLTestSuite) TestListEventsForPeriodWithMultipleEvents() {
	e1 := generateTestEvent()
	e2 := generateTestEvent()
	startTime, endTime := getStartEndTime(e1, e2)
	s.mockPeriodSelect(
		startTime,
		endTime,
		[]driver.Value{
			e1.ID,
			e1.Title,
			e1.StartTime,
			e1.EndTime,
			e1.NotifyTime,
			e1.Description,
			e1.UserID,
			e1.CreatedTime,
		},
		[]driver.Value{
			e2.ID,
			e2.Title,
			e2.StartTime,
			e2.EndTime,
			e2.NotifyTime,
			e2.Description,
			e2.UserID,
			e2.CreatedTime,
		},
	)
	events, err := s.repo.ListEventsForPeriod(startTime, endTime)
	s.NoError(err)
	s.Len(events, 2)
	s.Equal(events[0], e1)
	s.Equal(events[1], e2)
}

func (s *eventMockSQLTestSuite) TestListEventsForPeriodWithEventOutsidePeriod() {
	e := s.setEventInDB()
	result, err := s.repo.GetEvent(e.ID)
	s.NoError(err)
	s.Equal(e, result)
	endTime := e.EndTime.Add(time.Minute)
	s.mockPeriodSelect(endTime, endTime)
	events, err := s.repo.ListEventsForPeriod(endTime, endTime)
	s.NoError(err)
	s.Empty(events)
}

func TestRunMockSQLEventSuite(t *testing.T) {
	suite.Run(t, new(eventMockSQLTestSuite))
}

type eventCacheTestSuite struct {
	suite.Suite
	repo domain.EventRepository
}

func (s *eventCacheTestSuite) SetupSuite() {
	s.repo = NewEventCacheRepository()
}

func (s *eventCacheTestSuite) TearDownTest() {
	cacheDB.Clear()
}

func (s *eventCacheTestSuite) TestAddEvent() {
	event := generateTestEvent()
	err := s.repo.AddEvent(event)
	s.NoError(err)
	result, err := s.repo.GetEvent(event.ID)
	s.NoError(err)
	s.Equal(event, result)
}

func (s *eventCacheTestSuite) TestAddExistEvent() {
	event := generateTestEvent()
	err := s.repo.AddEvent(event)
	s.NoError(err)
	err = s.repo.AddEvent(event)
	s.Error(err)
}

func (s *eventCacheTestSuite) TestUpdateEvent() {
	event := generateTestEvent()
	err := s.repo.AddEvent(event)
	s.NoError(err)
	event.Title = "NewTitle"
	err = s.repo.UpdateEvent(event)
	s.NoError(err)
	updatedEvent, err := s.repo.GetEvent(event.ID)
	s.NoError(err)
	s.Equal("NewTitle", updatedEvent.Title)
}

func (s *eventCacheTestSuite) TestUpdateNonExistEvent() {
	event := generateTestEvent()
	err := s.repo.UpdateEvent(event)
	s.Error(err)
}

func (s *eventCacheTestSuite) TestGetNonExistEvent() {
	_, err := s.repo.GetEvent(faker.UUIDHyphenated(options.WithGenerateUniqueValues(true)))
	s.Error(err)
}

func (s *eventCacheTestSuite) TestDeleteEvent() {
	event := generateTestEvent()
	err := s.repo.AddEvent(event)
	s.NoError(err)
	err = s.repo.DeleteEvent(event.ID)
	s.NoError(err)
	_, err = s.repo.GetEvent(event.ID)
	s.Error(err)
}

func (s *eventCacheTestSuite) TestDeleteNonExistentEvent() {
	err := s.repo.DeleteEvent(faker.UUIDHyphenated(options.WithGenerateUniqueValues(true)))
	s.Error(err)
}

func (s *eventCacheTestSuite) TestListEventsForPeriod() {
	event1 := generateTestEvent()
	event2 := generateTestEvent()
	err := s.repo.AddEvent(event1)
	s.NoError(err)
	err = s.repo.AddEvent(event2)
	s.NoError(err)
	startTime, endTime := getStartEndTime(event1, event2)
	events, err := s.repo.ListEventsForPeriod(startTime, endTime)
	s.NoError(err)
	s.Len(events, 2)
}

func (s *eventCacheTestSuite) TestListEventsForPeriodNoEvents() {
	events, err := s.repo.ListEventsForPeriod(time.Now(), time.Now())
	s.NoError(err)
	s.Len(events, 0)
}

func TestRunCacheEventSuite(t *testing.T) {
	suite.Run(t, new(eventCacheTestSuite))
}
