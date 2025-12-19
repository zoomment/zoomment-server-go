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

// ListVisitors returns visitor count for a page and records the visit if fingerprint is provided
// GET /api/visitors?pageId=xxx
func ListVisitors(c *gin.Context) {
	pageID := c.Query("pageId")
	fingerprint := c.GetHeader("fingerprint")

	if pageID == "" {
		errors.BadRequest("pageId is required").Response(c)
		return
	}

	response, err := getPageVisitors(pageID, fingerprint)
	if err != nil {
		errors.ErrDatabaseError.Response(c)
		return
	}

	c.JSON(http.StatusOK, response)
}

// getPageVisitors fetches visitor count and records visit if fingerprint provided
func getPageVisitors(pageID, fingerprint string) (*models.VisitorCountResponse, error) {
	// Record visit if fingerprint is provided
	if fingerprint != "" {
		if err := recordVisit(pageID, fingerprint); err != nil {
			return nil, err
		}
	}

	// Count unique visitors for this page
	count, err := mgm.Coll(&models.Visitor{}).CountDocuments(mgm.Ctx(), bson.M{"pageId": pageID})
	if err != nil {
		return nil, err
	}

	return &models.VisitorCountResponse{
		Count: int(count),
	}, nil
}

// recordVisit records a visit for the given pageId and fingerprint (idempotent)
func recordVisit(pageID, fingerprint string) error {
	parsedURL, err := url.Parse("https://" + pageID)
	if err != nil {
		return err
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
		return mgm.Coll(newVisitor).Create(newVisitor)
	}

	// If visitor already exists or other error occurred, return it
	return err
}
