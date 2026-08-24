package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	JWTSecret    string
	JWTExpiredIn time.Duration

	RedisHost string
	RedisPort string

	RabbitMQHost     string
	RabbitMQPort     string
	RabbitMQUser     string
	RabbitMQPassword string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	duration, err := time.ParseDuration(getEnv("JWT_EXPIRED_IN"))
	if err != nil {
		log.Fatalf("invalid JWT_EXPIRED_IN: %v", err)
	}

	cfg := &Config{
		AppPort: getEnv("APP_PORT"),

		DBHost:     getEnv("DB_HOST"),
		DBPort:     getEnv("DB_PORT"),
		DBUser:     getEnv("DB_USER"),
		DBPassword: getEnv("DB_PASSWORD"),
		DBName:     getEnv("DB_NAME"),

		JWTSecret:    getEnv("JWT_SECRET"),
		JWTExpiredIn: duration,

		RedisHost: getEnv("REDIS_HOST"),
		RedisPort: getEnv("REDIS_PORT"),

		RabbitMQHost:     getEnv("RABBITMQ_HOST"),
		RabbitMQPort:     getEnv("RABBITMQ_PORT"),
		RabbitMQUser:     getEnv("RABBITMQ_USER"),
		RabbitMQPassword: getEnv("RABBITMQ_PASSWORD"),
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func getEnv(key string) string {
	return os.Getenv(key)
}

func (c *Config) Validate() error {
	required := map[string]string{
		"APP_PORT":          c.AppPort,
		"DB_HOST":           c.DBHost,
		"DB_PORT":           c.DBPort,
		"DB_USER":           c.DBUser,
		"DB_PASSWORD":       c.DBPassword,
		"DB_NAME":           c.DBName,
		"REDIS_HOST":        c.RedisHost,
		"REDIS_PORT":        c.RedisPort,
		"RABBITMQ_HOST":     c.RabbitMQHost,
		"RABBITMQ_PORT":     c.RabbitMQPort,
		"RABBITMQ_USER":     c.RabbitMQUser,
		"RABBITMQ_PASSWORD": c.RabbitMQPassword,
	}

	for key, value := range required {
		if value == "" {
			return fmt.Errorf("missing required environment variable: %s", key)
		}
	}

	return nil
}
