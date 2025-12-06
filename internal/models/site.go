package models

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Site represents a registered website
type Site struct {
	mgm.DefaultModel `bson:",inline"`

	UserID   primitive.ObjectID `bson:"userId" json:"userId"`
	Domain   string             `bson:"domain" json:"domain"`
	Verified bool               `bson:"verified" json:"verified"`
}

// CollectionName returns the MongoDB collection name
func (s *Site) CollectionName() string {
	return "sites"
}

// NewSite creates a new site with default values
func NewSite(userID primitive.ObjectID, domain string) *Site {
	return &Site{
		UserID:   userID,
		Domain:   domain,
		Verified: false,
	}
}
