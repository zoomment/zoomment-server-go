package logger

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init initializes the global logger
// Call this once at startup
func Init(isDevelopment bool) {
	// Set global log level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if isDevelopment {
		// Pretty console output for development
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
		})
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		// JSON output for production (easy to parse by log aggregators)
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}
}

// Info logs an info message
func Info(msg string) {
	log.Info().Msg(msg)
}

// Debug logs a debug message
func Debug(msg string) {
	log.Debug().Msg(msg)
}

// Error logs an error message with error details
func Error(err error, msg string) {
	log.Error().Err(err).Msg(msg)
}

// Warn logs a warning message
func Warn(msg string) {
	log.Warn().Msg(msg)
}

// WithField logs with additional fields
// Usage: logger.WithField("user_id", "123").Info("User logged in")
func WithField(key string, value any) *zerolog.Event {
	return log.Info().Interface(key, value)
}

// GinLogger returns a Gin middleware for logging HTTP requests
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log after request is processed
		latency := time.Since(start)
		status := c.Writer.Status()

		event := log.Info()
		if status >= 400 && status < 500 {
			event = log.Warn()
		} else if status >= 500 {
			event = log.Error()
		}

		event.
			Str("method", c.Request.Method).
			Str("path", path).
			Str("query", query).
			Int("status", status).
			Dur("latency", latency).
			Str("ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent()).
			Msg("HTTP Request")
	}
}

// RequestLogger returns a logger with request context
func RequestLogger(c *gin.Context) zerolog.Logger {
	return log.With().
		Str("request_id", c.GetHeader("X-Request-ID")).
		Str("method", c.Request.Method).
		Str("path", c.Request.URL.Path).
		Logger()
}

