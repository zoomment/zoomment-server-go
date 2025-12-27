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
	"zoomment-server/internal/constants"
	"zoomment-server/internal/errors"
	"zoomment-server/internal/logger"
	"zoomment-server/internal/middleware"
	"zoomment-server/internal/models"
	"zoomment-server/internal/repository"
	"zoomment-server/internal/services/mailer"
	"zoomment-server/internal/utils"
	"zoomment-server/internal/validators"
)

// ListComments returns comments for a page or domain with pagination
// GET /api/comments?pageId=xxx&limit=10&skip=0&sort=asc|desc
// Returns parent comments with repliesCount for each
func ListComments(c *gin.Context) {
	pageID := c.Query("pageId")
	domain := c.Query("domain")

	// Validation: require at least one parameter
	if pageID == "" && domain == "" {
		errors.BadRequest("pageId or domain is required").Response(c)
		return
	}

	// Validate length limits
	if len(pageID) > constants.MaxPageIDLength || len(domain) > constants.MaxDomainLength {
		errors.BadRequest("Bad request").Response(c)
		return
	}

	// Parse pagination parameters
	limit, skip := repository.ParsePagination(c.Query("limit"), c.Query("skip"))
	sortOrder := c.Query("sort")
	if sortOrder != "desc" {
		sortOrder = "asc" // Default to oldest first
	}

	// Fetch paginated comments with reply counts
	response, err := repository.GetPaginatedComments(pageID, domain, limit, skip, sortOrder)
	if err != nil {
		logger.Error(err, "Failed to fetch comments")
		errors.ErrDatabaseError.Response(c)
		return
	}

	c.JSON(http.StatusOK, response)
}

// ListReplies returns replies for a specific comment with pagination
// GET /api/comments/:commentId/replies?limit=10&skip=0
func ListReplies(c *gin.Context) {
	commentID := c.Param("commentId")

	// Validate comment ID format
	if _, err := primitive.ObjectIDFromHex(commentID); err != nil {
		errors.BadRequest("Invalid comment ID").Response(c)
		return
	}

	// Parse pagination parameters
	limit, skip := repository.ParsePagination(c.Query("limit"), c.Query("skip"))

	// Fetch replies
	response, err := repository.GetRepliesForComment(commentID, limit, skip)
	if err != nil {
		logger.Error(err, "Failed to fetch replies")
		errors.ErrDatabaseError.Response(c)
		return
	}

	c.JSON(http.StatusOK, response)
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
		}

		err = mgm.Coll(comment).Create(comment)
		if err != nil {
			logger.Error(err, "Failed to create comment")
			errors.ErrDatabaseError.Response(c)
			return
		}

		// Return 200 OK with _id instead of id
		c.JSON(http.StatusOK, CommentToResponse(comment))

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
					"exp":   time.Now().Add(constants.JWTExpirationHours * time.Hour).Unix(),
				})
				tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
				if err != nil {
					logger.Error(err, "Failed to sign JWT token for verification email")
					// Continue without verification email rather than failing
					return
				}

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
						Date:    comment.CreatedAt.Format(constants.DateFormat),
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

	c.JSON(http.StatusOK, NewDeletedResponse(commentID))
}

// ListCommentsBySite returns all comments for a site with pagination
// GET /api/comments/sites/:siteId?limit=10&skip=0
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

	// Parse pagination parameters
	limit, skip := repository.ParsePagination(c.Query("limit"), c.Query("skip"))

	// Get total count
	total, err := mgm.Coll(&models.Comment{}).CountDocuments(mgm.Ctx(), bson.M{"domain": site.Domain})
	if err != nil {
		errors.ErrDatabaseError.Response(c)
		return
	}

	// Find comments for this domain with pagination
	var comments []models.Comment
	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))
	err = mgm.Coll(&models.Comment{}).SimpleFind(&comments, bson.M{"domain": site.Domain}, opts)
	if err != nil {
		errors.ErrDatabaseError.Response(c)
		return
	}

	// Return empty array instead of null
	if comments == nil {
		comments = []models.Comment{}
	}

	c.JSON(http.StatusOK, NewPaginatedCommentsResponse(comments, total, limit, skip))
}
