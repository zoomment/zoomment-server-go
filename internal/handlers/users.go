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
	"zoomment-server/internal/constants"
	"zoomment-server/internal/errors"
	"zoomment-server/internal/logger"
	"zoomment-server/internal/middleware"
	"zoomment-server/internal/models"
	"zoomment-server/internal/services/mailer"
	"zoomment-server/internal/utils"
	"zoomment-server/internal/validators"
)

// AuthUser handles magic link authentication
// POST /api/users/auth
func AuthUser(cfg *config.Config) gin.HandlerFunc {
	// Create mailer instance
	mailService := mailer.New(cfg)

	return func(c *gin.Context) {
		var req validators.AuthRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			errors.BadRequest("Invalid email").Response(c)
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
					errors.ErrDatabaseError.Response(c)
					return
				}
			} else {
				logger.Error(err, "Database error finding user")
				errors.ErrDatabaseError.Response(c)
				return
			}
		}

		// Generate JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":    user.ID.Hex(),
			"email": user.Email,
			"name":  user.Name,
			"exp":   time.Now().Add(constants.JWTExpirationHours * time.Hour).Unix(), // 1 year
		})

		tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			logger.Error(err, "Failed to generate token")
			errors.ErrInternalError.Response(c)
			return
		}

		// Send magic link email (async - don't block the response)
		go func() {
			if err := mailService.SendMagicLink(email, tokenString); err != nil {
				// Error is already logged in mailer.SendMagicLink
				// We don't fail the request, just log the error
			}
		}()

		c.JSON(http.StatusOK, MessageResponse{Message: "Magic link sent to your email"})
	}
}

// GetProfile returns the current user's profile
// GET /api/users/profile
func GetProfile(c *gin.Context) {
	user := middleware.GetUser(c)

	if user == nil {
		errors.ErrForbidden.Response(c)
		return
	}

	c.JSON(http.StatusOK, UserProfileResponse{
		ID:    user.ID.Hex(),
		Name:  user.Name,
		Email: user.Email,
	})
}

// DeleteUser deletes the current user and their sites
// DELETE /api/users
func DeleteUser(c *gin.Context) {
	user := middleware.GetUser(c)

	if user == nil {
		errors.ErrForbidden.Response(c)
		return
	}

	userID := user.ID.Hex()

	// Delete user first (same order as Node.js)
	result, err := mgm.Coll(&models.User{}).DeleteOne(mgm.Ctx(), bson.M{"_id": user.ID})
	if err != nil {
		logger.Error(err, "Failed to delete user")
		errors.ErrDatabaseError.Response(c)
		return
	}

	// Delete user's sites
	mgm.Coll(&models.Site{}).DeleteMany(mgm.Ctx(), bson.M{"userId": user.ID})

	if result.DeletedCount == 0 {
		errors.NotFound("Account").Response(c)
		return
	}

	c.JSON(http.StatusOK, NewDeletedResponse(userID))
}
