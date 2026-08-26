package errs

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")

	ErrEventNotFound          = errors.New("event not found")
	ErrEventNotOwned          = errors.New("you do not own this event")
	ErrInvalidEventTime       = errors.New("end time must be after start time")
	ErrInvalidEventStatus     = errors.New("invalid event status")
	ErrEventCannotBePublished = errors.New("event cannot be published")
)
