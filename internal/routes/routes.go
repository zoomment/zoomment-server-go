package routes

import (
	"github.com/gin-gonic/gin"

	"zoomment-server/internal/config"
	"zoomment-server/internal/handlers"
	"zoomment-server/internal/middleware"
)

// Setup configures all API routes
// This keeps main.go clean and routes organized
func Setup(router *gin.Engine, cfg *config.Config) {
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes group
	api := router.Group("/api")
	{
		// Comments routes
		setupCommentRoutes(api, cfg)

		// Users routes
		setupUserRoutes(api, cfg)

		// Sites routes
		setupSiteRoutes(api, cfg)

		// Reactions routes
		setupReactionRoutes(api)

		// Visitors routes
		setupVisitorRoutes(api)

		// Votes routes
		setupVoteRoutes(api)
	}
}

// setupCommentRoutes configures /api/comments routes
func setupCommentRoutes(api *gin.RouterGroup, cfg *config.Config) {
	comments := api.Group("/comments")
	{
		// Register both with and without trailing slash since RedirectTrailingSlash is disabled
		comments.GET("/", handlers.ListComments)
		comments.GET("", handlers.ListComments)
		comments.POST("/", handlers.AddComment(cfg))
		comments.POST("", handlers.AddComment(cfg))
		comments.DELETE("/:id", handlers.DeleteComment)
		// Load more replies for a specific comment
		comments.GET("/:commentId/replies", handlers.ListReplies)
		// Node.js uses access('admin') for this route
		comments.GET("/sites/:siteId", middleware.Access("admin"), handlers.ListCommentsBySite)
	}
}

// setupUserRoutes configures /api/users routes
func setupUserRoutes(api *gin.RouterGroup, cfg *config.Config) {
	users := api.Group("/users")
	{
		users.POST("/auth", handlers.AuthUser(cfg))
		users.GET("/profile", middleware.Access(), handlers.GetProfile)
		users.DELETE("/", middleware.Access(), handlers.DeleteUser)
	}
}

// setupSiteRoutes configures /api/sites routes
func setupSiteRoutes(api *gin.RouterGroup, cfg *config.Config) {
	sites := api.Group("/sites")
	{
		sites.GET("/", middleware.Access("admin"), handlers.ListSites)
		sites.POST("/", middleware.Access("admin"), handlers.AddSite(cfg))
		sites.DELETE("/:id", middleware.Access("admin"), handlers.DeleteSite)
	}
}

// setupReactionRoutes configures /api/reactions routes
func setupReactionRoutes(api *gin.RouterGroup) {
	reactions := api.Group("/reactions")
	{
		// Register both with and without trailing slash since RedirectTrailingSlash is disabled
		reactions.GET("/", handlers.ListReactions)
		reactions.GET("", handlers.ListReactions)
		reactions.POST("/", handlers.AddReaction)
		reactions.POST("", handlers.AddReaction)
	}
}

// setupVisitorRoutes configures /api/visitors routes
func setupVisitorRoutes(api *gin.RouterGroup) {
	visitors := api.Group("/visitors")
	{
		// Register both with and without trailing slash since RedirectTrailingSlash is disabled
		visitors.GET("/", handlers.GetVisitorCount)
		visitors.GET("", handlers.GetVisitorCount)
		visitors.POST("/", handlers.TrackVisitor)
		visitors.POST("", handlers.TrackVisitor)
		visitors.GET("/domain", handlers.GetVisitorsByDomain)
	}
}

// setupVoteRoutes configures /api/votes routes
func setupVoteRoutes(api *gin.RouterGroup) {
	votes := api.Group("/votes")
	{
		// Register both with and without trailing slash since RedirectTrailingSlash is disabled
		votes.POST("/", handlers.Vote)
		votes.POST("", handlers.Vote)
		votes.GET("/", handlers.GetVotesBulk)
		votes.GET("", handlers.GetVotesBulk)
		votes.GET("/:commentId", handlers.GetVote)
	}
}

