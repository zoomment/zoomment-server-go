package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
// In Go, struct fields that start with uppercase are "exported" (public)
// Fields that start with lowercase are private to the package
type Config struct {
	Port        string
	MongoDBURI  string
	JWTSecret   string
	DashboardURL string
	BrandName   string
	AdminEmail  string
	BotEmail    EmailConfig
}

// EmailConfig holds SMTP configuration
type EmailConfig struct {
	Address  string
	Password string
	Host     string
	Port     int
}

// Load reads environment variables and returns a Config struct
// In Go, functions return values. Multiple return values are common.
// The pattern (value, error) is idiomatic Go.
func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if file doesn't exist)
	// _ means we're ignoring the error intentionally
	_ = godotenv.Load()

	// getEnv is a helper function defined below
	// It gets an env var or returns a default value
	port := getEnv("PORT", "8080")
	mongoURI := getEnv("MONGODB_URI", "mongodb://localhost:27017/zoomment")
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key")
	dashboardURL := getEnv("DASHBOARD_URL", "http://localhost:3000")
	brandName := getEnv("BRAND_NAME", "Zoomment")
	adminEmail := getEnv("ADMIN_EMAIL_ADDR", "")

	// Parse email port as integer
	emailPort, err := strconv.Atoi(getEnv("BOT_EMAIL_PORT", "465"))
	if err != nil {
		emailPort = 465
	}

	// Create and return the config
	// &Config{} creates a pointer to a new Config struct
	// In Go, we often return pointers to avoid copying large structs
	config := &Config{
		Port:        port,
		MongoDBURI:  mongoURI,
		JWTSecret:   jwtSecret,
		DashboardURL: dashboardURL,
		BrandName:   brandName,
		AdminEmail:  adminEmail,
		BotEmail: EmailConfig{
			Address:  getEnv("BOT_EMAIL_ADDR", ""),
			Password: getEnv("BOT_EMAIL_PASS", ""),
			Host:     getEnv("BOT_EMAIL_HOST", "smtp.gmail.com"),
			Port:     emailPort,
		},
	}

	return config, nil
}

// getEnv gets an environment variable or returns a default value
// This is a private function (lowercase first letter)
func getEnv(key, defaultValue string) string {
	// os.Getenv returns empty string if not found
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// MustLoad loads config and panics if it fails
// "Must" prefix is a Go convention for functions that panic on error
func MustLoad() *Config {
	config, err := Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	return config
}

