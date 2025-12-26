package handlers

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"zoomment-server/internal/errors"
	"zoomment-server/internal/models"
)

// TrackVisitor records a page visit (requires fingerprint)
// POST /api/visitors
func TrackVisitor(c *gin.Context) {
	fingerprint := c.GetHeader("fingerprint")
	if fingerprint == "" {
		errors.BadRequest("Fingerprint required for tracking").Response(c)
		return
	}

	var pageID string = c.Query("pageId")

	if pageID == "" {
		errors.BadRequest("pageId is required").Response(c)
		return
	}
	
	// Extract domain from pageId
	parsedURL, err := url.Parse("https://" + pageID)
	if err != nil {
		errors.BadRequest("Invalid pageId").Response(c)
		return
	}
	domain := parsedURL.Hostname()

	// Check if visitor already exists
	existingVisitor := &models.Visitor{}
	query := bson.M{
		"pageId":      pageID,
		"fingerprint": fingerprint,
		"domain":      domain,
	}

	err = mgm.Coll(existingVisitor).First(query, existingVisitor)
	if err == mongo.ErrNoDocuments {
		// No existing visit - create new one
		newVisitor := &models.Visitor{
			PageID:      pageID,
			Fingerprint: fingerprint,
			Domain:      domain,
		}
		if err := mgm.Coll(newVisitor).Create(newVisitor); err != nil {
			errors.ErrDatabaseError.Response(c)
			return
		}
	} else if err != nil {
		errors.ErrDatabaseError.Response(c)
		return
	}

	// Get total unique visitors for this page
	count, err := mgm.Coll(&models.Visitor{}).CountDocuments(mgm.Ctx(), bson.M{"pageId": pageID})
	if err != nil {
		errors.ErrDatabaseError.Response(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pageId": pageID,
		"count":  count,
	})
}

// GetVisitorCount returns visitor count for a page (no tracking)
// GET /api/visitors?pageId=xxx
func GetVisitorCount(c *gin.Context) {
	pageID := c.Query("pageId")

	if pageID == "" {
		errors.BadRequest("pageId is required").Response(c)
		return
	}

	// Count unique visitors for this page
	count, err := mgm.Coll(&models.Visitor{}).CountDocuments(mgm.Ctx(), bson.M{"pageId": pageID})
	if err != nil {
		errors.ErrDatabaseError.Response(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pageId": pageID,
		"count":  count,
	})
}

// GetVisitorsByDomain returns page view counts grouped by pageId for a domain
// GET /api/visitors/domain?domain=xxx
func GetVisitorsByDomain(c *gin.Context) {
	domain := c.Query("domain")

	if domain == "" {
		errors.BadRequest("domain is required").Response(c)
		return
	}

	// Get page view counts grouped by pageId, sorted by most viewed first
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"domain": domain}}},
		{{Key: "$group", Value: bson.M{
			"_id":   "$pageId",
			"count": bson.M{"$sum": 1},
		}}},
		{{Key: "$sort", Value: bson.M{"count": -1}}},
		{{Key: "$project", Value: bson.M{
			"pageId": "$_id",
			"count":  1,
			"_id":    0,
		}}},
	}

	cursor, err := mgm.Coll(&models.Visitor{}).Aggregate(mgm.Ctx(), pipeline)
	if err != nil {
		errors.ErrDatabaseError.Response(c)
		return
	}
	defer cursor.Close(mgm.Ctx())

	var pages []struct {
		PageID string `bson:"pageId" json:"pageId"`
		Count  int    `bson:"count" json:"count"`
	}

	if err := cursor.All(mgm.Ctx(), &pages); err != nil {
		errors.ErrDatabaseError.Response(c)
		return
	}

	// Return empty array if no pages found
	if pages == nil {
		pages = []struct {
			PageID string `bson:"pageId" json:"pageId"`
			Count  int    `bson:"count" json:"count"`
		}{}
	}

	c.JSON(http.StatusOK, gin.H{
		"domain": domain,
		"pages":  pages,
	})
}
