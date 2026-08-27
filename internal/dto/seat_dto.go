package dto

type CreateSeatRequest struct {
	SeatNumber string `json:"seat_number" binding:"required,max=20"`
	Price      int64  `json:"price" binding:"required,min=0"`
}

type CreateSeatsRequest struct {
	Seats []CreateSeatRequest `json:"seats" binding:"required,min=1"`
}

type SeatResponse struct {
	ID         uint64 `json:"id"`
	EventID    uint64 `json:"event_id"`
	SeatNumber string `json:"seat_number"`
	Price      int64  `json:"price"`
	Status     string `json:"status"`
}
