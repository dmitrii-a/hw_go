package domain

import (
	"time"
)

// Event entity.
type Event struct {
	ID          string
	Title       string
	StartTime   time.Time
	EndTime     time.Time
	NotifyTime  time.Time
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
	e.EndTime = e.EndTime.UTC().Truncate(truncateTime)
	e.NotifyTime = e.NotifyTime.UTC().Truncate(truncateTime)
}

// Notification entity.
type Notification struct {
	EventID    string
	EventTitle string
	EventDate  time.Time
	UserToSend int64
}
