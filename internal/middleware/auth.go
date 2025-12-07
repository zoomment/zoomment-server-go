package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"zoomment-server/internal/config"
	"zoomment-server/internal/models"
)

// Auth middleware extracts user from JWT token in header
// Similar to your auth() middleware in Express
func Auth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from header (like req.headers.token)
		tokenString := c.GetHeader("token")

		if tokenString == "" {
			// No token - continue as guest
			c.Next()
			return
		}

		// Parse and validate JWT token
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			// Invalid token - continue as guest
			c.Next()
			return
		}

		// Extract claims (payload) from token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Next()
			return
		}

		// Get user ID from claims
		userIDStr, ok := claims["id"].(string)
		if !ok {
			c.Next()
			return
		}

		// Convert string to ObjectID
		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			c.Next()
			return
		}

		// Find user in database
		user := &models.User{}
		err = mgm.Coll(user).FindByID(userID, user)
		if err != nil {
			c.Next()
			return
		}

		// Update user verification status if needed
		if !user.IsVerified {
			user.IsVerified = true
			mgm.Coll(user).Update(user)
		}

		// Store user in context (like req.user = user in Express)
		c.Set("user", user)
		c.Next()
	}
}

// GetUser retrieves the user from context
// Returns nil if no user is authenticated
func GetUser(c *gin.Context) *models.User {
	userInterface, exists := c.Get("user")
	if !exists {
		return nil
	}

	user, ok := userInterface.(*models.User)
	if !ok {
		return nil
	}

	return user
}

// Access middleware checks if user is authenticated and has required role
// Similar to your access() middleware in Express
func Access(level ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := GetUser(c)

		// No user - forbidden
		if user == nil {
			c.AbortWithStatusJSON(403, gin.H{"message": "Forbidden"})
			return
		}

		// If no level specified, just require authentication
		if len(level) == 0 {
			c.Next()
			return
		}

		// Check role level
		requiredLevel := level[0]

		// SuperAdmin (role=2) can access everything
		if user.Role == models.RoleSuperAdmin {
			c.Next()
			return
		}

		// Admin (role=1) can access "admin" level
		if requiredLevel == "admin" && user.Role == models.RoleAdmin {
			c.Next()
			return
		}

		// Otherwise, forbidden
		c.AbortWithStatusJSON(403, gin.H{"message": "Forbidden"})
	}
}
