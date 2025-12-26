package models

import (
	"github.com/kamva/mgm/v3"
)

// Vote represents a user's vote (upvote/downvote) on a comment
type Vote struct {
	mgm.DefaultModel `bson:",inline"`

	CommentID   string `bson:"commentId" json:"commentId"`
	Fingerprint string `bson:"fingerprint" json:"fingerprint"`
	Value       int    `bson:"value" json:"value"` // 1 = upvote, -1 = downvote
}

// CollectionName returns the MongoDB collection name
func (v *Vote) CollectionName() string {
	return "votes"
}

// VoteResponse is the response format for vote operations
type VoteResponse struct {
	CommentID string `json:"commentId"`
	Upvotes   int    `json:"upvotes"`
	Downvotes int    `json:"downvotes"`
	Score     int    `json:"score"`
	UserVote  int    `json:"userVote"`
}

// BulkVoteResponse is the response format for bulk vote queries
type BulkVoteResponse map[string]VoteResponse

