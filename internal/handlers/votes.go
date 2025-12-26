package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"zoomment-server/internal/errors"
	"zoomment-server/internal/models"
)

// VoteRequest represents the request body for voting
type VoteRequest struct {
	CommentID string `json:"commentId" binding:"required"`
	Value     int    `json:"value" binding:"required,oneof=1 -1"`
}

// Vote handles voting on a comment (upvote/downvote)
// POST /api/votes
// - If no vote exists: create vote
// - If same vote exists: remove vote (toggle off)
// - If opposite vote exists: update vote
func Vote(c *gin.Context) {
	fingerprint := c.GetHeader("fingerprint")
	if fingerprint == "" {
		errors.BadRequest("Fingerprint required for voting").Response(c)
		return
	}

	var req VoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.BadRequest("commentId is required and value must be 1 or -1").Response(c)
		return
	}

	// Verify comment exists
	commentObjID, err := primitive.ObjectIDFromHex(req.CommentID)
	if err != nil {
		errors.BadRequest("Invalid comment ID").Response(c)
		return
	}

	comment := &models.Comment{}
	if err := mgm.Coll(comment).FindByID(commentObjID, comment); err != nil {
		errors.NotFound("Comment").Response(c)
		return
	}

	// Handle vote logic
	if err := processVote(req.CommentID, fingerprint, req.Value); err != nil {
		errors.ErrDatabaseError.Response(c)
		return
	}

	// Return updated vote counts
	response, err := calculateVoteCounts(req.CommentID, fingerprint)
	if err != nil {
		errors.ErrDatabaseError.Response(c)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetVote returns vote counts for a single comment
// GET /api/votes/:commentId
func GetVote(c *gin.Context) {
	commentID := c.Param("commentId")
	fingerprint := c.GetHeader("fingerprint")

	response, err := calculateVoteCounts(commentID, fingerprint)
	if err != nil {
		errors.ErrDatabaseError.Response(c)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetVotesBulk returns vote counts for multiple comments
// GET /api/votes?commentIds=id1,id2,id3
func GetVotesBulk(c *gin.Context) {
	commentIdsParam := c.Query("commentIds")
	fingerprint := c.GetHeader("fingerprint")

	if commentIdsParam == "" {
		errors.BadRequest("commentIds is required").Response(c)
		return
	}

	// Parse and filter comment IDs
	filteredIds := parseCommentIds(commentIdsParam)
	if len(filteredIds) == 0 {
		errors.BadRequest("commentIds is required").Response(c)
		return
	}

	// Fetch all votes for these comments
	var votes []models.Vote
	if err := mgm.Coll(&models.Vote{}).SimpleFind(&votes, bson.M{
		"commentId": bson.M{"$in": filteredIds},
	}); err != nil {
		errors.ErrDatabaseError.Response(c)
		return
	}

	// Calculate counts per comment
	result := make(map[string]models.VoteResponse)
	for _, id := range filteredIds {
		result[id] = aggregateVotes(id, votes, fingerprint)
	}

	c.JSON(http.StatusOK, result)
}

// ========================================
// Helper Functions
// ========================================

// processVote handles the vote creation/update/deletion logic
func processVote(commentID, fingerprint string, value int) error {
	query := bson.M{
		"commentId":   commentID,
		"fingerprint": fingerprint,
	}

	existingVote := &models.Vote{}
	err := mgm.Coll(existingVote).First(query, existingVote)

	if err == mongo.ErrNoDocuments {
		// Create new vote
		newVote := &models.Vote{
			CommentID:   commentID,
			Fingerprint: fingerprint,
			Value:       value,
		}
		return mgm.Coll(newVote).Create(newVote)
	}

	if err != nil {
		return err
	}

	if existingVote.Value == value {
		// Same vote - remove it (toggle off)
		_, err = mgm.Coll(existingVote).DeleteOne(mgm.Ctx(), bson.M{"_id": existingVote.ID})
		return err
	}

	// Different vote - update it
	existingVote.Value = value
	return mgm.Coll(existingVote).Update(existingVote)
}

// calculateVoteCounts fetches and calculates vote counts for a comment
func calculateVoteCounts(commentID, fingerprint string) (*models.VoteResponse, error) {
	var votes []models.Vote
	if err := mgm.Coll(&models.Vote{}).SimpleFind(&votes, bson.M{"commentId": commentID}); err != nil {
		return nil, err
	}

	response := aggregateVotes(commentID, votes, fingerprint)
	return &response, nil
}

// aggregateVotes calculates vote counts from a slice of votes
func aggregateVotes(commentID string, votes []models.Vote, fingerprint string) models.VoteResponse {
	var upvotes, downvotes, userVote int

	for _, vote := range votes {
		if vote.CommentID != commentID {
			continue
		}
		if vote.Value == 1 {
			upvotes++
		} else if vote.Value == -1 {
			downvotes++
		}
		if fingerprint != "" && vote.Fingerprint == fingerprint {
			userVote = vote.Value
		}
	}

	return models.VoteResponse{
		CommentID: commentID,
		Upvotes:   upvotes,
		Downvotes: downvotes,
		Score:     upvotes - downvotes,
		UserVote:  userVote,
	}
}

// parseCommentIds splits and trims the comma-separated comment IDs
func parseCommentIds(param string) []string {
	var result []string
	for _, id := range strings.Split(param, ",") {
		if id = strings.TrimSpace(id); id != "" {
			result = append(result, id)
		}
	}
	return result
}
