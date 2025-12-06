package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"zoomment-server/internal/config"
	"zoomment-server/internal/database"
	"zoomment-server/internal/logger"
	"zoomment-server/internal/middleware"
	"zoomment-server/internal/routes"

	_ "zoomment-server/docs" // Import swagger docs
)

func main() {
	// Load configuration
	cfg := config.MustLoad()

	// Initialize logger
	isDev := os.Getenv("GIN_MODE") != "release"
	logger.Init(isDev)

	logger.Info("🚀 Starting " + cfg.BrandName + " Server...")

	// Connect to MongoDB
	if err := database.Connect(cfg.MongoDBURI); err != nil {
		logger.Error(err, "Failed to connect to MongoDB")
		os.Exit(1)
	}

	// Set Gin mode
	if !isDev {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.New()
	router.Use(logger.GinLogger())
	router.Use(gin.Recovery())
	router.Use(cors.Default())
	router.Use(middleware.Auth(cfg))

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup API routes
	routes.Setup(router, cfg)

	// Start server
	go func() {
		addr := ":" + cfg.Port
		logger.Info("🌐 Server running on http://localhost" + addr)
		logger.Info("📚 Swagger docs: http://localhost" + addr + "/swagger/index.html")
		if err := router.Run(addr); err != nil {
			logger.Error(err, "Server error")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("👋 Shutting down server...")
}
