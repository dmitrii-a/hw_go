package domain

import (
	"time"

	"github.com/google/uuid"
)

// Event entity.
type Event struct {
	ID          string
	Title       string
	StartTime   time.Time
	EndTime     *time.Time
	NotifyTime  *time.Time
	Description string
	UserID      int64
	CreatedTime *time.Time
}

// NormalizeTime set UTC and truncates time to milliseconds.
func (e *Event) NormalizeTime() {
	const truncateTime = time.Millisecond
	if e.CreatedTime != nil {
		createdTime := e.CreatedTime.UTC().Truncate(truncateTime)
		e.CreatedTime = &createdTime
	}
	e.StartTime = e.StartTime.UTC().Truncate(truncateTime)
	if e.EndTime != nil {
		t := e.EndTime.UTC().Truncate(truncateTime)
		e.EndTime = &t
	}
	if e.NotifyTime != nil {
		t := e.NotifyTime.UTC().Truncate(truncateTime)
		e.NotifyTime = &t
	}
}

func (e *Event) NewUUID() string {
	return uuid.New().String()
}

func (e *Event) Validate() error {
	if e.EndTime != nil && e.StartTime.After(*e.EndTime) {
		return ErrEndTime
	}
	if e.NotifyTime != nil && e.StartTime.After(*e.NotifyTime) {
		return ErrNotifyTime
	}
	if _, err := uuid.Parse(e.ID); err != nil {
		return ErrUUID
	}
	return nil
}

// Notification entity.
type Notification struct {
	EventID    string
	EventTitle string
	EventDate  time.Time
	UserToSend int64
}
