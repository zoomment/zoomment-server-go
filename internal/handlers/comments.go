package handlers

import (
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"zoomment-server/internal/config"
	"zoomment-server/internal/errors"
	"zoomment-server/internal/logger"
	"zoomment-server/internal/middleware"
	"zoomment-server/internal/models"
	"zoomment-server/internal/repository"
	"zoomment-server/internal/services/mailer"
	"zoomment-server/internal/utils"
	"zoomment-server/internal/validators"
)

// ListComments returns comments for a page or domain
// GET /api/comments?pageId=xxx or ?domain=xxx
// Uses MongoDB aggregation to fetch comments WITH replies in a single query (solves N+1)
func ListComments(c *gin.Context) {
	pageID := c.Query("pageId")
	domain := c.Query("domain")

	// Validation: require at least one parameter
	if pageID == "" && domain == "" {
		errors.BadRequest("Bad request").Response(c)
		return
	}

	// Validate length limits
	if len(pageID) > 500 || len(domain) > 253 {
		errors.BadRequest("Bad request").Response(c)
		return
	}

	// Fetch comments with replies in a SINGLE query (no N+1 problem!)
	comments, err := repository.GetCommentsWithReplies(pageID, domain)
	if err != nil {
		logger.Error(err, "Failed to fetch comments")
		errors.ErrDatabaseError.Response(c)
		return
	}

	// Get current user email for "isOwn" flag
	user := middleware.GetUser(c)
	currentEmail := ""
	if user != nil {
		currentEmail = user.Email
	}

	// Convert to public response
	result := make([]repository.CommentPublicResponse, 0, len(comments))
	for _, comment := range comments {
		result = append(result, comment.ToPublicResponse(currentEmail))
	}

	c.JSON(http.StatusOK, result)
}

// AddComment creates a new comment
// POST /api/comments
func AddComment(cfg *config.Config) gin.HandlerFunc {
	// Create mailer instance
	mailService := mailer.New(cfg)

	return func(c *gin.Context) {
		var req validators.AddCommentRequest

		// Validate request body
		if err := c.ShouldBindJSON(&req); err != nil {
			errors.BadRequest("Invalid request body").Response(c)
			return
		}

		// Parse URL to get domain
		parsedURL, err := url.Parse(req.PageURL)
		if err != nil {
			errors.BadRequest("Invalid page URL").Response(c)
			return
		}

		// Sanitize inputs to prevent XSS
		email := utils.CleanEmail(req.Email)
		author := utils.SanitizeStrict(utils.CleanName(req.Author)) // Remove ALL HTML from name
		body := utils.SanitizeComment(req.Body)                      // Allow safe HTML in body

		// Get current user
		user := middleware.GetUser(c)
		isVerified := user != nil && user.Email == email

		// Create comment
		comment := &models.Comment{
			PageURL:    parsedURL.String(),
			PageID:     req.PageID,
			Domain:     parsedURL.Hostname(),
			Body:       body,
			Author:     author,
			Email:      email,
			Gravatar:   utils.GenerateGravatar(email),
			ParentID:   req.ParentID,
			IsVerified: isVerified,
			Secret:     utils.GenerateSecret(),
			Owner: &models.CommentOwner{
				Name:     author,
				Email:    email,
				Gravatar: utils.GenerateGravatar(email),
			},
		}

		err = mgm.Coll(comment).Create(comment)
		if err != nil {
			logger.Error(err, "Failed to create comment")
			errors.ErrDatabaseError.Response(c)
			return
		}

		// Return 200 OK (same as Node.js)
		c.JSON(http.StatusOK, comment)

		// Send email notifications asynchronously (don't block the response)
		go func() {
			// Send verification email to guest users (not authenticated)
			if user == nil {
				// Find or create user for the commenter
				commenterUser := &models.User{}
				err := mgm.Coll(commenterUser).First(bson.M{"email": email}, commenterUser)
				if err != nil {
					// Create new user
					commenterUser = models.NewUser(email)
					commenterUser.Name = author
					mgm.Coll(commenterUser).Create(commenterUser)
				}

				// Generate token for email verification
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"id":    commenterUser.ID.Hex(),
					"email": email,
					"name":  author,
					"exp":   time.Now().Add(365 * 24 * time.Hour).Unix(),
				})
				tokenString, _ := token.SignedString([]byte(cfg.JWTSecret))

				// Send verification email
				mailService.SendEmailVerification(email, tokenString, comment.PageURL)
			}

			// Send notification to site owner
			site := &models.Site{}
			err := mgm.Coll(site).First(bson.M{"domain": comment.Domain}, site)
			if err == nil {
				// Found the site, get the owner
				siteOwner := &models.User{}
				err := mgm.Coll(siteOwner).FindByID(site.UserID, siteOwner)
				if err == nil && siteOwner.Email != email {
					// Don't notify if the commenter is the site owner
					mailService.SendCommentNotification(siteOwner.Email, mailer.CommentData{
						Author:  author,
						Date:    comment.CreatedAt.Format("02 Jan 2006 - 15:04"),
						PageURL: comment.PageURL,
						Body:    body,
					})
				}
			}
		}()
	}
}

// DeleteComment removes a comment
// DELETE /api/comments/:id?secret=xxx
func DeleteComment(c *gin.Context) {
	commentID := c.Param("id")
	secret := c.Query("secret")

	// Convert string to ObjectID
	objID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		errors.BadRequest("Invalid comment ID").Response(c)
		return
	}

	// Build query
	query := bson.M{"_id": objID}

	// Check authorization
	user := middleware.GetUser(c)
	if secret != "" {
		// Guest deletion with secret
		query["secret"] = secret
	} else if user != nil {
		// Authenticated user deletion
		query["email"] = utils.CleanEmail(user.Email)
	} else {
		errors.ErrForbidden.Response(c)
		return
	}

	// Delete comment
	result, err := mgm.Coll(&models.Comment{}).DeleteOne(mgm.Ctx(), query)
	if err != nil {
		errors.ErrDatabaseError.Response(c)
		return
	}

	if result.DeletedCount == 0 {
		errors.NotFound("Comment").Response(c)
		return
	}

	// Return same format as Node.js
	c.JSON(http.StatusOK, gin.H{"_id": commentID})
}

// ListCommentsBySite returns all comments for a site
// GET /api/comments/sites/:siteId
func ListCommentsBySite(c *gin.Context) {
	siteID := c.Param("siteId")
	user := middleware.GetUser(c)

	// Convert string to ObjectID
	objID, err := primitive.ObjectIDFromHex(siteID)
	if err != nil {
		errors.NotFound("Site").Response(c)
		return
	}

	// Find the site
	site := &models.Site{}
	err = mgm.Coll(site).FindByID(objID, site)
	if err != nil {
		errors.NotFound("Site").Response(c)
		return
	}

	// Check ownership (same as Node.js)
	if site.UserID != user.ID {
		errors.NotFound("Site").Response(c)
		return
	}

	// Find all comments for this domain
	var comments []models.Comment
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	err = mgm.Coll(&models.Comment{}).SimpleFind(&comments, bson.M{"domain": site.Domain}, opts)
	if err != nil {
		errors.ErrDatabaseError.Response(c)
		return
	}

	// Return empty array instead of null
	if comments == nil {
		comments = []models.Comment{}
	}

	c.JSON(http.StatusOK, comments)
}
