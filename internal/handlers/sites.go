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
	"zoomment-server/internal/middleware"
	"zoomment-server/internal/models"
	"zoomment-server/internal/services/metadata"
)

// AddSiteRequest is the request body for adding a site
type AddSiteRequest struct {
	URL string `json:"url" binding:"required,url"`
}

// ListSites returns all sites for the current user
// GET /api/sites
func ListSites(c *gin.Context) {
	user := middleware.GetUser(c)

	var sites []models.Site
	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}})
	err := mgm.Coll(&models.Site{}).SimpleFind(&sites, bson.M{"userId": user.ID}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch sites"})
		return
	}

	c.JSON(http.StatusOK, sites)
}

// AddSite registers a new site after verifying ownership via meta tag
// POST /api/sites
func AddSite(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AddSiteRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid URL"})
			return
		}

		user := middleware.GetUser(c)

		// Parse URL
		parsedURL, err := url.Parse(req.URL)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid URL"})
			return
		}

		domain := parsedURL.Hostname()

		// Fetch and verify the zoomment meta tag
		token, err := metadata.FetchSiteToken(parsedURL.String())
		if err != nil || token != user.ID.Hex() {
			c.JSON(http.StatusNotFound, gin.H{"message": "Meta tag not found"})
			return
		}

		// Check if site already exists
		existingSite := &models.Site{}
		err = mgm.Coll(existingSite).First(bson.M{"domain": domain}, existingSite)
		if err == nil {
			// Node.js uses 401 here (unusual but we match it for compatibility)
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Website already exists"})
			return
		}
		if err != mongo.ErrNoDocuments {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
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
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create site"})
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
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid site ID"})
		return
	}

	// Delete only if user owns the site
	result, err := mgm.Coll(&models.Site{}).DeleteOne(mgm.Ctx(), bson.M{
		"_id":    objID,
		"userId": user.ID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete site"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Site not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"_id": siteID})
}

