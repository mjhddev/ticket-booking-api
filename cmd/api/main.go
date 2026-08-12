package main

import (
	"os"

	"github.com/mjhddev/ticket-booking-api/internal/config"
	"github.com/mjhddev/ticket-booking-api/internal/database"
	"github.com/mjhddev/ticket-booking-api/internal/logger"
	"github.com/mjhddev/ticket-booking-api/internal/routes"
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

	router := routes.SetupRouter()

	log.Info("Starting HTTP server", "port", cfg.AppPort)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Error("server stopped", "error", err)
		os.Exit(1)
	}

	_ = db
}
