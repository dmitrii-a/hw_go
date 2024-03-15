package domain

import "errors"

var (
	ErrEventExist    = errors.New("event already exists")
	ErrEventNotExist = errors.New("event doesn't exist")
	ErrEventCreate   = errors.New("event creation failed")
	ErrEndTime       = errors.New("end time must be greater than start time")
	ErrNotifyTime    = errors.New("notify time must be greater than start time")
	ErrUUID          = errors.New("invalid UUID")
)
