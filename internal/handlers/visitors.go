package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"zoomment-server/internal/errors"
	"zoomment-server/internal/models"
)

// TrackVisitor records a page visit (requires fingerprint)
// POST /api/visitors?pageId=xxx
func TrackVisitor(c *gin.Context) {
	fingerprint := c.GetHeader("fingerprint")
	if fingerprint == "" {
		errors.BadRequest("Fingerprint required for tracking").Response(c)
		return
	}

	pageID := c.Query("pageId")
	if pageID == "" {
		errors.BadRequest("pageId is required").Response(c)
		return
	}

	domain, err := ExtractDomainFromPageID(pageID)
	if err != nil {
		errors.BadRequest("Invalid pageId").Response(c)
		return
	}

	// Upsert visitor (create if not exists)
	query := bson.M{
		"pageId":      pageID,
		"fingerprint": fingerprint,
		"domain":      domain,
	}

	existingVisitor := &models.Visitor{}
	if err := mgm.Coll(existingVisitor).First(query, existingVisitor); err == mongo.ErrNoDocuments {
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

	count, err := mgm.Coll(&models.Visitor{}).CountDocuments(mgm.Ctx(), bson.M{"pageId": pageID})
	if err != nil {
		errors.ErrDatabaseError.Response(c)
		return
	}

	c.JSON(http.StatusOK, NewVisitorCountResponse(pageID, count))
}

// GetVisitorCount returns visitor count for a page
// GET /api/visitors?pageId=xxx
func GetVisitorCount(c *gin.Context) {
	pageID := c.Query("pageId")
	if pageID == "" {
		errors.BadRequest("pageId is required").Response(c)
		return
	}

	count, err := mgm.Coll(&models.Visitor{}).CountDocuments(mgm.Ctx(), bson.M{"pageId": pageID})
	if err != nil {
		errors.ErrDatabaseError.Response(c)
		return
	}

	c.JSON(http.StatusOK, NewVisitorCountResponse(pageID, count))
}

// GetVisitorsByDomain returns page view counts grouped by pageId for a domain
// GET /api/visitors/domain?domain=xxx
func GetVisitorsByDomain(c *gin.Context) {
	domain := c.Query("domain")
	if domain == "" {
		errors.BadRequest("domain is required").Response(c)
		return
	}

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

	var pages []PageViewCount
	if err := cursor.All(mgm.Ctx(), &pages); err != nil {
		errors.ErrDatabaseError.Response(c)
		return
	}

	// Ensure empty array instead of null
	if pages == nil {
		pages = []PageViewCount{}
	}

	c.JSON(http.StatusOK, NewDomainPagesResponse(domain, pages))
}
