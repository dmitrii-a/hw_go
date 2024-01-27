package domain

import "errors"

var (
	ErrEventExist    = errors.New("event already exists")
	ErrEventNotExist = errors.New("event does not exist")
)
