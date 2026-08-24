package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mjhddev/ticket-booking-api/internal/response"
	"github.com/mjhddev/ticket-booking-api/internal/utils"
)

func AuthMiddleware(jwtManager *utils.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			response.Fail(
				c,
				http.StatusUnauthorized,
				"Authorization header is required",
				nil,
			)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Fail(
				c,
				http.StatusUnauthorized,
				"Invalid authorization header",
				nil,
			)
			c.Abort()
			return
		}

		token := parts[1]

		claims, err := jwtManager.ParseToken(token)
		if err != nil {
			response.Fail(
				c,
				http.StatusUnauthorized,
				"Invalid or expired token",
				nil,
			)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}
