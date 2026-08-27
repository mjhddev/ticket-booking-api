package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mjhddev/ticket-booking-api/internal/dto"
	"github.com/mjhddev/ticket-booking-api/internal/errs"
	"github.com/mjhddev/ticket-booking-api/internal/response"
	"github.com/mjhddev/ticket-booking-api/internal/service"
)

type SeatHandler struct {
	seatService service.SeatService
}

func NewSeatHandler(seatService service.SeatService) *SeatHandler {
	return &SeatHandler{
		seatService: seatService,
	}
}

func (h *SeatHandler) CreateSeats(c *gin.Context) {
	eventID, err := getIDParam(c)
	if err != nil {
		response.Fail(
			c,
			http.StatusBadRequest,
			"Invalid event ID",
			nil,
		)
		return
	}

	organizerID, err := getUserID(c)
	if err != nil {
		response.Fail(
			c,
			http.StatusUnauthorized,
			"Invalid user information",
			nil,
		)
		return
	}

	var req dto.CreateSeatsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(
			c,
			http.StatusBadRequest,
			"Invalid request body",
			err.Error(),
		)
		return
	}

	seats, err := h.seatService.CreateSeats(
		c.Request.Context(),
		organizerID,
		eventID,
		req,
	)
	if err != nil {
		handleSeatError(c, err)
		return
	}

	response.Success(
		c,
		http.StatusCreated,
		"Seats created successfully",
		seats,
	)
}

func (h *SeatHandler) GetSeatsByEvent(c *gin.Context) {
	eventID, err := getIDParam(c)
	if err != nil {
		response.Fail(
			c,
			http.StatusBadRequest,
			"Invalid event ID",
			nil,
		)
		return
	}

	seats, err := h.seatService.GetSeatsByEvent(
		c.Request.Context(),
		eventID,
	)
	if err != nil {
		handleSeatError(c, err)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"Seats retrieved successfully",
		seats,
	)
}

func handleSeatError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, errs.ErrEventNotFound):
		response.Fail(
			c,
			http.StatusNotFound,
			"Event not found",
			nil,
		)

	case errors.Is(err, errs.ErrEventNotOwned):
		response.Fail(
			c,
			http.StatusForbidden,
			"You do not own this event",
			nil,
		)

	case errors.Is(err, errs.ErrSeatAlreadyExists):
		response.Fail(
			c,
			http.StatusConflict,
			"Seat already exists",
			nil,
		)

	case errors.Is(err, errs.ErrSeatNotFound):
		response.Fail(
			c,
			http.StatusBadRequest,
			"No valid seats provided",
			nil,
		)

	default:
		response.Fail(
			c,
			http.StatusInternalServerError,
			"Internal server error",
			nil,
		)
	}
}
