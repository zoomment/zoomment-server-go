package repository

import (
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DefaultLimit = 10
	MaxLimit     = 50
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
	IsOwn        bool                    `json:"isOwn"`
	Body         string                  `json:"body"`
	Author       string                  `json:"author"`
	Gravatar     string                  `json:"gravatar"`
	ParentID     *string                 `json:"parentId"`
	CreatedAt    time.Time               `json:"createdAt"`
	IsVerified   bool                    `json:"isVerified"`
	RepliesCount int                     `json:"repliesCount,omitempty"`
	Replies      []CommentPublicResponse `json:"replies,omitempty"`
}

// PaginatedCommentsResponse is the response format for paginated comments list
type PaginatedCommentsResponse struct {
	Comments []CommentPublicResponse `json:"comments"`
	Total    int64                   `json:"total"`
	Limit    int                     `json:"limit"`
	Skip     int                     `json:"skip"`
	HasMore  bool                    `json:"hasMore"`
}

// PaginatedRepliesResponse is the response format for paginated replies list
type PaginatedRepliesResponse struct {
	Replies []CommentPublicResponse `json:"replies"`
	Total   int64                   `json:"total"`
	Limit   int                     `json:"limit"`
	Skip    int                     `json:"skip"`
	HasMore bool                    `json:"hasMore"`
}

// ParsePagination parses limit and skip from query parameters
func ParsePagination(limitStr, skipStr string) (limit, skip int) {
	limit = DefaultLimit
	skip = 0

	if limitStr != "" {
		if parsed := parseInt(limitStr); parsed > 0 {
			limit = parsed
			if limit > MaxLimit {
				limit = MaxLimit
			}
		}
	}

	if skipStr != "" {
		if parsed := parseInt(skipStr); parsed >= 0 {
			skip = parsed
		}
	}

	return limit, skip
}

// parseInt safely parses a string to int
func parseInt(s string) int {
	var result int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + int(c-'0')
		} else {
			return 0
		}
	}
	return result
}

// GetPaginatedComments fetches parent comments with pagination and reply counts
func GetPaginatedComments(pageID, domain string, limit, skip int, sortOrder string) (*PaginatedCommentsResponse, error) {
	// Build match condition for parent comments only
	matchCondition := bson.M{"parentId": nil}
	if pageID != "" {
		matchCondition["pageId"] = pageID
	} else if domain != "" {
		matchCondition["domain"] = domain
	}

	// Determine sort order
	sortDirection := 1 // asc (oldest first) is default
	if sortOrder == "desc" {
		sortDirection = -1
	}

	coll := mgm.Coll(&commentModel{})

	// Get total count
	total, err := coll.CountDocuments(mgm.Ctx(), matchCondition)
	if err != nil {
		return nil, err
	}

	// Get parent comments with pagination
	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: sortDirection}}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	cursor, err := coll.Find(mgm.Ctx(), matchCondition, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(mgm.Ctx())

	var comments []CommentWithReplies
	if err := cursor.All(mgm.Ctx(), &comments); err != nil {
		return nil, err
	}

	// Get reply counts for each comment
	if len(comments) > 0 {
		commentIds := make([]string, len(comments))
		for i, c := range comments {
			commentIds[i] = c.ID.Hex()
		}

		replyCounts, err := getReplyCountsForComments(commentIds)
		if err != nil {
			return nil, err
		}

		// Build response with reply counts
		result := make([]CommentPublicResponse, 0, len(comments))
		for _, comment := range comments {
			response := comment.ToPublicResponseWithoutReplies("")
			response.RepliesCount = replyCounts[comment.ID.Hex()]
			result = append(result, response)
		}

		return &PaginatedCommentsResponse{
			Comments: result,
			Total:    total,
			Limit:    limit,
			Skip:     skip,
			HasMore:  int64(skip+len(comments)) < total,
		}, nil
	}

	return &PaginatedCommentsResponse{
		Comments: []CommentPublicResponse{},
		Total:    total,
		Limit:    limit,
		Skip:     skip,
		HasMore:  false,
	}, nil
}

// GetRepliesForComment fetches replies for a specific comment with pagination
func GetRepliesForComment(commentID string, limit, skip int) (*PaginatedRepliesResponse, error) {
	coll := mgm.Coll(&commentModel{})
	matchCondition := bson.M{"parentId": commentID}

	// Get total count
	total, err := coll.CountDocuments(mgm.Ctx(), matchCondition)
	if err != nil {
		return nil, err
	}

	// Get replies with pagination (always sort by oldest first for replies)
	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: 1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	cursor, err := coll.Find(mgm.Ctx(), matchCondition, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(mgm.Ctx())

	var replies []CommentWithReplies
	if err := cursor.All(mgm.Ctx(), &replies); err != nil {
		return nil, err
	}

	// Build response
	result := make([]CommentPublicResponse, 0, len(replies))
	for _, reply := range replies {
		result = append(result, reply.ToPublicResponseWithoutReplies(""))
	}

	return &PaginatedRepliesResponse{
		Replies: result,
		Total:   total,
		Limit:   limit,
		Skip:    skip,
		HasMore: int64(skip+len(replies)) < total,
	}, nil
}

// getReplyCountsForComments returns a map of commentId -> replyCount
func getReplyCountsForComments(commentIds []string) (map[string]int, error) {
	coll := mgm.Coll(&commentModel{})

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"parentId": bson.M{"$in": commentIds}}}},
		{{Key: "$group", Value: bson.M{
			"_id":   "$parentId",
			"count": bson.M{"$sum": 1},
		}}},
	}

	cursor, err := coll.Aggregate(mgm.Ctx(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(mgm.Ctx())

	var results []struct {
		ID    string `bson:"_id"`
		Count int    `bson:"count"`
	}
	if err := cursor.All(mgm.Ctx(), &results); err != nil {
		return nil, err
	}

	countMap := make(map[string]int)
	for _, r := range results {
		countMap[r.ID] = r.Count
	}

	return countMap, nil
}

// GetCommentsWithReplies fetches comments with their replies in a SINGLE query
// This solves the N+1 problem using MongoDB $lookup (like SQL JOIN)
// DEPRECATED: Use GetPaginatedComments instead for the new API
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

// ToPublicResponseWithoutReplies converts CommentWithReplies to public response without embedding replies
func (c *CommentWithReplies) ToPublicResponseWithoutReplies(currentUserEmail string) CommentPublicResponse {
	isOwn := currentUserEmail != "" && currentUserEmail == c.Email

	response := CommentPublicResponse{
		ID:         c.ID,
		Author:     c.Author,
		Gravatar:   c.Gravatar,
		Body:       c.Body,
		ParentID:   c.ParentID,
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

	return response
}

// commentModel is a helper struct for mgm.Coll()
type commentModel struct {
	mgm.DefaultModel `bson:",inline"`
}

func (c *commentModel) CollectionName() string {
	return "comments"
}
