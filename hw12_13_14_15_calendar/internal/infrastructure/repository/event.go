package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/common"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/domain"
)

type eventDBRepository struct{}

const truncateTime = time.Millisecond

// NewEventDBRepository returns a new instance of a eventDBRepository.
func NewEventDBRepository() domain.EventRepository {
	return &eventDBRepository{}
}

// normalizeTime set UTC and truncates time to milliseconds.
func normalizeTime(e *domain.Event) {
	if e.CreatedTime != nil {
		createdTime := e.CreatedTime.UTC().Truncate(truncateTime)
		e.CreatedTime = &createdTime
	} else {
		createdTime := time.Now().UTC().Truncate(truncateTime)
		e.CreatedTime = &createdTime
	}
	e.StartTime = e.StartTime.UTC().Truncate(truncateTime)
	e.EndTime = e.EndTime.UTC().Truncate(truncateTime)
	e.NotifyTime = e.NotifyTime.UTC().Truncate(truncateTime)
}

// AddEvent adds a new event to the database.
func (repo *eventDBRepository) AddEvent(event *domain.Event) error {
	normalizeTime(event)
	query := `INSERT INTO event (id, title, start_time, end_time, notify_time, description, user_id, 
              created_time, updated_time) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := db.Exec(
		query,
		event.ID,
		event.Title,
		event.StartTime,
		event.EndTime,
		event.NotifyTime,
		event.Description,
		event.UserID,
		event.CreatedTime,
		event.CreatedTime,
	)
	if common.IsErr(err) {
		return err
	}
	return nil
}

// UpdateEvent updates an existing event in the database.
func (repo *eventDBRepository) UpdateEvent(event *domain.Event) error {
	now := time.Now()
	query := `UPDATE event SET (
                  title, start_time, end_time, notify_time, description, user_id, created_time, updated_time
              ) = ($1, $2, $3, $4, $5, $6, $7, $8) WHERE id = $9`
	result, err := db.Exec(
		query,
		event.Title,
		event.StartTime,
		event.EndTime,
		event.NotifyTime,
		event.Description,
		event.UserID,
		event.CreatedTime,
		now,
		event.ID,
	)
	if common.IsErr(err) {
		return err
	}
	_, err = result.RowsAffected()
	return err
}

// GetEvent returns an event by ID.
func (repo *eventDBRepository) GetEvent(eventID string) (*domain.Event, error) {
	var e domain.Event
	query := `SELECT id, title, start_time, end_time, notify_time,
			  description, user_id, created_time FROM event WHERE id = $1`
	row := db.QueryRow(query, eventID)
	if row.Err() != nil {
		return nil, row.Err()
	}
	err := row.Scan(
		&e.ID,
		&e.Title,
		&e.StartTime,
		&e.EndTime,
		&e.NotifyTime,
		&e.Description,
		&e.UserID,
		&e.CreatedTime,
	)
	if common.IsErr(err) {
		return nil, err
	}
	createdTime := e.CreatedTime.UTC()
	e.CreatedTime = &createdTime
	normalizeTime(&e)
	if common.IsErr(err) {
		return nil, err
	}
	return &e, nil
}

// DeleteEvent removes an event by ID.
func (repo *eventDBRepository) DeleteEvent(eventID string) error {
	result, err := db.Exec(
		"DELETE FROM event WHERE id = $1", eventID,
	)
	if common.IsErr(err) {
		return err
	}
	_, err = result.RowsAffected()
	return err
}

// ListEventsForPeriod returns a list of events for a period of time.
func (repo *eventDBRepository) ListEventsForPeriod(
	startTime, endTime time.Time,
) ([]*domain.Event, error) {
	query := `SELECT id, title, start_time, end_time, notify_time, description, user_id, 
       		  created_time FROM event WHERE start_time >= $1 AND end_time <= $2`
	rows, err := db.Query(query, startTime, endTime)
	if common.IsErr(err) {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if common.IsErr(err) {
			common.Logger.Error().Err(err).Msg("error closing rows")
		}
	}(rows)
	var events []*domain.Event
	for rows.Next() {
		var e domain.Event
		if err := rows.Scan(
			&e.ID,
			&e.Title,
			&e.StartTime,
			&e.EndTime,
			&e.NotifyTime,
			&e.Description,
			&e.UserID,
			&e.CreatedTime,
		); common.IsErr(err) {
			return nil, err
		}
		normalizeTime(&e)
		events = append(events, &e)
	}

	if err := rows.Err(); common.IsErr(err) {
		return nil, err
	}

	return events, nil
}

type eventCacheRepository struct{}

// NewEventCacheRepository returns a new instance of a eventCacheRepository.
func NewEventCacheRepository() domain.EventRepository {
	return &eventCacheRepository{}
}

// AddEvent adds a new event to the cache.
func (repo *eventCacheRepository) AddEvent(event *domain.Event) error {
	key := []byte(event.ID)
	if _, err := cacheDB.Get(key); err == nil {
		return domain.ErrEventExist
	}
	data, err := json.Marshal(event)
	if common.IsErr(err) {
		return err
	}
	if err := cacheDB.Set(key, data, 0); common.IsErr(err) {
		return err
	}
	return nil
}

// UpdateEvent updates an existing event in the cache.
func (repo *eventCacheRepository) UpdateEvent(event *domain.Event) error {
	key := []byte(event.ID)
	if _, err := cacheDB.Get(key); common.IsErr(err) {
		return domain.ErrEventNotExist
	}
	data, err := json.Marshal(event)
	if common.IsErr(err) {
		return err
	}
	if err := cacheDB.Set(key, data, 0); common.IsErr(err) {
		return err
	}
	return nil
}

// GetEvent returns an event by ID.
func (repo *eventCacheRepository) GetEvent(eventID string) (*domain.Event, error) {
	key := []byte(eventID)
	data, err := cacheDB.Get(key)
	if common.IsErr(err) {
		return nil, domain.ErrEventNotExist
	}
	event := &domain.Event{}
	err = json.Unmarshal(data, event)
	if common.IsErr(err) {
		return nil, err
	}
	return event, nil
}

// DeleteEvent removes an event by ID.
func (repo *eventCacheRepository) DeleteEvent(eventID string) error {
	key := []byte(eventID)
	if _, err := cacheDB.Get(key); common.IsErr(err) {
		return domain.ErrEventNotExist
	}
	if affected := cacheDB.Del(key); !affected {
		return errors.New("event deletion failed")
	}
	return nil
}

// ListEventsForPeriod returns a list of events for a period of time.
func (repo *eventCacheRepository) ListEventsForPeriod(
	startTime, endTime time.Time,
) ([]*domain.Event, error) {
	keys := cacheDB.Keys()
	var result []*domain.Event
	for _, key := range keys {
		data, err := cacheDB.Get(key)
		if common.IsErr(err) {
			return nil, err
		}
		event := &domain.Event{}
		err = json.Unmarshal(data, event)
		if common.IsErr(err) {
			return nil, err
		}
		// For time equality in if statement
		startTime = startTime.Add(-time.Millisecond)
		endTime = endTime.Add(time.Millisecond)
		fmt.Println(startTime, event.StartTime, endTime, event.EndTime)
		if event.StartTime.After(startTime) && event.EndTime.Before(endTime) {
			result = append(result, event)
		}
	}
	return result, nil
}
