package models

import (
	"github.com/kamva/mgm/v3"
)

// Comment represents a comment on a page
type Comment struct {
	mgm.DefaultModel `bson:",inline"`

	// ParentID for threaded replies (nil = top-level comment)
	ParentID *string `bson:"parentId,omitempty" json:"parentId,omitempty"`

	// Author info
	Author   string `bson:"author" json:"author"`
	Email    string `bson:"email" json:"email"`
	Gravatar string `bson:"gravatar" json:"gravatar"`

	// Comment content
	Body string `bson:"body" json:"body"`

	// Location info
	Domain  string `bson:"domain" json:"domain"`
	PageURL string `bson:"pageUrl" json:"pageUrl"`
	PageID  string `bson:"pageId" json:"pageId"`

	// Status
	IsVerified bool `bson:"isVerified" json:"isVerified"`

	// Secret for guest deletion (not exposed in JSON)
	Secret string `bson:"secret" json:"-"`

	// Deprecated owner field (for backward compatibility)
	Owner *CommentOwner `bson:"owner,omitempty" json:"owner,omitempty"`
}

// CommentOwner is the deprecated nested owner structure
type CommentOwner struct {
	Name     string `bson:"name" json:"name"`
	Email    string `bson:"email" json:"email"`
	Gravatar string `bson:"gravatar" json:"gravatar"`
}

// CollectionName returns the MongoDB collection name
func (c *Comment) CollectionName() string {
	return "comments"
}

// CommentPublicData is what we return to clients (no sensitive data)
type CommentPublicData struct {
	ID         interface{}         `json:"_id"`
	Author     string              `json:"author"`
	Gravatar   string              `json:"gravatar"`
	Body       string              `json:"body"`
	ParentID   *string             `json:"parentId,omitempty"`
	IsVerified bool                `json:"isVerified"`
	IsOwn      bool                `json:"isOwn"`
	CreatedAt  interface{}         `json:"createdAt"`
	Owner      *CommentOwner       `json:"owner,omitempty"`
	Replies    []CommentPublicData `json:"replies,omitempty"`
}

// ToPublicData converts a Comment to public-safe data
func (c *Comment) ToPublicData(currentUserEmail string) CommentPublicData {
	isOwn := currentUserEmail != "" && currentUserEmail == c.Email

	return CommentPublicData{
		ID:         c.ID,
		Author:     c.Author,
		Gravatar:   c.Gravatar,
		Body:       c.Body,
		ParentID:   c.ParentID,
		IsVerified: c.IsVerified,
		IsOwn:      isOwn,
		CreatedAt:  c.CreatedAt,
		Owner: &CommentOwner{
			Name:     c.Author,
			Gravatar: c.Gravatar,
		},
	}
}
