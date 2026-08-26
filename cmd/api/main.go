package main

import (
	"os"

	"github.com/mjhddev/ticket-booking-api/internal/config"
	"github.com/mjhddev/ticket-booking-api/internal/database"
	"github.com/mjhddev/ticket-booking-api/internal/handler"
	"github.com/mjhddev/ticket-booking-api/internal/logger"
	"github.com/mjhddev/ticket-booking-api/internal/repository"
	"github.com/mjhddev/ticket-booking-api/internal/routes"
	"github.com/mjhddev/ticket-booking-api/internal/service"
	"github.com/mjhddev/ticket-booking-api/internal/utils"
)

func main() {

	log := logger.New()

	log.Info("Loading configuration")

	cfg, err := config.Load()
	if err != nil {
		log.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	log.Info("Connecting PostgreSQL")

	db, err := database.Connect(cfg)
	if err != nil {
		log.Error("failed to connect database", "error", err)
		os.Exit(1)
	}

	log.Info("Database connected")

	jwtManager := utils.NewJWTManager(
		cfg.JWTSecret,
		cfg.JWTExpiredIn,
	)

	// =========================
	// Repository
	// =========================
	userRepository := repository.NewUserRepository(db)
	eventRepository := repository.NewEventRepository(db)

	// =========================
	// Service
	// =========================
	authService := service.NewAuthService(userRepository, jwtManager)
	eventService := service.NewEventService(eventRepository)

	// =========================
	// Handler
	// =========================
	authHandler := handler.NewAuthHandler(authService)
	profileHandler := handler.NewProfileHandler()
	eventHandler := handler.NewEventHandler(eventService)

	// =========================
	// Router
	// =========================
	router := routes.SetupRouter(routes.Handlers{
		Auth:    authHandler,
		Profile: profileHandler,
		Event:   eventHandler,
	}, jwtManager)

	log.Info("Starting HTTP server", "port", cfg.AppPort)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
