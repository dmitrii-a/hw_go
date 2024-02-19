package repository

import (
	"database/sql/driver"
	"fmt"
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
	e := tests.GenerateTestEvent()
	err := s.repo.Add(e)
	s.NoError(err)
	return e
}

func (s *eventDBTestSuite) TestNonExistedEvent() {
	_ = s.setEventInDB()
	eventID := faker.UUIDHyphenated(options.WithGenerateUniqueValues(true))
	event, err := s.repo.Get(eventID)
	s.Error(err)
	s.Nil(event)
}

func (s *eventDBTestSuite) TestAddEvent() {
	e := tests.GenerateTestEvent()
	err := s.repo.Add(e)
	s.NoError(err)
}

func (s *eventDBTestSuite) TestAddEventWithExistingID() {
	e := tests.GenerateTestEvent()
	err := s.repo.Add(e)
	s.NoError(err)
	err = s.repo.Add(e)
	s.Error(err)
}

func (s *eventDBTestSuite) TestUpdateEvent() {
	e := tests.GenerateTestEvent()
	err := s.repo.Add(e)
	s.NoError(err)
	err = s.repo.Update(e)
	s.NoError(err)
}

func (s *eventDBTestSuite) TestGetEvent() {
	e := s.setEventInDB()
	result, err := s.repo.Get(e.ID)
	s.NoError(err)
	s.NotNil(result)
	s.Equal(e, result)
}

func (s *eventDBTestSuite) TestDeleteEvent() {
	e := s.setEventInDB()
	result, err := s.repo.Get(e.ID)
	s.NoError(err)
	s.NotNil(result)
	err = s.repo.Delete(e.ID)
	s.NoError(err)
}

func (s *eventDBTestSuite) TestListEventsForPeriodWithNoEvents() {
	startTime := time.Now()
	endTime := time.Now().Add(time.Hour)
	events, err := s.repo.ListForPeriod(startTime, endTime)
	s.NoError(err)
	s.Empty(events)
}

func (s *eventDBTestSuite) TestListEventsForPeriodWithSingleEvent() {
	e := s.setEventInDB()
	events, err := s.repo.ListForPeriod(e.StartTime, e.EndTime)
	s.NoError(err)
	s.Len(events, 1)
	s.Equal(e, events[0])
}

func (s *eventDBTestSuite) TestListEventsForPeriodWithMultipleEvents() {
	e1 := s.setEventInDB()
	e2 := s.setEventInDB()
	startTime, endTime := tests.GetEventStartEndTime(e1, e2)
	events, err := s.repo.ListForPeriod(startTime, endTime)
	s.NoError(err)
	s.Len(events, 2)
	s.Equal(events[0], e1)
	s.Equal(events[1], e2)
}

func (s *eventDBTestSuite) TestListEventsForPeriodWithEventOutsidePeriod() {
	e := s.setEventInDB()
	endTime := e.EndTime.Add(time.Minute)
	events, err := s.repo.ListForPeriod(endTime, endTime)
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
	e := tests.GenerateTestEvent()
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
	event, err := s.repo.Get(eventID)
	s.Error(err)
	s.Nil(event)
}

func (s *eventMockSQLTestSuite) TestAddEvent() {
	e := tests.GenerateTestEvent()
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
	err := s.repo.Add(e)
	s.NoError(err)
}

func (s *eventMockSQLTestSuite) TestAddEventWithExistingID() {
	e := tests.GenerateTestEvent()
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
	err := s.repo.Add(e)
	s.ErrorIs(err, duplicateErr)
}

func (s *eventMockSQLTestSuite) TestUpdateEvent() {
	e := tests.GenerateTestEvent()
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
	err := s.repo.Update(e)
	s.NoError(err)
}

func (s *eventMockSQLTestSuite) TestGetEvent() {
	e := s.setEventInDB()
	result, err := s.repo.Get(e.ID)
	s.NoError(err)
	s.NotNil(result)
	s.Equal(e, result)
}

