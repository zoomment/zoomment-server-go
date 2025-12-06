package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"zoomment-server/internal/config"
	"zoomment-server/internal/logger"
	"zoomment-server/internal/middleware"
	"zoomment-server/internal/models"
	"zoomment-server/internal/services/mailer"
	"zoomment-server/internal/utils"
)

// AuthRequest is the request body for authentication
type AuthRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// AuthUser handles magic link authentication
// POST /api/users/auth
func AuthUser(cfg *config.Config) gin.HandlerFunc {
	// Create mailer instance
	mailService := mailer.New(cfg)

	return func(c *gin.Context) {
		var req AuthRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email"})
			return
		}

		email := utils.CleanEmail(req.Email)

		// Find or create user
		user := &models.User{}
		err := mgm.Coll(user).First(bson.M{"email": email}, user)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				// Create new user
				user = models.NewUser(email)
				if err := mgm.Coll(user).Create(user); err != nil {
					logger.Error(err, "Failed to create user")
					c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user"})
					return
				}
			} else {
				logger.Error(err, "Database error finding user")
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
				return
			}
		}

		// Generate JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":    user.ID.Hex(),
			"email": user.Email,
			"name":  user.Name,
			"exp":   time.Now().Add(365 * 24 * time.Hour).Unix(), // 1 year
		})

		tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			logger.Error(err, "Failed to generate token")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate token"})
			return
		}

		// Send magic link email (async - don't block the response)
		go func() {
			if err := mailService.SendMagicLink(email, tokenString); err != nil {
				// Error is already logged in mailer.SendMagicLink
				// We don't fail the request, just log the error
			}
		}()

		// Return empty object (same as Node.js)
		c.JSON(http.StatusOK, gin.H{})
	}
}

// GetProfile returns the current user's profile
// GET /api/users/profile
func GetProfile(c *gin.Context) {
	user := middleware.GetUser(c)

	if user == nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Forbidden"})
		return
	}

	// Match Node.js response format exactly
	c.JSON(http.StatusOK, gin.H{
		"name":  user.Name,
		"email": user.Email,
		"id":    user.ID.Hex(),
	})
}

// DeleteUser deletes the current user and their sites
// DELETE /api/users
func DeleteUser(c *gin.Context) {
	user := middleware.GetUser(c)

	if user == nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Forbidden"})
		return
	}

	userID := user.ID.Hex()

	// Delete user first (same order as Node.js)
	result, err := mgm.Coll(&models.User{}).DeleteOne(mgm.Ctx(), bson.M{"_id": user.ID})
	if err != nil {
		logger.Error(err, "Failed to delete user")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete user"})
		return
	}

	// Delete user's sites
	mgm.Coll(&models.Site{}).DeleteMany(mgm.Ctx(), bson.M{"userId": user.ID})

	if result.DeletedCount > 0 {
		c.JSON(http.StatusOK, gin.H{"_id": userID})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"message": "Account not found"})
	}
}
