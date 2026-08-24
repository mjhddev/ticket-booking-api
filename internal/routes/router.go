package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mjhddev/ticket-booking-api/internal/handler"
	"github.com/mjhddev/ticket-booking-api/internal/middleware"
	"github.com/mjhddev/ticket-booking-api/internal/utils"
)

type Handlers struct {
	Auth    *handler.AuthHandler
	Profile *handler.ProfileHandler
}

func SetupRouter(h Handlers, jwtManager *utils.JWTManager) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api/v1")

	auth := api.Group("/auth")
	{
		auth.POST("/register", h.Auth.Register)
		auth.POST("/login", h.Auth.Login)
	}

	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(jwtManager))
	{
		protected.GET("/profile", h.Profile.Get)
	}

	return router
}
