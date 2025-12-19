package repository

import (
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CommentWithReplies represents a comment with its replies embedded
// This is the result of our aggregation query
type CommentWithReplies struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	ParentID   *string            `bson:"parentId" json:"parentId,omitempty"`
	Author     string             `bson:"author" json:"author"`
	Email      string             `bson:"email" json:"-"` // Hidden from JSON
	Gravatar   string             `bson:"gravatar" json:"gravatar"`
	Body       string             `bson:"body" json:"body"`
	Domain     string             `bson:"domain" json:"-"`
	PageURL    string             `bson:"pageUrl" json:"-"`
	PageID     string             `bson:"pageId" json:"-"`
	IsVerified bool               `bson:"isVerified" json:"isVerified"`
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`

	// Replies are fetched via $lookup aggregation
	Replies []CommentWithReplies `bson:"replies" json:"replies"`

	// Owner for backward compatibility
	Owner *struct {
		Name     string `bson:"name" json:"name"`
		Gravatar string `bson:"gravatar" json:"gravatar"`
	} `bson:"owner" json:"owner,omitempty"`
}

// CommentPublicResponse is what we send to clients
// Matches Node.js API response format exactly
type CommentPublicResponse struct {
	ID         primitive.ObjectID      `json:"_id"`
	Owner      *struct {
		Name     string `json:"name"`
		Gravatar string `json:"gravatar"`
	} `json:"owner,omitempty"`
	IsOwn      bool                    `json:"isOwn"`
	Body       string                  `json:"body"`
	Author     string                  `json:"author"`
	Gravatar   string                  `json:"gravatar"`
	ParentID   *string                 `json:"parentId"` // Added for Node.js compatibility
	CreatedAt  time.Time               `json:"createdAt"`
	IsVerified bool                    `json:"isVerified"`
	Replies    []CommentPublicResponse `json:"replies,omitempty"`
}

// GetCommentsWithReplies fetches comments with their replies in a SINGLE query
// This solves the N+1 problem using MongoDB $lookup (like SQL JOIN)
func GetCommentsWithReplies(pageID, domain string) ([]CommentWithReplies, error) {
	// Build match condition
	matchCondition := bson.M{"parentId": nil}
	if pageID != "" {
		matchCondition["pageId"] = pageID
	} else if domain != "" {
		matchCondition["domain"] = domain
	}

	// MongoDB Aggregation Pipeline
	// This is like a series of data transformations
	pipeline := mongo.Pipeline{
		// Stage 1: Match top-level comments (no parent)
		{{Key: "$match", Value: matchCondition}},

		// Stage 2: Sort by creation date (newest first)
		{{Key: "$sort", Value: bson.D{{Key: "createdAt", Value: -1}}}},

		// Stage 3: Lookup (JOIN) replies from the same collection
		{{Key: "$lookup", Value: bson.M{
			"from": "comments", // Same collection (self-join)
			"let":  bson.M{"parentId": bson.M{"$toString": "$_id"}}, // Convert ObjectID to string
			"pipeline": mongo.Pipeline{
				// Match replies where parentId equals this comment's ID
				{{Key: "$match", Value: bson.M{
					"$expr": bson.M{"$eq": bson.A{"$parentId", "$$parentId"}},
				}}},
				// Sort replies by date (newest first)
				{{Key: "$sort", Value: bson.D{{Key: "createdAt", Value: -1}}}},
			},
			"as": "replies", // Output field name
		}}},
	}

	// Execute aggregation
	coll := mgm.Coll(&commentModel{})
	cursor, err := coll.Aggregate(mgm.Ctx(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(mgm.Ctx())

	// Decode results
	var comments []CommentWithReplies
	if err := cursor.All(mgm.Ctx(), &comments); err != nil {
		return nil, err
	}

	return comments, nil
}

// ToPublicResponse converts CommentWithReplies to public response
// Matches Node.js getCommentPublicData() output exactly
func (c *CommentWithReplies) ToPublicResponse(currentUserEmail string) CommentPublicResponse {
	isOwn := currentUserEmail != "" && currentUserEmail == c.Email

	response := CommentPublicResponse{
		ID:         c.ID,
		Author:     c.Author,
		Gravatar:   c.Gravatar,
		Body:       c.Body,
		ParentID:   c.ParentID, // Include parentId for Node.js compatibility
		IsVerified: c.IsVerified,
		IsOwn:      isOwn,
		CreatedAt:  c.CreatedAt,
	}

	// Add owner for backward compatibility (same as Node.js)
	response.Owner = &struct {
		Name     string `json:"name"`
		Gravatar string `json:"gravatar"`
	}{
		Name:     c.Author,
		Gravatar: c.Gravatar,
	}

	// Convert replies
	if len(c.Replies) > 0 {
		response.Replies = make([]CommentPublicResponse, 0, len(c.Replies))
		for _, reply := range c.Replies {
			response.Replies = append(response.Replies, reply.ToPublicResponse(currentUserEmail))
		}
	}

	return response
}

// commentModel is a helper struct for mgm.Coll()
type commentModel struct {
	mgm.DefaultModel `bson:",inline"`
}

func (c *commentModel) CollectionName() string {
	return "comments"
}

