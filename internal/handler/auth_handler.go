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

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(
			c,
			http.StatusBadRequest,
			"Invalid request body",
			err.Error(),
		)
		return
	}

	res, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {

		if errors.Is(err, errs.ErrEmailAlreadyExists) {
			response.Fail(
				c,
				http.StatusConflict,
				"Email already exists",
				nil,
			)
			return
		}

		response.Fail(
			c,
			response.StatusCode(err),
			err.Error(),
			nil,
		)
		return
	}

	response.Success(
		c,
		http.StatusCreated,
		"User registered successfully",
		res,
	)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(
			c,
			http.StatusBadRequest,
			"Invalid request body",
			err.Error(),
		)
		return
	}

	res, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {

		switch {
		case errors.Is(err, errs.ErrInvalidCredentials):
			response.Fail(
				c,
				http.StatusUnauthorized,
				"Invalid email or password",
				nil,
			)
			return

		default:
			response.Fail(
				c,
				http.StatusInternalServerError,
				"Internal server error",
				nil,
			)
			return
		}
	}

	response.Success(
		c,
		http.StatusOK,
		"Login successful",
		res,
	)
}
