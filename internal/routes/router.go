package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mjhddev/ticket-booking-api/internal/handler"
)

type Handlers struct {
	Auth *handler.AuthHandler
}

func SetupRouter(h Handlers) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.Auth.Register)
			auth.POST("/login", h.Auth.Login)
		}
	}

	return router
}
