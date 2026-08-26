package dto

import "time"

type CreateEventRequest struct {
	Title       string    `json:"title" binding:"required,min=3,max=150"`
	Description string    `json:"description" binding:"max=5000"`
	Location    string    `json:"location" binding:"required,max=255"`
	StartTime   time.Time `json:"start_time" binding:"required"`
	EndTime     time.Time `json:"end_time" binding:"required"`
}

type EventResponse struct {
	ID          uint64    `json:"id"`
	OrganizerID uint64    `json:"organizer_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateEventRequest struct {
	Title       string    `json:"title" binding:"required,min=3,max=150"`
	Description string    `json:"description" binding:"max=5000"`
	Location    string    `json:"location" binding:"required,max=255"`
	StartTime   time.Time `json:"start_time" binding:"required"`
	EndTime     time.Time `json:"end_time" binding:"required"`
}
