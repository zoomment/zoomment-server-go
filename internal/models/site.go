package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Site represents a registered website
type Site struct {
	BaseModel `bson:",inline"`

	UserID   primitive.ObjectID `bson:"userId" json:"userId"`
	Domain   string             `bson:"domain" json:"domain"`
	Verified bool               `bson:"verified" json:"verified"`
}

// CollectionName returns the MongoDB collection name
func (s *Site) CollectionName() string {
	return "sites"
}
