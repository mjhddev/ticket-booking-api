package errs

import "errors"

var (
	// Auth
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")

	// Event
	ErrEventNotFound          = errors.New("event not found")
	ErrEventNotOwned          = errors.New("you do not own this event")
	ErrInvalidEventTime       = errors.New("end time must be after start time")
	ErrInvalidEventStatus     = errors.New("invalid event status")
	ErrEventCannotBePublished = errors.New("event cannot be published")

	// Seat
	ErrSeatAlreadyExists = errors.New("seat already exists")
	ErrSeatNotFound      = errors.New("seat not found")
	ErrSeatNotAvailable  = errors.New("seat is not available")

	// Booking
	ErrInvalidBooking       = errors.New("invalid booking")
	ErrBookingNotFound      = errors.New("booking not found")
	ErrBookingForbidden     = errors.New("you do not have access to this booking")
	ErrInvalidBookingStatus = errors.New("invalid booking status")
)
