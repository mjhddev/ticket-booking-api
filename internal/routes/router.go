package routes

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/mjhddev/ticket-booking-api/internal/handler"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	healthHandler := handlers.NewHealthHandler()
	router.GET("/health", healthHandler.Health)

	return router
}
