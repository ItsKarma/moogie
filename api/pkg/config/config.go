package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	// App configuration
	AppEnv  string
	AppPort string

	// Database configuration
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// CORS configuration
	AllowedOrigins []string
}

// Load loads the configuration from environment variables and .env file
func Load() *Config {
	// Load .env file if it exists (ignore error in production)
	if err := godotenv.Load(); err != nil && os.Getenv("APP_ENV") != "production" {
		log.Println("Warning: .env file not found, using environment variables only")
	}

	dbPort, err := strconv.Atoi(getEnvOrDefault("DB_PORT", "5432"))
	if err != nil {
		log.Fatal("Invalid DB_PORT value")
	}

	return &Config{
		AppEnv:  getEnvOrDefault("APP_ENV", "development"),
		AppPort: getEnvOrDefault("APP_PORT", "8080"),

		DBHost:     getEnvOrDefault("DB_HOST", "localhost"),
		DBPort:     dbPort,
		DBUser:     getEnvOrDefault("DB_USER", "moogie"),
		DBPassword: getEnvOrDefault("DB_PASSWORD", "moogie"),
		DBName:     getEnvOrDefault("DB_NAME", "moogie"),
		DBSSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),

		AllowedOrigins: parseAllowedOrigins(getEnvOrDefault("ALLOWED_ORIGINS", "http://localhost:3000")),
	}
}

func parseAllowedOrigins(origins string) []string {
	if origins == "" {
		return []string{"http://localhost:3000"}
	}
	return strings.Split(origins, ",")
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
