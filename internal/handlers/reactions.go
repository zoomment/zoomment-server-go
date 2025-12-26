package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"zoomment-server/internal/constants"
	"zoomment-server/internal/errors"
	"zoomment-server/internal/models"
	"zoomment-server/internal/validators"
)

// ========================================
// Response Types
// ========================================

// UserReactionResponse matches Node.js format: { reaction: "..." }
type UserReactionResponse struct {
	Reaction string `json:"reaction"`
}

// PageReactionsResponse matches Node.js response format exactly
type PageReactionsResponse struct {
	Aggregation  []models.ReactionAggregation `json:"aggregation"`
	UserReaction *UserReactionResponse        `json:"userReaction"`
}

// ========================================
// Handlers
// ========================================

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
	var req validators.AddReactionRequest
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

	// Limit reaction length
	reaction := req.Reaction
	if len(reaction) > constants.MaxReactionLength {
		reaction = reaction[:constants.MaxReactionLength]
	}

	domain, err := ExtractDomainFromPageID(req.PageID)
	if err != nil {
		errors.BadRequest("Invalid pageId").Response(c)
		return
	}

	// Process reaction (create/update/delete)
	if err := processReaction(req.PageID, fingerprint, domain, reaction); err != nil {
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

// ========================================
// Helper Functions
// ========================================

// processReaction handles the reaction creation/update/deletion logic
func processReaction(pageID, fingerprint, domain, reaction string) error {
	query := bson.M{
		"pageId":      pageID,
		"fingerprint": fingerprint,
		"domain":      domain,
	}

	existingReaction := &models.Reaction{}
	err := mgm.Coll(existingReaction).First(query, existingReaction)

	if err == mongo.ErrNoDocuments {
		// Create new reaction
		newReaction := &models.Reaction{
			PageID:      pageID,
			Fingerprint: fingerprint,
			Domain:      domain,
			Reaction:    reaction,
		}
		return mgm.Coll(newReaction).Create(newReaction)
	}

	if err != nil {
		return err
	}

	if existingReaction.Reaction == reaction {
		// Same reaction - remove it (toggle off)
		return mgm.Coll(existingReaction).Delete(existingReaction)
	}

	// Different reaction - update it
	existingReaction.Reaction = reaction
	return mgm.Coll(existingReaction).Update(existingReaction)
}

// getPageReactions fetches reaction counts and user's reaction for a page
func getPageReactions(pageID, fingerprint string) (*PageReactionsResponse, error) {
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

	// Get user's reaction
	var userReaction *UserReactionResponse
	if fingerprint != "" {
		reaction := &models.Reaction{}
		if err := mgm.Coll(reaction).First(bson.M{
			"pageId":      pageID,
			"fingerprint": fingerprint,
		}, reaction); err == nil {
			userReaction = &UserReactionResponse{Reaction: reaction.Reaction}
		}
	}

	return &PageReactionsResponse{
		Aggregation:  aggregation,
		UserReaction: userReaction,
	}, nil
}
