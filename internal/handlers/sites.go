package handlers

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"zoomment-server/internal/config"
	"zoomment-server/internal/errors"
	"zoomment-server/internal/middleware"
	"zoomment-server/internal/models"
	"zoomment-server/internal/services/metadata"
	"zoomment-server/internal/validators"
)

// ListSites returns all sites for the current user
// GET /api/sites
func ListSites(c *gin.Context) {
	user := middleware.GetUser(c)

	var sites []models.Site
	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}})
	err := mgm.Coll(&models.Site{}).SimpleFind(&sites, bson.M{"userId": user.ID}, opts)
	if err != nil {
		errors.ErrDatabaseError.Response(c)
		return
	}

	c.JSON(http.StatusOK, sites)
}

// AddSite registers a new site after verifying ownership via meta tag
// POST /api/sites
func AddSite(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req validators.AddSiteRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			errors.BadRequest("Invalid URL").Response(c)
			return
		}

		user := middleware.GetUser(c)

		// Parse URL
		parsedURL, err := url.Parse(req.URL)
		if err != nil {
			errors.BadRequest("Invalid URL").Response(c)
			return
		}

		domain := parsedURL.Hostname()

		// Fetch and verify the zoomment meta tag
		token, err := metadata.FetchSiteToken(parsedURL.String())
		if err != nil || token != user.ID.Hex() {
			errors.NotFound("Meta tag").Response(c)
			return
		}

		// Check if site already exists
		existingSite := &models.Site{}
		err = mgm.Coll(existingSite).First(bson.M{"domain": domain}, existingSite)
		if err == nil {
			// Node.js uses 401 here (unusual but we match it for compatibility)
			errors.Conflict("Website already exists").Response(c)
			return
		}
		if err != mongo.ErrNoDocuments {
			errors.ErrDatabaseError.Response(c)
			return
		}

		// Create new site
		site := &models.Site{
			UserID:   user.ID,
			Domain:   domain,
			Verified: true,
		}

		err = mgm.Coll(site).Create(site)
		if err != nil {
			errors.ErrDatabaseError.Response(c)
			return
		}

		c.JSON(http.StatusOK, site)
	}
}

// DeleteSite removes a site
// DELETE /api/sites/:id
func DeleteSite(c *gin.Context) {
	siteID := c.Param("id")
	user := middleware.GetUser(c)

	// Convert string to ObjectID
	objID, err := primitive.ObjectIDFromHex(siteID)
	if err != nil {
		errors.BadRequest("Invalid site ID").Response(c)
		return
	}

	// Delete only if user owns the site
	result, err := mgm.Coll(&models.Site{}).DeleteOne(mgm.Ctx(), bson.M{
		"_id":    objID,
		"userId": user.ID,
	})
	if err != nil {
		errors.ErrDatabaseError.Response(c)
		return
	}

	if result.DeletedCount == 0 {
		errors.NotFound("Site").Response(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"_id": siteID})
}

