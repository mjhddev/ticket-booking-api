package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mjhddev/ticket-booking-api/internal/response"
)

type ProfileHandler struct{}

func NewProfileHandler() *ProfileHandler {
	return &ProfileHandler{}
}

func (h *ProfileHandler) Get(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Fail(
			c,
			http.StatusUnauthorized,
			"User information not found",
			nil,
		)
		return
	}

	email, _ := c.Get("email")
	role, _ := c.Get("role")

	response.Success(
		c,
		http.StatusOK,
		"Profile retrieved successfully",
		gin.H{
			"user_id": userID,
			"email":   email,
			"role":    role,
		},
	)
}
