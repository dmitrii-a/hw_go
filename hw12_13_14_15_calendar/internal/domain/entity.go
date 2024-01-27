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
	UserID      int
	CreatedTime *time.Time
}

// Notification entity.
type Notification struct {
	EventID    string
	EventTitle string
	EventDate  time.Time
	UserToSend string
}
