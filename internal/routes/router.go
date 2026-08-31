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
	Event   *handler.EventHandler
	Seat    *handler.SeatHandler
	Booking *handler.BookingHandler
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

	events := protected.Group("/events")
	{
		events.POST(
			"",
			middleware.RequireRole("admin", "organizer"),
			h.Event.Create,
		)

		events.GET("", h.Event.GetPublished)

		events.GET("/:id", h.Event.GetByID)

		events.PUT(
			"/:id",
			middleware.RequireRole("admin", "organizer"),
			h.Event.Update,
		)

		events.DELETE(
			"/:id",
			middleware.RequireRole("admin", "organizer"),
			h.Event.Delete,
		)

		events.PATCH(
			"/:id/publish",
			middleware.RequireRole("admin", "organizer"),
			h.Event.Publish,
		)

		events.POST(
			"/:id/seats",
			middleware.RequireRole("admin", "organizer"),
			h.Seat.CreateSeats,
		)

		events.GET(
			"/:id/seats",
			h.Seat.GetSeatsByEvent,
		)
	}
	bookings := protected.Group("/bookings")
	{
		bookings.POST("", h.Booking.CreateBooking)
		bookings.GET("/:id", h.Booking.GetBookingByID)
		bookings.DELETE("/:id", h.Booking.CancelBooking)
	}

	return router
}
