package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/common"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/domain"
)

type eventDBRepository struct{}

// NewEventDBRepository returns a new instance of a eventDBRepository.
func NewEventDBRepository() domain.EventRepository {
	return &eventDBRepository{}
}

// Add adds a new event to the database.
func (repo *eventDBRepository) Add(event *domain.Event) error {
	createdTime := time.Now().UTC()
	event.CreatedTime = &createdTime
	event.NormalizeTime()
	query := `INSERT INTO event (id, title, start_time, end_time, notify_time, description, user_id, 
              created_time, updated_time) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	result, err := db.Exec(
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
	count, err := result.RowsAffected()
	if common.IsErr(err) {
		return err
	}
	if count == 0 {
		return domain.ErrEventCreate
	}
	return nil
}

// Update updates an existing event in the database.
func (repo *eventDBRepository) Update(event *domain.Event) error {
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
	rowsAffected, err := result.RowsAffected()
	if common.IsErr(err) {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrEventNotExist
	}
	return err
}

// Get returns an event by ID.
func (repo *eventDBRepository) Get(eventID string) (*domain.Event, error) {
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrEventNotExist
		}
		return nil, err
	}
	e.NormalizeTime()
	if common.IsErr(err) {
		return nil, err
	}
	return &e, nil
}

// Delete removes an event by ID.
func (repo *eventDBRepository) Delete(eventID string) error {
	result, err := db.Exec(
		"DELETE FROM event WHERE id = $1", eventID,
	)
	if common.IsErr(err) {
		return err
	}
	_, err = result.RowsAffected()
	return err
}

// DeleteEventBeforeDate removes an event before date.
func (repo *eventDBRepository) DeleteEventBeforeDate(date time.Time) error {
	result, err := db.Exec(
		"DELETE FROM event WHERE start_time <= $1", date,
	)
	if common.IsErr(err) {
		return err
	}
	_, err = result.RowsAffected()
	return err
}

func (repo *eventDBRepository) getEvents(
	query string, args ...interface{},
) ([]*domain.Event, error) {
	rows, err := db.Query(query, args...)
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
		e.NormalizeTime()
		events = append(events, &e)
	}
	if err := rows.Err(); common.IsErr(err) {
		return nil, err
	}
	return events, nil
}

// GetEventsByPeriod returns a list of events for a period of time.
func (repo *eventDBRepository) GetEventsByPeriod(
	startTime, endTime time.Time,
) ([]*domain.Event, error) {
	query := `SELECT id, title, start_time, end_time, notify_time, description, user_id, 
       		  created_time FROM event WHERE start_time >= $1 AND end_time <= $2`
	return repo.getEvents(query, startTime, endTime)
}

// GetEventsByNotifyTime returns a list of events by notify time.
func (repo *eventDBRepository) GetEventsByNotifyTime(
	startTime, endTime time.Time,
) ([]*domain.Event, error) {
	query := `SELECT id, title, start_time, end_time, notify_time, description, user_id, 
	   		  created_time FROM event WHERE notify_time >= $1 AND notify_time <= $2`
	return repo.getEvents(query, startTime, endTime)
}

type eventCacheRepository struct{}

// NewEventCacheRepository returns a new instance of a eventCacheRepository.
func NewEventCacheRepository() domain.EventRepository {
	return &eventCacheRepository{}
}

// Add adds a new event to the cache.
func (repo *eventCacheRepository) Add(event *domain.Event) error {
	key := []byte(event.ID)
	createdTime := time.Now().UTC()
	event.CreatedTime = &createdTime
	event.NormalizeTime()
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

// Update updates an existing event in the cache.
func (repo *eventCacheRepository) Update(event *domain.Event) error {
	key := []byte(event.ID)
	event.NormalizeTime()
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

// Get returns an event by ID.
func (repo *eventCacheRepository) Get(eventID string) (*domain.Event, error) {
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

// Delete removes an event by ID.
func (repo *eventCacheRepository) Delete(eventID string) error {
	key := []byte(eventID)
	if _, err := cacheDB.Get(key); common.IsErr(err) {
		return domain.ErrEventNotExist
	}
	if affected := cacheDB.Del(key); !affected {
		return errors.New("event deletion failed")
	}
	return nil
}

// DeleteEventBeforeDate removes an event before date.
func (repo *eventCacheRepository) DeleteEventBeforeDate(date time.Time) error {
	keys := cacheDB.Keys()
	for _, key := range keys {
		data, err := cacheDB.Get(key)
		if common.IsErr(err) {
			return err
		}
		event := &domain.Event{}
		err = json.Unmarshal(data, event)
		if common.IsErr(err) {
			return err
		}
		if event.StartTime.Before(date) || event.StartTime.Equal(date) {
			if affected := cacheDB.Del(key); !affected {
				return errors.New("event deletion failed")
			}
		}
	}
	return nil
}

// GetEventsByPeriod returns a list of events for a period of time.
func (repo *eventCacheRepository) GetEventsByPeriod(
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
		if event.StartTime.After(startTime) && event.EndTime.Before(endTime) {
			result = append(result, event)
		}
	}
	return result, nil
}

// GetEventsByNotifyTime returns a list of events by notify time.
func (repo *eventCacheRepository) GetEventsByNotifyTime(
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
		if event.NotifyTime.After(startTime) && event.NotifyTime.Before(endTime) {
			result = append(result, event)
		}
	}
	return result, nil
}
