package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mjhddev/ticket-booking-api/internal/dto"
	"github.com/mjhddev/ticket-booking-api/internal/errs"
	"github.com/mjhddev/ticket-booking-api/internal/response"
	"github.com/mjhddev/ticket-booking-api/internal/service"
)

type EventHandler struct {
	eventService service.EventService
}

func NewEventHandler(eventService service.EventService) *EventHandler {
	return &EventHandler{
		eventService: eventService,
	}
}

func (h *EventHandler) Create(c *gin.Context) {
	var req dto.CreateEventRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(
			c,
			http.StatusBadRequest,
			"Invalid request body",
			err.Error(),
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

	event, err := h.eventService.Create(
		c.Request.Context(),
		organizerID,
		req,
	)
	if err != nil {
		handleEventError(c, err)
		return
	}

	response.Success(
		c,
		http.StatusCreated,
		"Event created successfully",
		event,
	)
}

func (h *EventHandler) GetByID(c *gin.Context) {
	id, err := getIDParam(c)
	if err != nil {
		response.Fail(
			c,
			http.StatusBadRequest,
			"Invalid event ID",
			nil,
		)
		return
	}

	event, err := h.eventService.GetByID(
		c.Request.Context(),
		id,
	)
	if err != nil {
		handleEventError(c, err)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"Event retrieved successfully",
		event,
	)
}

func (h *EventHandler) GetPublished(c *gin.Context) {
	events, err := h.eventService.GetPublished(
		c.Request.Context(),
	)
	if err != nil {
		handleEventError(c, err)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"Events retrieved successfully",
		events,
	)
}

func (h *EventHandler) Update(c *gin.Context) {
	id, err := getIDParam(c)
	if err != nil {
		response.Fail(
			c,
			http.StatusBadRequest,
			"Invalid event ID",
			nil,
		)
		return
	}

	var req dto.UpdateEventRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(
			c,
			http.StatusBadRequest,
			"Invalid request body",
			err.Error(),
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

	event, err := h.eventService.Update(
		c.Request.Context(),
		organizerID,
		id,
		req,
	)
	if err != nil {
		handleEventError(c, err)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"Event updated successfully",
		event,
	)
}

func (h *EventHandler) Delete(c *gin.Context) {
	id, err := getIDParam(c)
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

	if err := h.eventService.Delete(
		c.Request.Context(),
		organizerID,
		id,
	); err != nil {
		handleEventError(c, err)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"Event deleted successfully",
		nil,
	)
}

func (h *EventHandler) Publish(c *gin.Context) {
	id, err := getIDParam(c)
	if err != nil {
		response.Fail(
			c,
			http.StatusBadRequest,
			"Invalid event ID",
			nil,
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

	roleValue, exists := c.Get("role")
	if !exists {
		response.Fail(
			c,
			http.StatusForbidden,
			"Role information not found",
			nil,
		)
		return
	}

	role, ok := roleValue.(string)
	if !ok {
		response.Fail(
			c,
			http.StatusForbidden,
			"Invalid role",
			nil,
		)
		return
	}

	isAdmin := role == "admin"

	if err := h.eventService.Publish(
		c.Request.Context(),
		userID,
		isAdmin,
		id,
	); err != nil {
		handleEventError(c, err)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"Event published successfully",
		nil,
	)
}

func getUserID(c *gin.Context) (uint64, error) {
	value, exists := c.Get("user_id")
	if !exists {
		return 0, errors.New("user_id not found")
	}

	userID, ok := value.(uint64)
	if !ok {
		return 0, errors.New("invalid user_id")
	}

	return userID, nil
}

func getIDParam(c *gin.Context) (uint64, error) {
	return strconv.ParseUint(c.Param("id"), 10, 64)
}

func handleEventError(c *gin.Context, err error) {
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

	case errors.Is(err, errs.ErrInvalidEventTime):
		response.Fail(
			c,
			http.StatusBadRequest,
			"End time must be after start time",
			nil,
		)

	case errors.Is(err, errs.ErrEventCannotBePublished):
		response.Fail(
			c,
			http.StatusBadRequest,
			"Event cannot be published",
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
