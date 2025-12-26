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
// If no vote exists: create vote
// If same vote exists: remove vote (toggle off)
// If opposite vote exists: update vote
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

	// Validate value
	if req.Value != 1 && req.Value != -1 {
		errors.BadRequest("value must be 1 (upvote) or -1 (downvote)").Response(c)
		return
	}

	// Check if comment exists
	commentObjID, err := primitive.ObjectIDFromHex(req.CommentID)
	if err != nil {
		errors.BadRequest("Invalid comment ID").Response(c)
		return
	}

	comment := &models.Comment{}
	err = mgm.Coll(comment).FindByID(commentObjID, comment)
	if err != nil {
		errors.NotFound("Comment").Response(c)
		return
	}

	// Find existing vote
	existingVote := &models.Vote{}
	err = mgm.Coll(existingVote).First(bson.M{
		"commentId":   req.CommentID,
		"fingerprint": fingerprint,
	}, existingVote)

	if err == mongo.ErrNoDocuments {
		// Create new vote
		newVote := &models.Vote{
			CommentID:   req.CommentID,
			Fingerprint: fingerprint,
			Value:       req.Value,
		}
		if err := mgm.Coll(newVote).Create(newVote); err != nil {
			errors.ErrDatabaseError.Response(c)
			return
		}
	} else if err == nil {
		if existingVote.Value == req.Value {
			// Same vote - remove it (toggle off)
			if _, err := mgm.Coll(existingVote).DeleteOne(mgm.Ctx(), bson.M{"_id": existingVote.ID}); err != nil {
				errors.ErrDatabaseError.Response(c)
				return
			}
		} else {
			// Different vote - update it
			existingVote.Value = req.Value
			if err := mgm.Coll(existingVote).Update(existingVote); err != nil {
				errors.ErrDatabaseError.Response(c)
				return
			}
		}
	} else {
		errors.ErrDatabaseError.Response(c)
		return
	}

	// Get updated vote counts
	response, err := getVoteResponse(req.CommentID, fingerprint)
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

	response, err := getVoteResponse(commentID, fingerprint)
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

	commentIds := strings.Split(commentIdsParam, ",")
	if len(commentIds) == 0 {
		errors.BadRequest("commentIds is required").Response(c)
		return
	}

	// Filter out empty strings
	var filteredIds []string
	for _, id := range commentIds {
		if id = strings.TrimSpace(id); id != "" {
			filteredIds = append(filteredIds, id)
		}
	}

	if len(filteredIds) == 0 {
		errors.BadRequest("commentIds is required").Response(c)
		return
	}

	// Get all votes for these comments
	var votes []models.Vote
	err := mgm.Coll(&models.Vote{}).SimpleFind(&votes, bson.M{
		"commentId": bson.M{"$in": filteredIds},
	})
	if err != nil {
		errors.ErrDatabaseError.Response(c)
		return
	}

	// Calculate counts per comment
	result := make(map[string]models.VoteResponse)

	for _, id := range filteredIds {
		var upvotes, downvotes, userVote int

		for _, vote := range votes {
			if vote.CommentID == id {
				if vote.Value == 1 {
					upvotes++
				} else if vote.Value == -1 {
					downvotes++
				}
				if fingerprint != "" && vote.Fingerprint == fingerprint {
					userVote = vote.Value
				}
			}
		}

		result[id] = models.VoteResponse{
			CommentID: id,
			Upvotes:   upvotes,
			Downvotes: downvotes,
			Score:     upvotes - downvotes,
			UserVote:  userVote,
		}
	}

	c.JSON(http.StatusOK, result)
}

// getVoteResponse calculates vote counts for a comment
func getVoteResponse(commentID, fingerprint string) (*models.VoteResponse, error) {
	var votes []models.Vote
	err := mgm.Coll(&models.Vote{}).SimpleFind(&votes, bson.M{"commentId": commentID})
	if err != nil {
		return nil, err
	}

	var upvotes, downvotes, userVote int
	for _, vote := range votes {
		if vote.Value == 1 {
			upvotes++
		} else if vote.Value == -1 {
			downvotes++
		}
		if fingerprint != "" && vote.Fingerprint == fingerprint {
			userVote = vote.Value
		}
	}

	return &models.VoteResponse{
		CommentID: commentID,
		Upvotes:   upvotes,
		Downvotes: downvotes,
		Score:     upvotes - downvotes,
		UserVote:  userVote,
	}, nil
}

