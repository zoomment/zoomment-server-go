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
}

// CollectionName returns the MongoDB collection name
func (c *Comment) CollectionName() string {
	return "comments"
}
