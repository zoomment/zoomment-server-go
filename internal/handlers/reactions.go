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

// AddReactionRequest is the request body for adding a reaction
type AddReactionRequest struct {
	PageID   string `json:"pageId" binding:"required"`
	Reaction string `json:"reaction" binding:"required"`
}

// ListReactions returns reactions for a page
// GET /api/reactions?pageId=xxx
func ListReactions(c *gin.Context) {
	pageID := c.Query("pageId")
	fingerprint := c.GetHeader("fingerprint")

	if pageID == "" {
		errors.BadRequest("pageId is required").Response(c)
		return
	}

	response, err := getPageReactions(pageID, fingerprint)
	if err != nil {
		errors.ErrDatabaseError.Response(c)
		return
	}

	c.JSON(http.StatusOK, response)
}

// AddReaction adds or toggles a reaction
// POST /api/reactions
func AddReaction(c *gin.Context) {
	var req AddReactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.BadRequest("Invalid request").Response(c)
		return
	}

	fingerprint := c.GetHeader("fingerprint")
	if fingerprint == "" {
		// Node.js uses 500 and sends plain text - keep for compatibility
		c.String(http.StatusInternalServerError, "Fingerprint required for reacting.")
		return
	}

	// Limit reaction length (just in case)
	reaction := req.Reaction
	if len(reaction) > 20 {
		reaction = reaction[:20]
	}

	// Parse domain from pageId
	parsedURL, err := url.Parse("https://" + req.PageID)
	if err != nil {
		errors.BadRequest("Invalid pageId").Response(c)
		return
	}
	domain := parsedURL.Hostname()

	// Find existing reaction
	existingReaction := &models.Reaction{}
	query := bson.M{
		"pageId":      req.PageID,
		"fingerprint": fingerprint,
		"domain":      domain,
	}

	err = mgm.Coll(existingReaction).First(query, existingReaction)

	if err == mongo.ErrNoDocuments {
		// No existing reaction - create new one
		newReaction := &models.Reaction{
			PageID:      req.PageID,
			Fingerprint: fingerprint,
			Domain:      domain,
			Reaction:    reaction,
		}
		mgm.Coll(newReaction).Create(newReaction)
	} else if err == nil {
		// Existing reaction found
		if existingReaction.Reaction == reaction {
			// Same reaction - remove it (toggle off)
			mgm.Coll(existingReaction).Delete(existingReaction)
		} else {
			// Different reaction - update it
			existingReaction.Reaction = reaction
			mgm.Coll(existingReaction).Update(existingReaction)
		}
	} else {
		errors.ErrDatabaseError.Response(c)
		return
	}

	// Return updated reactions
	response, err := getPageReactions(req.PageID, fingerprint)
	if err != nil {
		errors.ErrDatabaseError.Response(c)
		return
	}

	c.JSON(http.StatusOK, response)
}

// UserReactionResponse matches Node.js format: { reaction: "..." }
type UserReactionResponse struct {
	Reaction string `json:"reaction"`
}

// PageReactionsResponseNodeJS matches Node.js response format exactly
type PageReactionsResponseNodeJS struct {
	Aggregation  []models.ReactionAggregation `json:"aggregation"`
	UserReaction *UserReactionResponse        `json:"userReaction"`
}

// getPageReactions fetches reaction counts and user's reaction for a page
func getPageReactions(pageID, fingerprint string) (*PageReactionsResponseNodeJS, error) {
	// Aggregate reaction counts
	pipeline := []bson.M{
		{"$match": bson.M{"pageId": pageID}},
		{"$group": bson.M{
			"_id":   "$reaction",
			"count": bson.M{"$sum": 1},
		}},
	}

	cursor, err := mgm.Coll(&models.Reaction{}).Aggregate(mgm.Ctx(), pipeline)
	if err != nil {
		return nil, err
	}

	var aggregation []models.ReactionAggregation
	if err := cursor.All(mgm.Ctx(), &aggregation); err != nil {
		return nil, err
	}

	// Get user's reaction - Node.js returns { reaction: "..." } object
	var userReaction *UserReactionResponse
	if fingerprint != "" {
		reaction := &models.Reaction{}
		err := mgm.Coll(reaction).First(bson.M{
			"pageId":      pageID,
			"fingerprint": fingerprint,
		}, reaction)
		if err == nil {
			userReaction = &UserReactionResponse{Reaction: reaction.Reaction}
		}
	}

	return &PageReactionsResponseNodeJS{
		Aggregation:  aggregation,
		UserReaction: userReaction,
	}, nil
}

