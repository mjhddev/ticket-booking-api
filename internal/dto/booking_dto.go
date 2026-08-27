package dto

import "time"

type CreateBookingRequest struct {
	EventID uint64   `json:"event_id" binding:"required"`
	SeatIDs []uint64 `json:"seat_ids" binding:"required,min=1"`
}

type BookingItemResponse struct {
	ID         uint64 `json:"id"`
	SeatID     uint64 `json:"seat_id"`
	SeatNumber string `json:"seat_number"`
	Price      int64  `json:"price"`
}

type BookingResponse struct {
	ID        uint64                `json:"id"`
	UserID    uint64                `json:"user_id"`
	EventID   uint64                `json:"event_id"`
	Status    string                `json:"status"`
	ExpiresAt time.Time             `json:"expires_at"`
	Items     []BookingItemResponse `json:"items"`
	CreatedAt time.Time             `json:"created_at"`
}
