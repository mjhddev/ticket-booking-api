package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mjhddev/ticket-booking-api/internal/response"
)

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get("role")
		if !exists {
			response.Fail(
				c,
				http.StatusForbidden,
				"Role information not found",
				nil,
			)
			c.Abort()
			return
		}

		userRole, ok := roleValue.(string)
		if !ok {
			response.Fail(
				c,
				http.StatusForbidden,
				"Invalid role",
				nil,
			)
			c.Abort()
			return
		}

		for _, role := range roles {
			if userRole == role {
				c.Next()
				return
			}
		}

		response.Fail(
			c,
			http.StatusForbidden,
			"You do not have permission to access this resource",
			nil,
		)
		c.Abort()
	}
}