func (s *eventMockSQLTestSuite) TestDeleteEvent() {
	e := s.setEventInDB()
	result, err := s.repo.Get(e.ID)
	s.NoError(err)
	s.NotNil(result)
	s.mock.ExpectExec("^DELETE FROM event WHERE id = \\$1$").
		WithArgs(e.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = s.repo.Delete(e.ID)
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
	events, err := s.repo.ListForPeriod(startTime, endTime)
	s.NoError(err)
	s.Empty(events)
}

func (s *eventMockSQLTestSuite) TestListEventsForPeriodWithSingleEvent() {
	e := s.setEventInDB()
	result, err := s.repo.Get(e.ID)
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
	events, err := s.repo.ListForPeriod(e.StartTime, e.EndTime)
	s.NoError(err)
	s.Len(events, 1)
	s.Equal(e, events[0])
}

func (s *eventMockSQLTestSuite) TestListEventsForPeriodWithMultipleEvents() {
	e1 := tests.GenerateTestEvent()
	e2 := tests.GenerateTestEvent()
	startTime, endTime := tests.GetEventStartEndTime(e1, e2)
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
	events, err := s.repo.ListForPeriod(startTime, endTime)
	s.NoError(err)
	s.Len(events, 2)
	s.Equal(events[0], e1)
	s.Equal(events[1], e2)
}

func (s *eventMockSQLTestSuite) TestListEventsForPeriodWithEventOutsidePeriod() {
	e := s.setEventInDB()
	result, err := s.repo.Get(e.ID)
	s.NoError(err)
	s.Equal(e, result)
	endTime := e.EndTime.Add(time.Minute)
	s.mockPeriodSelect(endTime, endTime)
	events, err := s.repo.ListForPeriod(endTime, endTime)
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
	event := tests.GenerateTestEvent()
	err := s.repo.Add(event)
	s.NoError(err)
	result, err := s.repo.Get(event.ID)
	s.NoError(err)
	s.Equal(event, result)
}

func (s *eventCacheTestSuite) TestAddExistEvent() {
	event := tests.GenerateTestEvent()
	err := s.repo.Add(event)
	s.NoError(err)
	err = s.repo.Add(event)
	s.Error(err)
}

func (s *eventCacheTestSuite) TestUpdateEvent() {
	event := tests.GenerateTestEvent()
	err := s.repo.Add(event)
	s.NoError(err)
	event.Title = "NewTitle"
	err = s.repo.Update(event)
	s.NoError(err)
	updatedEvent, err := s.repo.Get(event.ID)
	s.NoError(err)
	s.Equal("NewTitle", updatedEvent.Title)
}

func (s *eventCacheTestSuite) TestUpdateNonExistEvent() {
	event := tests.GenerateTestEvent()
	err := s.repo.Update(event)
	s.Error(err)
}

func (s *eventCacheTestSuite) TestGetNonExistEvent() {
	_, err := s.repo.Get(faker.UUIDHyphenated(options.WithGenerateUniqueValues(true)))
	s.Error(err)
}

func (s *eventCacheTestSuite) TestDeleteEvent() {
	event := tests.GenerateTestEvent()
	err := s.repo.Add(event)
	s.NoError(err)
	err = s.repo.Delete(event.ID)
	s.NoError(err)
	_, err = s.repo.Get(event.ID)
	s.Error(err)
}

func (s *eventCacheTestSuite) TestDeleteNonExistentEvent() {
	err := s.repo.Delete(faker.UUIDHyphenated(options.WithGenerateUniqueValues(true)))
	s.Error(err)
}

func (s *eventCacheTestSuite) TestListEventsForPeriod() {
	event1 := tests.GenerateTestEvent()
	event2 := tests.GenerateTestEvent()
	err := s.repo.Add(event1)
	s.NoError(err)
	err = s.repo.Add(event2)
	s.NoError(err)
	startTime, endTime := tests.GetEventStartEndTime(event1, event2)
	events, err := s.repo.ListForPeriod(startTime, endTime)
	s.NoError(err)
	s.Len(events, 2)
}

func (s *eventCacheTestSuite) TestListEventsForPeriodNoEvents() {
	events, err := s.repo.ListForPeriod(time.Now(), time.Now())
	s.NoError(err)
	s.Len(events, 0)
}

func TestRunCacheEventSuite(t *testing.T) {
	suite.Run(t, new(eventCacheTestSuite))
}
