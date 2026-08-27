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

type BookingHandler struct {
	bookingService service.BookingService
}

func NewBookingHandler(bookingService service.BookingService) *BookingHandler {
	return &BookingHandler{
		bookingService: bookingService,
	}
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var req dto.CreateBookingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(
			c,
			http.StatusBadRequest,
			"Invalid request body",
			err.Error(),
		)
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		response.Fail(
			c,
			http.StatusUnauthorized,
			"Invalid user information",
			nil,
		)
		return
	}

	booking, err := h.bookingService.CreateBooking(
		c.Request.Context(),
		userID,
		req,
	)
	if err != nil {
		handleBookingError(c, err)
		return
	}

	response.Success(
		c,
		http.StatusCreated,
		"Booking created successfully",
		booking,
	)
}

func handleBookingError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, errs.ErrEventNotFound):
		response.Fail(
			c,
			http.StatusNotFound,
			"Event not found",
			nil,
		)

	case errors.Is(err, errs.ErrSeatNotFound):
		response.Fail(
			c,
			http.StatusNotFound,
			"Seat not found",
			nil,
		)

	case errors.Is(err, errs.ErrSeatNotAvailable):
		response.Fail(
			c,
			http.StatusConflict,
			"Seat is not available",
			nil,
		)

	case errors.Is(err, errs.ErrInvalidBooking):
		response.Fail(
			c,
			http.StatusBadRequest,
			"Invalid booking",
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
